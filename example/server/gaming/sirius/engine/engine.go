package engine

import "time"

type RankUnit struct {
	ID    uint64
	Key   uint64
	Value []byte
}

type RankEngineConfig struct {
	// 排行榜需要保存排名的数量
	MaxSize uint32
	// 需要缓存的暂时不在排行榜上的冗余节点数量
	RedundantNodeNum uint32
	// 主榜ID, 仅用于快照榜配置
	PrimaryRankID uint32
	// 清空周期
	ClearPeriod TimePeriod
	// 生成快照周期
	SnapshotPeriod TimePeriod
	// 禁止更新数据周期
	// 客户端可以通过ByPassNoUpdate来强制更新数据
	NoUpdatePeriod TimePeriod
}

type RankEngine interface {
	Config() RankEngineConfig
	Size() uint32
	Get(id uint64) (bool, uint32, RankUnit)
	GetByRank(pos uint32) (bool, RankUnit)
	GetRange(pos, num uint32) []RankUnit
	Update(u RankUnit) (bool, uint32, RankUnit)
	Delete(id uint64) (bool, uint32, RankUnit)
	CreateSnapshot() RankEngine
	Clear()
	CopyFrom(rank RankEngine)

	LastClearTime() time.Time
	SetLastClearTime(t time.Time)

	LastSnapshotTime() time.Time
	SetLastSnapshotTime(t time.Time)
}

func NewRankEngine(config RankEngineConfig) RankEngine {
	return NewArrayRankEngine(config)
}
