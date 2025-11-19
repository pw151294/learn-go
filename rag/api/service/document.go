package service

import (
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/models"
	"iflytek.com/weipan4/learn-go/rag/pkg/constants"
)

type DocumentService interface {
	Create(doc *models.Document) error
	GetByID(id string) (*models.Document, error)
	Update(doc *models.Document) error
	Delete(id string) error
	List(offset, limit int) ([]models.Document, error)
	GetAll() ([]models.Document, error)
	ChangeStatus(status string, docs []models.Document) error
	WaitingDocuments(limit int) ([]models.Document, error)
	UpdateStatus(id string, status string) error
}

type DocumentServiceImpl struct {
	DocumentRepo models.DocumentRepo
}

func NewDocumentService(repo models.DocumentRepo) DocumentService {
	return &DocumentServiceImpl{
		DocumentRepo: repo,
	}
}

func (s *DocumentServiceImpl) Create(doc *models.Document) error {
	return s.DocumentRepo.Create(doc)
}

func (s *DocumentServiceImpl) GetByID(id string) (*models.Document, error) {
	return s.DocumentRepo.GetByID(id)
}

func (s *DocumentServiceImpl) Update(doc *models.Document) error {
	return s.DocumentRepo.Update(doc)
}

func (s *DocumentServiceImpl) Delete(id string) error {
	return s.DocumentRepo.Delete(id)
}

func (s *DocumentServiceImpl) List(offset, limit int) ([]models.Document, error) {
	return s.DocumentRepo.List(offset, limit)
}

func (s *DocumentServiceImpl) GetAll() ([]models.Document, error) {
	return s.DocumentRepo.GetAll()
}

func (s *DocumentServiceImpl) ChangeStatus(status string, docs []models.Document) error {
	return s.DocumentRepo.ChangeStatus(status, docs)
}

func (s *DocumentServiceImpl) WaitingDocuments(limit int) ([]models.Document, error) {
	return s.DocumentRepo.GetByStatus(constants.DocumentWaitingIndexing, limit)
}

func (s *DocumentServiceImpl) UpdateStatus(id string, status string) error {
	return s.DocumentRepo.UpdateStatus(id, status)
}
