package main

import (
	"github.com/streadway/amqp"
	"log"
)

func PublishExchange(exchangeName, types, routingKey, body string) error {
	conn, err := Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	//创建交换机
	err = ch.ExchangeDeclare(exchangeName, types, true, false, false, false, nil)
	if err != nil {
		log.Printf("create exchange failed: %v\n", err)
		return err
	}

	err = ch.Publish(exchangeName, routingKey, false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(body),
		DeliveryMode: amqp.Persistent,
	})

	return err
}

func ConsumeExchange(exchangeName, types, routingKey string, callback Callback) {
	conn, err := Connect()
	if err != nil {
		log.Printf("create connection failed: %v\n", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("create channel failed: %v\n", err)
		return
	}
	defer ch.Close()

	//创建交换机
	err = ch.ExchangeDeclare(exchangeName, types, true, false, false, false, nil)
	if err != nil {
		log.Printf("create exchange failed: %v\n", err)
		return
	}

	//创建队列并且和交换级绑定起来
	queue, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Printf("create queue failed: %v\n", err)
	}

	err = ch.QueueBind(queue.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		log.Printf("bind queue failed: %v\n", err)
		return
	}

	//消费消息
	deliveries, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("consume failed: %v\n", err)
		return
	}

	go func() {
		for d := range deliveries {
			body := string(d.Body)
			callback(body)
		}
	}()

	forever := make(chan struct{})
	log.Println("waiting for messages")
	<-forever
}
