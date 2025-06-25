package main

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	defaultProbeAddr = "223.5.5.5:80"
)

var HostInfo *HostInfoCache

type HostInfoCache struct {
	name string
	ip   string
	sync.RWMutex
}

func (c *HostInfoCache) GetHostname() string {
	c.RLock()
	n := c.name
	c.RUnlock()
	return n
}

func (c *HostInfoCache) GetIP() string {
	c.RLock()
	defer c.RUnlock()
	ip := c.ip
	return ip
}

func (c *HostInfoCache) SetHostname(name string) {
	if name == c.GetHostname() {
		return
	}

	c.Lock()
	c.name = name
	c.Unlock()
}

func (c *HostInfoCache) SetIP(ip string) {
	if ip == c.GetIP() {
		return
	}

	c.Lock()
	c.ip = ip
	c.Unlock()
}

func InitHostInfo() {
	hostname, err := os.Hostname()
	if err != nil {
		panic("初始化hostname失败")
	}

	ip, err := GetOutboundIP(nil)
	if err != nil {
		panic("初始化hostIP失败")
	}

	HostInfo = &HostInfoCache{
		name: hostname,
		ip:   fmt.Sprint(ip),
	}

	go HostInfo.update()
}

func (c *HostInfoCache) update() {
	for {
		time.Sleep(time.Minute)
		ip, err := GetOutboundIP(nil)
		if err != nil {
			log.Println("E! failed to get ip:", err)
		} else {
			HostInfo.SetIP(ip.String())
		}
		name, err := os.Hostname()
		if err != nil {
			log.Println("E! failed to get hostname:", err)
		} else {
			HostInfo.SetHostname(name)
			// 防止并发
			copyMap := make(map[string]string)
			copyMap["ident"] = name
			copyMap["datasleuth_ip"] = c.ip

		}

	}
}

func GetOutboundIP(urls []string) (net.IP, error) {
	addr := defaultProbeAddr
	for _, v := range urls {
		if len(v) != 0 {
			u, err := url.Parse(v)
			if err != nil {
				log.Printf("W! parse writers url %s error %s", v, err)
				continue
			} else {
				if strings.Contains(u.Host, "localhost") || strings.Contains(u.Host, "127.0.0.1") {
					continue
				}
				if len(u.Port()) == 0 {
					if u.Scheme == "http" {
						u.Host = u.Host + ":80"
					}
					if u.Scheme == "https" {
						u.Host = u.Host + ":443"
					}
				}
				addr = u.Host
				break
			}
		}
	}

	conn, err := net.Dial("udp", addr)
	if err != nil {
		ip, err := getLocalIP()
		if err != nil {
			return nil, fmt.Errorf("failed to get local ip: %v", err)
		}
		return ip, nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func getLocalIP() (net.IP, error) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifs {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Println("W! iface address error", err)
			continue
		}
		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok && ip.IP.IsLoopback() {
				continue
			} else {
				ip4 := ip.IP.To4()
				if ip4 == nil {
					continue
				}
				return ip4, nil
			}
		}
	}
	return nil, fmt.Errorf("no local ip found")
}
