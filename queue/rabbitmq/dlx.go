package main

import (
	"github.com/streadway/amqp"
	"log"
)

const dlExchangeTypes = "fanout"

// PublishDlx 死信队列生产端
func PublishDlx(exchangeAName, body string) error {
	conn, err := Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	//将消息投递到队列A中
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	//投递消息给交换机A
	err = ch.Publish(exchangeAName, "", false, false, amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(body),
		DeliveryMode: amqp.Persistent,
	})

	return err
}

// ConsumeDlx 死信队列
func ConsumeDlx(exchangeAName, queueAName, exchangeBName, queueBName string, ttl int, callback Callback) {
	//建立连接
	conn, err := Connect()
	if err != nil {
		log.Printf("create connect failed: %v\n", err)
		return
	}
	defer conn.Close()

	//创建channel
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("create channel failed: %v\n", err)
		return
	}
	defer ch.Close()

	//创建A交换机和A队列
	err = ch.ExchangeDeclare(exchangeAName, dlExchangeTypes, true, false, false, false, nil)
	if err != nil {
		log.Printf("create exchange A failed: %v\n", err)
		return
	}
	queueA, err := ch.QueueDeclare(queueAName, true, false, false, false,
		amqp.Table{
			"x-message-ttl":          int64(ttl),
			"x-dead-letter-exchange": exchangeBName,
			"x-dead-letter-queue":    queueBName,
		})
	if err != nil {
		log.Printf("create queue failed: %v\n", err)
		return
	}
	err = ch.QueueBind(queueA.Name, "", exchangeAName, false, nil)
	if err != nil {
		log.Printf("bind queue A to exchange A failed: %v\n", err)
		return
	}

	//创建B交换机和B队列
	err = ch.ExchangeDeclare(exchangeBName, dlExchangeTypes, true, false, false, false, nil)
	if err != nil {
		log.Printf("create exchange B failed: %v\n", err)
		return
	}
	queueB, err := ch.QueueDeclare(queueBName, true, false, false, false, nil)
	if err != nil {
		log.Printf("create queue B failed: %v\n", err)
		return
	}
	err = ch.QueueBind(queueB.Name, "", exchangeBName, false, nil)
	if err != nil {
		log.Printf("bind queue B to exchange B failed: %v\n", err)
		return
	}

	deliveries, err := ch.Consume(queueB.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("consume failed: %v\n", err)
		return
	}

	go func() {
		for d := range deliveries {
			body := string(d.Body)
			callback(body)
			d.Ack(false)
		}
	}()

	log.Printf("wating for messages......")
	forever := make(chan struct{})
	<-forever
}
