package service

import (
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/models"
)

type DocumentSegmentService interface {
	Create(seg *models.DocumentSegment) error
	GetByID(id string) (*models.DocumentSegment, error)
	Update(seg *models.DocumentSegment) error
	Delete(id string) error
	List(offset, limit int) ([]models.DocumentSegment, error)
	BatchUpdateStatus(ids []string, status string) error
	ListByIds(ids []string) ([]models.DocumentSegment, error)
	InsertBatch(segs []models.DocumentSegment) error
}

type DocumentSegmentServiceImpl struct {
	DocumentSegmentRepo models.DocumentSegmentRepo
}

func NewDocumentSegmentService(repo models.DocumentSegmentRepo) DocumentSegmentService {
	return &DocumentSegmentServiceImpl{
		DocumentSegmentRepo: repo,
	}
}

func (s *DocumentSegmentServiceImpl) Create(seg *models.DocumentSegment) error {
	return s.DocumentSegmentRepo.Create(seg)
}

func (s *DocumentSegmentServiceImpl) GetByID(id string) (*models.DocumentSegment, error) {
	return s.DocumentSegmentRepo.GetByID(id)
}

func (s *DocumentSegmentServiceImpl) Update(seg *models.DocumentSegment) error {
	return s.DocumentSegmentRepo.Update(seg)
}

func (s *DocumentSegmentServiceImpl) Delete(id string) error {
	return s.DocumentSegmentRepo.Delete(id)
}

func (s *DocumentSegmentServiceImpl) List(offset, limit int) ([]models.DocumentSegment, error) {
	return s.DocumentSegmentRepo.List(offset, limit)
}

func (s *DocumentSegmentServiceImpl) BatchUpdateStatus(ids []string, status string) error {
	return s.DocumentSegmentRepo.BatchUpdateStatus(ids, status)
}

func (s *DocumentSegmentServiceImpl) ListByIds(ids []string) ([]models.DocumentSegment, error) {
	return s.DocumentSegmentRepo.ListByIds(ids)
}

func (s *DocumentSegmentServiceImpl) InsertBatch(segs []models.DocumentSegment) error {
	return s.DocumentSegmentRepo.InsertBatch(segs)
}
