package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {
	//设置一个调度协程 定时检查所有的Cron任务 谁过期了就执行谁
	e1 := cronexpr.MustParse("*/5 * * * * *")
	cronJob1 := &CronJob{expr: e1, nextTime: e1.Next(time.Now())}

	e2 := cronexpr.MustParse("*/5 * * * * *")
	cronJob2 := &CronJob{expr: e2, nextTime: e2.Next(time.Now())}

	//任务注册到调度表
	scheduleMap := make(map[string]*CronJob)
	scheduleMap["job1"] = cronJob1
	scheduleMap["job2"] = cronJob2

	//启动一个调度的协程
	go func() {
		for {
			now := time.Now()
			for jobName, cronJob := range scheduleMap {
				// 判断调度的任务是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//启动一个协程 执行这个任务
					go func(jobName string) {
						fmt.Printf("job %v executed\n", jobName)
					}(jobName)

					//计算下一次调度时间
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println("next time: ", cronJob.nextTime)
				}
			}

			//睡眠100毫秒
			timer := time.NewTimer(100 * time.Millisecond)
			select {
			case <-timer.C: //将在100毫秒时可读 返回
			}
		}
	}()

	time.Sleep(5 * time.Minute)
}
