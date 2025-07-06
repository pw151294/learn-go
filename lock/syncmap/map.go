package syncmap

import (
	"sync"
	"time"
)

// SyncSet 使用sync_map实现一个sync_set
type SyncSet struct {
	setMap *sync.Map
}

func New() *SyncSet {
	return &SyncSet{&sync.Map{}}
}

func (s *SyncSet) Set(e string) {
	s.setMap.Store(e, struct{}{})
}

func (s *SyncSet) Del(e string) {
	s.setMap.Delete(e)
}

func (s *SyncSet) Exists(e string) bool {
	_, ok := s.setMap.Load(e)
	return ok
}

type timerStorage struct {
	idToTimer *sync.Map
}

func NewTimeStorage() *timerStorage {
	return &timerStorage{
		new(sync.Map),
	}
}

func (ts *timerStorage) Set(insId string, timer *time.Timer) {
	// 先删除再保存
	ts.idToTimer.Delete(insId)
	ts.idToTimer.Store(insId, timer)
}

func (ts *timerStorage) Del(insId string) {
	// 先停止计时器 再清空内存
	if v, ok := ts.idToTimer.Load(insId); ok {
		t := v.(*time.Timer)
		t.Stop()
	}
	ts.idToTimer.Delete(insId)
}

func (ts *timerStorage) Exists(insId string) bool {
	_, ok := ts.idToTimer.Load(insId)
	return ok
}
