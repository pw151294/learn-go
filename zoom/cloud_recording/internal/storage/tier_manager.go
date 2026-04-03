package storage

import (
	"time"

	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/configs"
)

type TierManager struct {
	cfg *configs.MinIOConfig
	lc  *configs.LifecycleConfig
}

func NewTierManager(minioCfg *configs.MinIOConfig, lcCfg *configs.LifecycleConfig) *TierManager {
	return &TierManager{cfg: minioCfg, lc: lcCfg}
}

// DetermineTier returns the storage tier for an object based on its age.
// 7 days -> hot; 7-30 days -> warm; 30+ days -> cold.
func (t *TierManager) DetermineTier(createdAt time.Time) string {
	age := time.Since(createdAt)
	if age < time.Duration(t.lc.HotDays)*24*time.Hour {
		return "hot"
	}
	if age < time.Duration(t.lc.WarmDays)*24*time.Hour {
		return "warm"
	}
	return "cold"
}

// GetBucketByTier returns the bucket name for the given tier.
func (t *TierManager) GetBucketByTier(tier string) string {
	switch tier {
	case "warm":
		return t.cfg.WarmBucket
	case "cold":
		return t.cfg.ColdBucket
	default:
		return t.cfg.HotBucket
	}
}
