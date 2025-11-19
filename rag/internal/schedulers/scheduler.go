package schedulers

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"iflytek.com/weipan4/learn-go/rag/api/service"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/embedding"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/models"
	"iflytek.com/weipan4/learn-go/rag/pkg/loggers"
)

type Scheduler struct {
	splittingScheduler *SplittingScheduler
	embeddingScheduler *EmbeddingScheduler
	storeScheduler     *EmbeddingStoreScheduler

	embeddingClient embedding.EmbeddingClient
	docService      service.DocumentService
	segmentService  service.DocumentSegmentService
	messageChan     chan Message
	embeddingChan   chan Embedding

	WorkerNum int
	period    int // 文本切分/向量化定时任务周期
}

func NewScheduler(db *gorm.DB, embeddingClient embedding.EmbeddingClient, workNum, capacity, period int) *Scheduler {
	docRepo := models.NewDocumentRepo(db)
	segmentRepo := models.NewDocumentSegmentRepo(db)
	docService := service.NewDocumentService(docRepo)
	segmentService := service.NewDocumentSegmentService(segmentRepo)

	return &Scheduler{
		docService:      docService,
		segmentService:  segmentService,
		embeddingClient: embeddingClient,
		messageChan:     make(chan Message, capacity),
		embeddingChan:   make(chan Embedding, capacity),
		WorkerNum:       workNum,
		period:          period,
	}
}

func (s *Scheduler) Start() {
	// 1. 初始化文本切分定时任务(producer)
	s.splittingScheduler = NewSplittingScheduler(s.docService, s.segmentService, s.messageChan)
	go s.splittingScheduler.StartSplittingTextTask(time.Duration(s.period) * time.Millisecond)

	// 2. 初始化文本向量化任务
	s.embeddingScheduler = NewEmbeddingScheduler(s.embeddingClient, s.docService, s.segmentService, s.messageChan, s.embeddingChan)
	rand.NewSource(time.Now().UnixNano())
	for i := range s.WorkerNum {
		go func() {
			loggers.Info(fmt.Sprintf("embedding task %d start", i+1))
			s.embeddingScheduler.startEmbeddingTextTask(s.messageChan)
		}()
	}

	// 3.初始化向量化结果存储任务
	s.storeScheduler = NewEmbeddingStoreScheduler(s.docService, s.segmentService, s.embeddingChan)
	for i := range s.WorkerNum {
		go func() {
			loggers.Info(fmt.Sprintf("embedding store task %d start", i+1))
			s.storeScheduler.startEmbeddingStoreTask(s.embeddingChan)
		}()
	}

	never := make(chan struct{})
	<-never
}

func (s *Scheduler) Close() {
	close(s.messageChan)
	close(s.embeddingChan)
}
