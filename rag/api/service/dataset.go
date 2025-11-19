package service

import (
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/models"
)

type DatasetService interface {
	Create(dataset *models.Dataset) error
	GetByID(id string) (*models.Dataset, error)
	Update(dataset *models.Dataset) error
	Delete(id string) error
	List(offset, limit int) ([]models.Dataset, error)
}

type DatasetServiceImpl struct {
	DatasetRepo models.DatasetRepo
}

func NewDatasetServiceImpl(datasetRepo models.DatasetRepo) DatasetService {
	return &DatasetServiceImpl{
		DatasetRepo: datasetRepo,
	}
}

func (s *DatasetServiceImpl) Create(dataset *models.Dataset) error {
	return s.DatasetRepo.Create(dataset)
}

func (s *DatasetServiceImpl) GetByID(id string) (*models.Dataset, error) {
	return s.DatasetRepo.GetByID(id)
}

func (s *DatasetServiceImpl) Update(dataset *models.Dataset) error {
	return s.DatasetRepo.Update(dataset)
}

func (s *DatasetServiceImpl) Delete(id string) error {
	return s.DatasetRepo.Delete(id)
}

func (s *DatasetServiceImpl) List(offset, limit int) ([]models.Dataset, error) {
	return s.DatasetRepo.List(offset, limit)
}
