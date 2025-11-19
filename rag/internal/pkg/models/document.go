package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	ID               string     `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	SubjectID        string     `gorm:"column:subject_id;type:varchar(64);not null" json:"subject_id"`
	DatasetID        string     `gorm:"column:dataset_id;type:varchar(64);not null" json:"dataset_id"`
	Position         int        `gorm:"column:position;not null" json:"position"`
	Name             string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	DocForm          string     `gorm:"column:doc_form;type:varchar(255);not null" json:"doc_form"`
	DataSourceType   string     `gorm:"column:data_source_type;type:varchar(255);not null" json:"data_source_type"`
	DataSourceInfo   string     `gorm:"column:data_source_info;type:varchar(1000)" json:"data_source_info"`
	Batch            string     `gorm:"column:batch;type:varchar(255);not null" json:"batch"`
	WordCount        int        `gorm:"column:word_count;not null" json:"word_count"`
	Tokens           int        `gorm:"column:tokens;not null" json:"tokens"`
	Status           string     `gorm:"column:status;type:varchar(255);not null" json:"status"`
	Error            string     `gorm:"column:error;type:text" json:"error"`
	StatusChangeTime time.Time  `gorm:"column:status_change_time;autoCreateTime" json:"status_change_time"`
	Enabled          bool       `gorm:"column:enabled;type:tinyint(1);not null" json:"enabled"`
	DocLanguage      string     `gorm:"column:doc_language;type:varchar(255)" json:"doc_language"`
	DocMetadata      string     `gorm:"column:doc_metadata;type:text" json:"doc_metadata"`
	IndexingLog      string     `gorm:"column:indexing_log;type:varchar(2000)" json:"indexing_log"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy        string     `gorm:"column:created_by;type:varchar(64)" json:"created_by"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy        string     `gorm:"column:updated_by;type:varchar(64)" json:"updated_by"`
	DisabledAt       *time.Time `gorm:"column:disabled_at" json:"disabled_at"`
	DisabledBy       string     `gorm:"column:disabled_by;type:varchar(64)" json:"disabled_by"`
	ExtInfo          string     `gorm:"column:ext_info;type:varchar(1000)" json:"ext_info"`
	Scope            string     `gorm:"column:scope;type:varchar(32);not null" json:"scope"`
	Env              string     `gorm:"column:env;type:varchar(20)" json:"env"`
	StatusWhenError  string     `gorm:"column:status_when_error;type:varchar(255)" json:"status_when_error"`
}

// TableName 指定表名
func (Document) TableName() string {
	return "rag_document"
}

type DocumentRepo interface {
	Create(doc *Document) error
	GetByID(id string) (*Document, error)
	Update(doc *Document) error
	Delete(id string) error
	List(offset, limit int) ([]Document, error)
	GetAll() ([]Document, error)
	GetByStatus(status string, limit int) ([]Document, error)
	ChangeStatus(status string, docs []Document) error
	UpdateStatus(id string, status string) error
}

type documentRepoImpl struct {
	db *gorm.DB
}

func NewDocumentRepo(db *gorm.DB) DocumentRepo {
	return &documentRepoImpl{db: db}
}

func (r *documentRepoImpl) Create(doc *Document) error {
	return r.db.Create(doc).Error
}

func (r *documentRepoImpl) GetByID(id string) (*Document, error) {
	var doc Document
	if err := r.db.First(&doc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *documentRepoImpl) Update(doc *Document) error {
	return r.db.Save(doc).Error
}

func (r *documentRepoImpl) Delete(id string) error {
	return r.db.Delete(&Document{}, "id = ?", id).Error
}

func (r *documentRepoImpl) List(offset, limit int) ([]Document, error) {
	var docs []Document
	if err := r.db.Offset(offset).Limit(limit).Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *documentRepoImpl) GetAll() ([]Document, error) {
	var docs []Document
	if err := r.db.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *documentRepoImpl) GetByStatus(status string, limit int) ([]Document, error) {
	var docs []Document
	query := r.db.Where("status = ?", status)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *documentRepoImpl) ChangeStatus(status string, docs []Document) error {
	now := time.Now()
	ids := make([]string, 0, len(docs))
	for _, doc := range docs {
		ids = append(ids, doc.ID)
	}
	return r.db.Model(&Document{}).
		Where("id IN ?", ids).
		Updates(map[string]interface{}{
			"status":             status,
			"status_change_time": now,
			"updated_at":         now,
		}).Error
}

func (r *documentRepoImpl) UpdateStatus(id string, status string) error {
	now := time.Now()
	return r.db.Model(&Document{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":             status,
			"status_change_time": now,
			"updated_at":         now,
		}).Error
}
