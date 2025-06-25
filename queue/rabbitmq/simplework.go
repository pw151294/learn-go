package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"path/filepath"
)

const (
	configDir  = "queue/rabbitmq"
	configFile = "rabbitmq.json"
	url        = "amqp://%s:%s@%s:%d/"
)

type Callback func(msg string)

type RabbitMQ struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func Connect() (*amqp.Connection, error) {
	//读取rabbitmq的配置文件
	fn := filepath.Join(configDir, configFile)
	bytes, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	var cfg RabbitMQ
	if err = json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	conn, err := amqp.Dial(fmt.Sprintf(url, cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Publish 发送端函数
func Publish(exchangeName, queueName, body string) error {
	conn, err := Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	//声明消息队列
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	//发送消息
	err = ch.Publish(exchangeName, queue.Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})

	return err
}

// Consume 接收者的方法
func Consume(queueName string, callback Callback) {
	conn, err := Connect()
	if err != nil {
		fmt.Printf("rabbitmq connect failed: %v\n", err)
		return
	}
	defer conn.Close()

	//创建通道channel
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("create rabbitmq channel failed: %v\n", err)
		return
	}
	defer ch.Close()

	//创建queue队列
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("declare rabbitmq failed: %v\n", err)
		return
	}

	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("consume message failed: %v\n", err)
		return
	}

	//这是个守护进程不能退出
	go func() {
		for d := range msgs {
			s := string(d.Body)
			callback(s)
			d.Ack(false)
		}
	}()

	fmt.Printf("waiting for messages")
	forever := make(chan struct{})
	<-forever
}
