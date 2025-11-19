package models

import (
	"time"

	"gorm.io/gorm"
)

type Dataset struct {
	ID                     string    `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	SubjectID              string    `gorm:"column:subject_id;type:varchar(64);not null" json:"subject_id"`
	Name                   string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description            string    `gorm:"column:description;type:varchar(500)" json:"description"`
	DataSourceType         string    `gorm:"column:data_source_type;type:varchar(255)" json:"data_source_type"`
	IndexingTechnique      string    `gorm:"column:indexing_technique;type:varchar(255)" json:"indexing_technique"`
	CreatedAt              time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy              string    `gorm:"column:created_by;type:varchar(64)" json:"created_by"`
	UpdatedAt              time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy              string    `gorm:"column:updated_by;type:varchar(64)" json:"updated_by"`
	EmbeddingModel         string    `gorm:"column:embedding_model;type:varchar(255)" json:"embedding_model"`
	EmbeddingModelProvider string    `gorm:"column:embedding_model_provider;type:varchar(255)" json:"embedding_model_provider"`
	RetrievalModel         string    `gorm:"column:retrieval_model;type:varchar(1000)" json:"retrieval_model"`
	BuiltInFieldEnabled    bool      `gorm:"column:built_in_field_enabled;type:tinyint(1);not null" json:"built_in_field_enabled"`
	ExtInfo                string    `gorm:"column:ext_info;type:varchar(1000)" json:"ext_info"`
	ProcessRule            string    `gorm:"column:process_rule;type:varchar(2000)" json:"process_rule"`
	Scope                  string    `gorm:"column:scope;type:varchar(32);not null" json:"scope"`
	Env                    string    `gorm:"column:env;type:varchar(20)" json:"env"`
}

// TableName 指定表名
func (Dataset) TableName() string {
	return "rag_dataset"
}

type DatasetRepo interface {
	Create(dataset *Dataset) error
	GetByID(id string) (*Dataset, error)
	Update(dataset *Dataset) error
	Delete(id string) error
	List(offset, limit int) ([]Dataset, error)
}

type datasetRepoImpl struct {
	db *gorm.DB
}

func NewDatasetRepo(db *gorm.DB) DatasetRepo {
	return &datasetRepoImpl{db: db}
}

func (r *datasetRepoImpl) Create(dataset *Dataset) error {
	return r.db.Create(dataset).Error
}

func (r *datasetRepoImpl) GetByID(id string) (*Dataset, error) {
	var dataset Dataset
	if err := r.db.First(&dataset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &dataset, nil
}

func (r *datasetRepoImpl) Update(dataset *Dataset) error {
	return r.db.Save(dataset).Error
}

func (r *datasetRepoImpl) Delete(id string) error {
	return r.db.Delete(&Dataset{}, "id = ?", id).Error
}

func (r *datasetRepoImpl) List(offset, limit int) ([]Dataset, error) {
	var datasets []Dataset
	if err := r.db.Offset(offset).Limit(limit).Find(&datasets).Error; err != nil {
		return nil, err
	}
	return datasets, nil
}
