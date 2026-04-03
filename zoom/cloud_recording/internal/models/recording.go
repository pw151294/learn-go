package models

import "time"

// Tier 表示存储层级
const (
	TierHot  = "hot"
	TierWarm = "warm"
	TierCold = "cold"
)

// Recording 表示一个录制文件
type Recording struct {
	RecordingID    string    `json:"recording_id"`
	Title          string    `json:"title"`
	UserID         string    `json:"user_id"`
	RoomID         string    `json:"room_id"`
	Bucket         string    `json:"bucket"`
	ObjectKey      string    `json:"object_key"`  // 格式：{userID}/{roomID}/{recordingID}
	Size           int64     `json:"size"`
	Duration       int       `json:"duration"`    // 秒
	Status         string    `json:"status"`      // hot/warm/cold/deleted
	CreatedAt      time.Time `json:"created_at"`
	LastAccessAt   time.Time `json:"last_access_at"`
	TierMigratedAt time.Time `json:"tier_migrated_at"`
}

// SearchQuery 搜索查询参数
type SearchQuery struct {
	Keyword string `form:"keyword"`
	UserID  string `form:"user_id"`
	RoomID  string `form:"room_id"`
	Status  string `form:"status"`
	Page    int    `form:"page"`
	Size    int    `form:"size"`
}
