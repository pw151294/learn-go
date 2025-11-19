package models

import (
	"time"

	"gorm.io/gorm"
)

type DocumentSegment struct {
	ID               string     `gorm:"column:id;primaryKey;type:varchar(64)" json:"id"`
	SubjectID        string     `gorm:"column:subject_id;type:varchar(64);not null" json:"subject_id"`
	DatasetID        string     `gorm:"column:dataset_id;type:varchar(64);not null" json:"dataset_id"`
	DocumentID       string     `gorm:"column:document_id;type:varchar(64);not null" json:"document_id"`
	ParentPosition   int        `gorm:"column:parent_position;not null" json:"parent_position"`
	ParentGroup      string     `gorm:"column:parent_group;type:varchar(64)" json:"parent_group"`
	SegmentPosition  int        `gorm:"column:segment_position;not null" json:"segment_position"`
	Content          string     `gorm:"column:content;type:text;not null" json:"content"`
	Answer           string     `gorm:"column:answer;type:text" json:"answer"`
	ImageURL         string     `gorm:"column:image_url;type:varchar(255)" json:"image_url"`
	BBoxPosition     string     `gorm:"column:bbox_position;type:text" json:"bbox_position"`
	WordCount        int        `gorm:"column:word_count;not null" json:"word_count"`
	IndexNodeID      string     `gorm:"column:index_node_id;type:varchar(255)" json:"index_node_id"`
	IndexNodeHash    string     `gorm:"column:index_node_hash;type:varchar(255)" json:"index_node_hash"`
	HitCount         int        `gorm:"column:hit_count;not null" json:"hit_count"`
	Enabled          bool       `gorm:"column:enabled;type:tinyint(1);not null" json:"enabled"`
	Status           string     `gorm:"column:status;type:varchar(255);not null" json:"status"`
	StatusChangeTime time.Time  `gorm:"column:status_change_time;autoCreateTime" json:"status_change_time"`
	Error            string     `gorm:"column:error;type:text" json:"error"`
	CreatedAt        time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	CreatedBy        string     `gorm:"column:created_by;type:varchar(64)" json:"created_by"`
	UpdatedAt        time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UpdatedBy        string     `gorm:"column:updated_by;type:varchar(64)" json:"updated_by"`
	DisabledAt       *time.Time `gorm:"column:disabled_at" json:"disabled_at"`
	DisabledBy       string     `gorm:"column:disabled_by;type:varchar(64)" json:"disabled_by"`
	IndexingLog      string     `gorm:"column:indexing_log;type:varchar(2000)" json:"indexing_log"`
	ExtInfo          string     `gorm:"column:ext_info;type:varchar(1000)" json:"ext_info"`
	Scope            string     `gorm:"column:scope;type:varchar(32);not null" json:"scope"`
	Env              string     `gorm:"column:env;type:varchar(20)" json:"env"`
	StatusWhenError  string     `gorm:"column:status_when_error;type:varchar(255)" json:"status_when_error"`
}

// TableName 指定表名
func (DocumentSegment) TableName() string {
	return "rag_document_segment"
}

type DocumentSegmentRepo interface {
	Create(seg *DocumentSegment) error
	GetByID(id string) (*DocumentSegment, error)
	Update(seg *DocumentSegment) error
	Delete(id string) error
	List(offset, limit int) ([]DocumentSegment, error)
	BatchUpdateStatus(ids []string, status string) error
	InsertBatch(segs []DocumentSegment) error
	ListByIds(ids []string) ([]DocumentSegment, error)
}

type documentSegmentRepoImpl struct {
	db *gorm.DB
}

func NewDocumentSegmentRepo(db *gorm.DB) DocumentSegmentRepo {
	return &documentSegmentRepoImpl{db: db}
}

func (r *documentSegmentRepoImpl) Create(seg *DocumentSegment) error {
	return r.db.Create(seg).Error
}

func (r *documentSegmentRepoImpl) GetByID(id string) (*DocumentSegment, error) {
	var seg DocumentSegment
	if err := r.db.First(&seg, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &seg, nil
}

func (r *documentSegmentRepoImpl) Update(seg *DocumentSegment) error {
	return r.db.Save(seg).Error
}

func (r *documentSegmentRepoImpl) Delete(id string) error {
	return r.db.Delete(&DocumentSegment{}, "id = ?", id).Error
}

func (r *documentSegmentRepoImpl) List(offset, limit int) ([]DocumentSegment, error) {
	var segs []DocumentSegment
	if err := r.db.Offset(offset).Limit(limit).Find(&segs).Error; err != nil {
		return nil, err
	}
	return segs, nil
}

func (r *documentSegmentRepoImpl) BatchUpdateStatus(ids []string, status string) error {
	return r.db.Model(&DocumentSegment{}).
		Where("id IN ?", ids).
		Update("status", status).Error
}

func (r *documentSegmentRepoImpl) InsertBatch(segs []DocumentSegment) error {
	return r.db.Create(&segs).Error
}

func (r *documentSegmentRepoImpl) ListByIds(ids []string) ([]DocumentSegment, error) {
	var segs []DocumentSegment
	if err := r.db.Where("id IN ?", ids).Find(&segs).Error; err != nil {
		return nil, err
	}
	return segs, nil
}
