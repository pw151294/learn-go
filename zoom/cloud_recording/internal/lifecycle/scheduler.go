package lifecycle

import (
	"context"
	"time"
)

type Scheduler struct {
	manager  *LifecycleManager
	interval time.Duration
	stopCh   chan struct{}
}

func NewScheduler(manager *LifecycleManager, intervalSeconds int) *Scheduler {
	return &Scheduler{
		manager:  manager,
		interval: time.Duration(intervalSeconds) * time.Second,
		stopCh:   make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.manager.RunOnce(context.Background())
			case <-s.stopCh:
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}
