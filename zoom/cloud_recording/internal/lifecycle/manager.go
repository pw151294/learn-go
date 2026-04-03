package lifecycle

import (
	"context"
	"fmt"
	"time"

	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/configs"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/index"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/models"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/storage"
)

type LifecycleManager struct {
	cfg     *configs.LifecycleConfig
	minio   *storage.MinIOClient
	recIdx  *index.RecordingIndex
	tierMgr *storage.TierManager
}

func NewLifecycleManager(
	cfg *configs.LifecycleConfig,
	minio *storage.MinIOClient,
	recIdx *index.RecordingIndex,
	tierMgr *storage.TierManager,
) *LifecycleManager {
	return &LifecycleManager{cfg: cfg, minio: minio, recIdx: recIdx, tierMgr: tierMgr}
}

// RunOnce 执行一次生命周期扫描
func (m *LifecycleManager) RunOnce(ctx context.Context) error {
	now := time.Now()

	// 1. hot → warm：status=hot 且 created_at < now-hot_days
	hotCutoff := now.Add(-time.Duration(m.cfg.HotDays) * 24 * time.Hour)
	if err := m.migrateAll(ctx, models.TierHot, models.TierWarm, hotCutoff); err != nil {
		return fmt.Errorf("hot→warm migration: %w", err)
	}

	// 2. warm → cold：status=warm 且 created_at < now-warm_days
	warmCutoff := now.Add(-time.Duration(m.cfg.WarmDays) * 24 * time.Hour)
	if err := m.migrateAll(ctx, models.TierWarm, models.TierCold, warmCutoff); err != nil {
		return fmt.Errorf("warm→cold migration: %w", err)
	}

	// 3. 删除：created_at < now-retention_days
	retentionCutoff := now.Add(-time.Duration(m.cfg.RetentionDays) * 24 * time.Hour)
	if err := m.deleteAll(ctx, retentionCutoff); err != nil {
		return fmt.Errorf("delete expired: %w", err)
	}

	return nil
}

// migrateAll 分页迁移所有符合条件的录制文件
func (m *LifecycleManager) migrateAll(ctx context.Context, fromTier, toTier string, cutoff time.Time) error {
	from := 0
	for {
		recs, err := m.recIdx.ListForMigration(ctx, fromTier, cutoff, from, 100)
		if err != nil {
			return err
		}
		if len(recs) == 0 {
			break
		}
		for _, rec := range recs {
			if err := m.migrateRecording(ctx, rec, toTier); err != nil {
				from++ // 迁移失败时偏移，避免无限重试同一条记录
				continue
			}
		}
		if len(recs) < 100 {
			break
		}
		from += len(recs)
	}
	return nil
}

func (m *LifecycleManager) migrateRecording(ctx context.Context, rec *models.Recording, toTier string) error {
	dstBucket := m.tierMgr.GetBucketByTier(toTier)
	if err := m.minio.MigrateObject(ctx, rec.Bucket, dstBucket, rec.ObjectKey); err != nil {
		return err
	}
	return m.recIdx.UpdateTierStatus(ctx, rec.RecordingID, toTier, dstBucket)
}

func (m *LifecycleManager) deleteAll(ctx context.Context, cutoff time.Time) error {
	from := 0
	for {
		recs, err := m.recIdx.ListForMigration(ctx, models.TierCold, cutoff, from, 100)
		if err != nil {
			return err
		}
		if len(recs) == 0 {
			break
		}
		for _, rec := range recs {
			m.deleteRecording(ctx, rec)
		}
		if len(recs) < 100 {
			break
		}
		from += len(recs)
	}
	return nil
}

func (m *LifecycleManager) deleteRecording(ctx context.Context, rec *models.Recording) {
	m.minio.DeleteObject(ctx, rec.Bucket, rec.ObjectKey)
	m.recIdx.UpdateTierStatus(ctx, rec.RecordingID, "deleted", "")
}
