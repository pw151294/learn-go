package nacos

import (
	"iflytek.com/weipan4/learn-go/net/host"
	"strconv"
)

const (
	Status            = "active"
	ActiveConnections = 200
)

type MetadataOptions func(*InstanceMetadata)

type InstanceMetadata struct {
	InstanceId        string `json:"instanceId"`
	Ip                string `json:"ip"`
	Port              int    `json:"port"`
	Status            string `json:"status"`
	ActiveConnections int    `json:"activeConnections"`
}

func NewMetadata(opts ...MetadataOptions) map[string]string {
	instMeta := newInstanceMetadata(opts...)
	if instMeta == nil {
		return nil
	}

	return map[string]string{
		"instanceId":  instMeta.InstanceId,
		"ip":          instMeta.Ip,
		"port":        strconv.Itoa(instMeta.Port),
		"status":      instMeta.Status,
		"activeConns": strconv.Itoa(instMeta.ActiveConnections),
	}
}

func newInstanceMetadata(opts ...MetadataOptions) *InstanceMetadata {
	instMeta := &InstanceMetadata{
		Ip:                host.HostInfo.GetIP(),
		Status:            Status,
		ActiveConnections: ActiveConnections,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(instMeta)
		}
	}

	return instMeta
}

func WithInstanceId(id string) MetadataOptions {
	return func(i *InstanceMetadata) {
		i.InstanceId = id
	}
}

func WithIp(ip string) MetadataOptions {
	return func(i *InstanceMetadata) {
		i.Ip = ip
	}
}

func WithPort(port int) MetadataOptions {
	return func(i *InstanceMetadata) {
		i.Port = port
	}
}

func WithStatus(status string) MetadataOptions {
	return func(i *InstanceMetadata) {
		i.Status = status
	}
}

func WithActiveConnections(activeConnections int) MetadataOptions {
	return func(i *InstanceMetadata) {
		i.ActiveConnections = activeConnections
	}
}
