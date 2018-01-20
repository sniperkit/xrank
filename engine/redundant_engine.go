package engine

import "time"

type RedundantRankEngine struct {
	config     RankEngineConfig
	underlying RankEngine
}

func NewRedundantRankEngine(config RankEngineConfig) *RedundantRankEngine {
	underlyingConfig := config
	underlyingConfig.MaxSize = config.MaxSize + config.RedundantNodeNum
	return &RedundantRankEngine{
		config:     config,
		underlying: NewArrayRankEngine(underlyingConfig),
	}
}

func (e *RedundantRankEngine) Config() RankEngineConfig {
	return e.config
}

func (e *RedundantRankEngine) Size() uint32 {
	if e.underlying.Size() > e.config.MaxSize {
		return e.config.MaxSize
	}
	return e.underlying.Size()
}

func (e *RedundantRankEngine) Get(id uint64) (bool, uint32, RankUnit) {
	exist, pos, u := e.underlying.Get(id)
	if pos >= e.config.MaxSize {
		return false, 0, RankUnit{}
	}
	return exist, pos, u
}

func (e *RedundantRankEngine) GetByRank(pos uint32) (bool, RankUnit) {
	if pos >= e.config.MaxSize {
		return false, RankUnit{}
	}
	return e.underlying.GetByRank(pos)
}

func (e *RedundantRankEngine) GetRange(pos, num uint32) []RankUnit {
	if pos >= e.Size() {
		return nil
	}
	n := num
	if pos+num >= e.Size() {
		n = e.Size() - pos
	}
	return e.underlying.GetRange(pos, n)
}

func (e *RedundantRankEngine) Update(u RankUnit) (bool, uint32, RankUnit) {
	exist, pos, u := e.underlying.Update(u)
	if pos >= e.config.MaxSize {
		return false, 0, RankUnit{}
	}
	return exist, pos, u
}

func (e *RedundantRankEngine) Delete(id uint64) (bool, uint32, RankUnit) {
	exist, pos, u := e.underlying.Delete(id)
	if pos >= e.config.MaxSize {
		return false, 0, RankUnit{}
	}
	return exist, pos, u
}

func (e *RedundantRankEngine) CreateSnapshot() RankEngine {
	snapshot := &RedundantRankEngine{
		config:     e.config,
		underlying: e.underlying.CreateSnapshot(),
	}
	return snapshot
}

func (e *RedundantRankEngine) Clear() {
	e.underlying.Clear()
}

func (e *RedundantRankEngine) CopyFrom(rank RankEngine) {
	e.underlying.CopyFrom(rank)
}

func (e *RedundantRankEngine) LastClearTime() time.Time {
	return e.underlying.LastClearTime()
}

func (e *RedundantRankEngine) SetLastClearTime(t time.Time) {
	e.underlying.SetLastClearTime(t)
}

func (e *RedundantRankEngine) LastSnapshotTime() time.Time {
	return e.underlying.LastSnapshotTime()
}

func (e *RedundantRankEngine) SetLastSnapshotTime(t time.Time) {
	e.underlying.SetLastSnapshotTime(t)
}
