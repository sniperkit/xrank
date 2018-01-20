package engine

import (
	"bytes"
	"sort"
	"time"
)

type ArrayRankUnit RankUnit
type ArrayRankUnitSlice []ArrayRankUnit

func (es ArrayRankUnitSlice) Len() int {
	return len(es)
}

func (es ArrayRankUnitSlice) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es ArrayRankUnitSlice) Less(i, j int) bool {
	return es[i].Key > es[j].Key
}

type ArrayRankEngine struct {
	config           RankEngineConfig
	data             ArrayRankUnitSlice
	lastClearTime    time.Time
	lastSnapshotTime time.Time
}

func NewArrayRankEngine(config RankEngineConfig) *ArrayRankEngine {
	return &ArrayRankEngine{config: config}
}

func (e *ArrayRankEngine) Config() RankEngineConfig {
	return e.config
}

func (e *ArrayRankEngine) Size() uint32 {
	return uint32(e.data.Len())
}

func (e *ArrayRankEngine) Get(id uint64) (bool, uint32, RankUnit) {
	for index := 0; index < e.data.Len(); index++ {
		if e.data[index].ID == id {
			return true, uint32(index), RankUnit(e.data[index])
		}
	}

	return false, 0, RankUnit{}
}

func (e *ArrayRankEngine) GetByRank(pos uint32) (bool, RankUnit) {
	if pos >= e.Size() {
		return false, RankUnit{}
	}
	return true, RankUnit(e.data[pos])
}

func (e *ArrayRankEngine) GetRange(pos, num uint32) []RankUnit {
	if pos >= e.Size() {
		return nil
	}
	n := num
	if pos+num >= e.Size() {
		n = e.Size() - pos
	}
	result := make([]RankUnit, n)
	for i := uint32(0); i < n; i++ {
		result[i] = RankUnit(e.data[pos+i])
	}
	return result
}

func (e *ArrayRankEngine) Update(u RankUnit) (bool, uint32, RankUnit) {
	aru := ArrayRankUnit(u)
	exist, index, old := e.Get(u.ID)
	if exist {
		e.data[index] = aru
	} else if e.config.MaxSize != 0 && e.Size() >= e.config.MaxSize {
		// 已经达到最大上限, 淘汰最后一个
		last := e.Size() - 1
		e.data[last] = aru
	} else {
		e.data = append(e.data, aru)
	}
	sort.Stable(e.data)
	return exist, index, old
}

func (e *ArrayRankEngine) Delete(id uint64) (bool, uint32, RankUnit) {
	exist, pos, u := e.Get(id)
	if exist {
		e.data = append(e.data[:pos], e.data[pos+1:]...)
	}
	return exist, pos, u
}

func (e *ArrayRankEngine) CreateSnapshot() RankEngine {
	snapshot := &ArrayRankEngine{
		config: e.config,
		data:   make(ArrayRankUnitSlice, len(e.data)),
	}
	for i := 0; i < len(e.data); i++ {
		buffer := bytes.NewBuffer(e.data[i].Value)
		snapshot.data[i] = ArrayRankUnit{
			ID:    e.data[i].ID,
			Key:   e.data[i].Key,
			Value: buffer.Bytes(),
		}
	}
	return snapshot
}

func (e *ArrayRankEngine) Clear() {
	e.data = nil
}

func (e *ArrayRankEngine) CopyFrom(rank RankEngine) {
	units := rank.GetRange(0, rank.Size())
	e.Clear()
	e.SetLastClearTime(rank.LastClearTime())
	e.SetLastSnapshotTime(rank.LastSnapshotTime())
	for i := 0; i < len(units); i++ {
		e.Update(units[i])
	}
}

func (e *ArrayRankEngine) LastClearTime() time.Time {
	return e.lastClearTime
}

func (e *ArrayRankEngine) SetLastClearTime(t time.Time) {
	e.lastClearTime = t
}

func (e *ArrayRankEngine) LastSnapshotTime() time.Time {
	return e.lastSnapshotTime
}

func (e *ArrayRankEngine) SetLastSnapshotTime(t time.Time) {
	e.lastSnapshotTime = t
}
