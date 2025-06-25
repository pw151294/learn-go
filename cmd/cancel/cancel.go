package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

func main() {
	//创建一个resultChan用来接收结果
	resultChan := make(chan *result, 1000)

	//执行一个cmd 让它在协程里面去执行 让它执行2秒
	//在1秒的时候杀死这个cmd
	ctx, cancelFunc := context.WithCancel(context.TODO())

	//context对象中内置了一个channel
	//cancelFunc执行的时候 会将context对象里的channel关闭掉

	go func() {
		cmd := exec.CommandContext(ctx, "bash", "-c", "sleep 2;echo hello")
		//cmd对象内部通过select语句来监听对象channel是否被关闭
		//select {
		//		case <-ctx.Done(): //监听到了done之后就会执行kill指令来杀死子进程
		//		}

		//执行任务 捕获输出
		output, err := cmd.CombinedOutput()

		//任务执行的输出传递给main写成
		resultChan <- &result{err: err, output: output}
	}()

	//继续往下走
	time.Sleep(1 * time.Second)

	//取消上下文
	cancelFunc()

	//在main协程里等待子协程的退出 并打印出任务执行的结果
	res := <-resultChan

	fmt.Println(res.err, string(res.output))
}
