package main

import (
	"context"
	"fmt"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"sync"
)

type Scheduler struct {
	workNum  int
	ctx      context.Context
	cancel   context.CancelFunc
	taskCh   chan Task
	resultCh chan *AdResult
	errCh    chan *Error
	wg       sync.WaitGroup
}

func NewScheduler(workNum, queueSize int) *Scheduler {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &Scheduler{
		workNum:  workNum,
		ctx:      ctx,
		cancel:   cancelFunc,
		taskCh:   make(chan Task, queueSize),
		resultCh: make(chan *AdResult, queueSize),
		errCh:    make(chan *Error, queueSize),
	}
}

func (s *Scheduler) Start() {
	for i := 0; i < s.workNum; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

func (s *Scheduler) Submit(task Task) {
	s.taskCh <- task
}

func (s *Scheduler) worker(idx int) {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			zap.GetLogger().Error("scheduler exit", "scheduler", idx)
			return
		case task := <-s.taskCh:
			if res, err := task.Execute(s.ctx); err != nil {
				s.errCh <- &Error{taskID: task.ID(), err: fmt.Errorf("task execute failed: %w", err)}
			} else {
				s.resultCh <- res
				zap.GetLogger().Info("task execute success",
					"task id", task.ID(), "worker", idx)
			}
		}
	}
}

func (s *Scheduler) Wait() {
	s.wg.Wait()
}

func (s *Scheduler) Stop() {
	s.cancel()
	s.Wait()
	close(s.taskCh)
	close(s.resultCh)
	close(s.errCh)
	zap.GetLogger().Info("scheduler stopped")
}
