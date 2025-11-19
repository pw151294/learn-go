package main

import (
	"bytes"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

const (
	host     = "172.30.34.73:22"
	user     = "root"
	password = "qwe!@#2022"
	cmd      = "ls -la /data/weipan4"
)

func main() {
	// 建立ssh连接远程服务器
	client, err := ssh.Dial("tcp", host, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建新会话
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	var buffer bytes.Buffer
	session.Stdout = &buffer
	if err := session.Run(cmd); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(buffer.String())
	}
}
