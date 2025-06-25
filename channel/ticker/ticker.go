package main

import (
	"log"
	"time"
)

func main() {
	t := time.NewTicker(3 * time.Second)
	defer t.Stop()
	ch := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Second)
		ch <- struct{}{}
	}()

loop:
	for {
		select {
		case <-ch:
			log.Println("done")
			break loop
		case <-t.C:
			log.Printf("current time is %s", time.Now())
		}
	}
}
