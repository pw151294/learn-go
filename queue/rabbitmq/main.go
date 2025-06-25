package main

import (
	"fmt"
	"log"
	"time"
)

const (
	exchangeName = ""
	queueName    = "learn-go"

	subscribeExchangeName = "learn.go.fanout"
	subscribeTypes        = "fanout"

	routingExchangeName = "learn.go.routing"
	routingTypes        = "direct"
	routingKey1         = "one"
	routingKey2         = "two"

	topicExchangeName = "learn.go.topic"
	topicTypes        = "topic"
	topic1RoutingKey1 = "topic1.key1"
	topic1RoutingKey2 = "topic1.key2"
	topic2RoutingKey1 = "topic2.key1"
	topic2RoutingKey2 = "topic2.key2"
	consumerGroup1    = "topic1.*"
	consumerGroup2    = "topic2.*"

	dlExchangeAName = "learn.go.dl.A"
	dlExchangeBName = "learn.go.dl.B"
	dlQueueAName    = "learn-go-dlx-A"
	dlQueueBName    = "learn-go-dlx-B"
	dlQueueTTL      = 10000
)

func dlx() {
	go func() {
		count := 0
		for {
			err := PublishDlx(dlExchangeAName,
				fmt.Sprintf("dead letter message, exchange name:%s, body:%d", dlExchangeAName, count))
			if err != nil {
				log.Printf("publish dead letter message failed, exchange name:%s", dlExchangeAName)
			}
			count++
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		count := 0
		for {
			err := PublishExchange(dlExchangeBName, subscribeTypes, "",
				fmt.Sprintf("dead letter message, exchange name:%s, body:%d", dlExchangeBName, count))
			if err != nil {
				log.Printf("publish dead letter message failed, exchange name:%s", dlExchangeBName)
			}
			count++
			time.Sleep(2 * time.Second)
		}
	}()

	ConsumeDlx(dlExchangeAName, dlQueueAName, dlExchangeBName, dlQueueBName, dlQueueTTL, func(msg string) {
		log.Printf("received message from dead letter queue, queue name:%s, body:%s", dlQueueBName, msg)
	})
	forever := make(chan struct{})
	<-forever
}

func topic() {
	//topic1消息
	go func() {
		count := 0
		for {
			if count%2 == 0 {
				err := PublishExchange(topicExchangeName, topicTypes,
					topic1RoutingKey1, fmt.Sprintf("topic message body:%d, routing key:%s", count, topic1RoutingKey1))
				if err != nil {
					log.Printf("publish message failed, routing key:%s", topic1RoutingKey1)
				}
			} else {
				err := PublishExchange(topicExchangeName, topicTypes,
					topic1RoutingKey2, fmt.Sprintf("topic message body:%d, routing key:%s", count, topic1RoutingKey2))
				if err != nil {
					log.Printf("publish message failed, routing key:%s", topic1RoutingKey2)
				}
			}
			count++
			time.Sleep(2 * time.Second)
		}
	}()
	//topic2消息
	go func() {
		count := 0
		for {
			if count%2 == 0 {
				err := PublishExchange(topicExchangeName, topicTypes,
					topic2RoutingKey1, fmt.Sprintf("topic message body:%d, routing key:%s", count, topic2RoutingKey1))
				if err != nil {
					log.Printf("publish message failed, routing key:%s", topic2RoutingKey1)
				}
			} else {
				err := PublishExchange(topicExchangeName, topicTypes,
					topic2RoutingKey2, fmt.Sprintf("topic message body:%d, routing key:%s", count, topic2RoutingKey2))
				if err != nil {
					log.Printf("publish message failed, routing key:%s", topic2RoutingKey2)
				}
			}
			count++
			time.Sleep(2 * time.Second)
		}
	}()

	//创建消费者
	go ConsumeExchange(topicExchangeName, topicTypes, consumerGroup1, func(msg string) {
		log.Printf("received message, topic:%s, body:%s", consumerGroup1, msg)
	})
	go ConsumeExchange(topicExchangeName, topicTypes, consumerGroup2, func(msg string) {
		log.Printf("received message, topic:%s, body:%s", consumerGroup2, msg)
	})

	log.Printf("begin cosuming message......")
	forever := make(chan struct{})
	<-forever
}

// 订阅模式消费
func subscribe() {
	go func() {
		count := 0
		for {
			err := PublishExchange(exchangeName, subscribeTypes,
				"", fmt.Sprintf("subscribe message body: %d", count))
			if err != nil {
				log.Printf("publish message failed: %v", err)
			}
			time.Sleep(1 * time.Second)
			count++
		}
	}()

	ConsumeExchange(subscribeExchangeName, subscribeTypes, "", func(msg string) {
		log.Printf("received message: %s", msg)
	})
}

func routing() {
	go func() {
		count := 1
		for {
			err := PublishExchange(routingExchangeName, routingTypes, routingKey1, fmt.Sprintf("routing message body: %d", 2*count-1))
			if err != nil {
				log.Printf("publish message failed, exchange name:%s, routing key:%s, err:%v", routingExchangeName, routingKey1, err)
			}
			err = PublishExchange(routingExchangeName, routingTypes, routingKey2, fmt.Sprintf("routing message body: %d", 2*count))
			if err != nil {
				log.Printf("publish message failed, exchange name:%s, routing key:%s, err:%v", routingExchangeName, routingKey2, err)
			}
			time.Sleep(1 * time.Second)
			count++
		}
	}()

	//创建消费者消费消息
	go ConsumeExchange(routingExchangeName, routingTypes, routingKey1, func(msg string) {
		log.Printf("received message, routing key: %s, message: %s", routingKey1, msg)
	})
	go ConsumeExchange(routingExchangeName, routingTypes, routingKey2, func(msg string) {
		log.Printf("received message, routing key: %s, message: %s", routingKey2, msg)
	})

	forever := make(chan struct{})
	<-forever
}

func work() {
	count := 0
	go func() {
		for {
			err := Publish(exchangeName, queueName, fmt.Sprintf("message body %d", count))
			if err != nil {
				fmt.Printf("failed to publish message: %v", err)
			}
			count++
			time.Sleep(1 * time.Second)
		}
	}()

	//Work工作模式消费消息
	go Consume(queueName, func(msg string) {
		log.Printf("first consumer received message: %s", msg)
	})
	go Consume(queueName, func(msg string) {
		log.Printf("second consumer received message: %s", msg)
	})
	forever := make(chan struct{})
	<-forever
}

func main() {
	//work()
	//subscribe()
	//routing()
	//topic()
	dlx()
}
