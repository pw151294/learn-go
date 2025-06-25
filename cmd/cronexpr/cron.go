package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"log"
	"time"
)

func main() {
	var expr *cronexpr.Expression
	expr, err := cronexpr.Parse("1-2 * * * *")
	if err != nil {
		log.Fatalf("parse err: %v", err)
	}

	//设置当前时间和下次被调度时间
	now := time.Now()
	timeNext := expr.Next(now)
	log.Println("timeNext:", timeNext)

	//等待这个定时器被调度
	time.AfterFunc(timeNext.Sub(now), func() {
		fmt.Printf("被调度了：%v", timeNext)
	})

	time.Sleep(5 * time.Second)
}
