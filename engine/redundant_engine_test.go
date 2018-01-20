package engine

import (
	"testing"
)

func TestRedundantEngineSize(t *testing.T) {
	u1 := RankUnit{ID: 1024, Key: 10, Value: []byte("Soldier76")}
	u2 := RankUnit{ID: 1025, Key: 12, Value: []byte("McCree")}
	u3 := RankUnit{ID: 1026, Key: 14, Value: []byte("Sombra")}
	//u4 := RankUnit{ID: 1025, Key: 12, Value: []byte("D.Va")}
	e := NewRedundantRankEngine(RankEngineConfig{MaxSize: 2, RedundantNodeNum: 1})
	e.Update(u3)
	e.Update(u2)
	e.Update(u1)

	if e.Size() != 2 {
		t.Errorf("Expect size 2, got: %d", e.Size())
	}

}

func TestRedundantEngineGet(t *testing.T) {
	u1 := RankUnit{ID: 1024, Key: 10, Value: []byte("Soldier76")}
	u2 := RankUnit{ID: 1025, Key: 12, Value: []byte("McCree")}
	u3 := RankUnit{ID: 1026, Key: 14, Value: []byte("Sombra")}
	//u4 := RankUnit{ID: 1025, Key: 12, Value: []byte("D.Va")}
	e := NewRedundantRankEngine(RankEngineConfig{MaxSize: 2, RedundantNodeNum: 1})
	e.Update(u2)
	e.Update(u1)
	e.Update(u3)

	exist, _, _ := e.Get(u1.ID)
	if exist {
		t.Errorf("Expect ID %d not exist", u1.ID)
	}

	exist, pos, u := e.Get(u2.ID)
	if !exist {
		t.Errorf("Expect ID %d exist", u2.ID)
	}
	if pos != 1 {
		t.Errorf("Expect ID %d rank 1, got: %d", u2.ID, pos)
	}
	if err := checkUnitEqual(u, u2); err != nil {
		t.Error(err)
	}

	exist, pos, u = e.Get(u3.ID)
	if !exist {
		t.Errorf("Expect ID %d exist", u3.ID)
	}
	if pos != 0 {
		t.Errorf("Expect ID %d rank 0, got: %d", u3.ID, pos)
	}
	if err := checkUnitEqual(u, u3); err != nil {
		t.Error(err)
	}
}

func TestRedundantEngineGetByRank(t *testing.T) {
	u1 := RankUnit{ID: 1024, Key: 10, Value: []byte("Soldier76")}
	u2 := RankUnit{ID: 1025, Key: 12, Value: []byte("McCree")}
	u3 := RankUnit{ID: 1026, Key: 14, Value: []byte("Sombra")}
	u4 := RankUnit{ID: 1025, Key: 20, Value: []byte("D.Va")}
	e := NewRedundantRankEngine(RankEngineConfig{MaxSize: 2, RedundantNodeNum: 1})
	e.Update(u2)
	e.Update(u1)
	e.Update(u4)
	e.Update(u3)

	exist, u := e.GetByRank(0)
	if !exist {
		t.Errorf("Expect rank 0 exist")
	}
	if err := checkUnitEqual(u, u4); err != nil {
		t.Error(err)
	}

	exist, u = e.GetByRank(1)
	if !exist {
		t.Errorf("Expect rank 1 exist")
	}
	if err := checkUnitEqual(u, u3); err != nil {
		t.Error(err)
	}

	exist, u = e.GetByRank(2)
	if exist {
		t.Error("Expect rank 2 not exist")
	}
}

func TestRedundantEngineGetRange(t *testing.T) {
	u1 := RankUnit{ID: 1024, Key: 10, Value: []byte("Soldier76")}
	u2 := RankUnit{ID: 1025, Key: 12, Value: []byte("McCree")}
	u3 := RankUnit{ID: 1026, Key: 14, Value: []byte("Sombra")}
	u4 := RankUnit{ID: 1025, Key: 20, Value: []byte("D.Va")}
	e := NewRedundantRankEngine(RankEngineConfig{MaxSize: 2, RedundantNodeNum: 1})
	e.Update(u2)
	e.Update(u1)
	e.Update(u4)
	e.Update(u3)

	units := e.GetRange(0, 10)
	if units == nil {
		t.Fatal("Expect non nil result")
	}
	if len(units) != 2 {
		t.Fatalf("Expect 2 units, got: %d", len(units))
	}
	if err := checkUnitEqual(units[0], u4); err != nil {
		t.Error(err)
	}
	if err := checkUnitEqual(units[1], u3); err != nil {
		t.Error(err)
	}

	units = e.GetRange(e.Size(), 1)
	if units != nil {
		t.Fatal("Expect nil result")
	}
}

func TestRedundantEngineUpdate(t *testing.T) {
	u1 := RankUnit{ID: 1024, Key: 10, Value: []byte("Soldier76")}
	u2 := RankUnit{ID: 1025, Key: 12, Value: []byte("McCree")}
	u3 := RankUnit{ID: 1026, Key: 4, Value: []byte("Sombra")}
	u4 := RankUnit{ID: 1026, Key: 20, Value: []byte("Sombra")}
	u5 := RankUnit{ID: 1025, Key: 1, Value: []byte("McCree")}
	e := NewRedundantRankEngine(RankEngineConfig{MaxSize: 2, RedundantNodeNum: 1})

	exist, _, _ := e.Update(u3)
	if exist {
		t.Errorf("Expect ID %d not exist", u3.ID)
	}
	exist, _, _ = e.Update(u1)
	if exist {
		t.Errorf("Expect ID %d not exist", u1.ID)
	}
	exist, _, _ = e.Update(u2)
	if exist {
		t.Errorf("Expect ID %d not exist", u2.ID)
	}
	// 更新之前u3是redundant节点
	exist, _, _ = e.Update(u4)
	if exist {
		t.Errorf("Expect ID %d not exist", u3.ID)
	}

	exist, _, u := e.Get(u4.ID)
	if !exist {
		t.Errorf("Expect ID %d exist", u4.ID)
	}

	exist, _, u = e.Update(u5)
	if !exist {
		t.Errorf("Expect ID %d exist", u5.ID)
	}
	if err := checkUnitEqual(u, u2); err != nil {
		t.Error(err)
	}
}

func TestRedundantEngineDelete(t *testing.T) {
	u1 := RankUnit{ID: 1024, Key: 10, Value: []byte("Soldier76")}
	u2 := RankUnit{ID: 1025, Key: 12, Value: []byte("McCree")}
	u3 := RankUnit{ID: 1026, Key: 4, Value: []byte("Sombra")}
	e := NewRedundantRankEngine(RankEngineConfig{MaxSize: 2, RedundantNodeNum: 1})
	e.Update(u1)
	e.Update(u2)
	e.Update(u3)

	exist, _, _ := e.Delete(u3.ID)
	if exist {
		t.Errorf("Expect ID %d not exist", u3.ID)
	}

	exist, pos, u := e.Delete(u1.ID)
	if !exist {
		t.Errorf("Expect ID %d exist", u1.ID)
	}
	if pos != 1 {
		t.Errorf("Expect ID %d rank 1, got: %d", u1.ID, pos)
	}
	if err := checkUnitEqual(u, u1); err != nil {
		t.Error(err)
	}

	if e.Size() != 1 {
		t.Errorf("Expect size 1, got: %d", e.Size())
	}
}

func TestRedundantEngineCreateSnapshot(t *testing.T) {
	u1 := RankUnit{ID: 1024, Key: 10, Value: []byte("Soldier76")}
	u2 := RankUnit{ID: 1025, Key: 12, Value: []byte("McCree")}
	u3 := RankUnit{ID: 1026, Key: 4, Value: []byte("Sombra")}
	u4 := RankUnit{ID: 1026, Key: 20, Value: []byte("D.Va")}
	u5 := RankUnit{ID: 1027, Key: 30, Value: []byte("Reaper")}
	e := NewRedundantRankEngine(RankEngineConfig{MaxSize: 2, RedundantNodeNum: 1})

	e.Update(u1)
	e.Update(u2)
	e.Update(u4)

	snapshot := e.CreateSnapshot()
	e.Update(u3)
	e.Update(u5)

	units := snapshot.GetRange(0, snapshot.Size())
	if len(units) != int(snapshot.Config().MaxSize) {
		t.Errorf("Expect size %d, got: %d", snapshot.Config().MaxSize, len(units))
	}
	if err := checkUnitEqual(u4, units[0]); err != nil {
		t.Error(err)
	}
	if err := checkUnitEqual(u2, units[1]); err != nil {
		t.Error(err)
	}
}

func TestRedundantEngineKeyDecrease(t *testing.T) {
	u1 := RankUnit{ID: 1024, Key: 10, Value: []byte("Soldier76")}
	u2 := RankUnit{ID: 1025, Key: 12, Value: []byte("McCree")}
	u3 := RankUnit{ID: 1026, Key: 4, Value: []byte("Sombra")}
	u4 := RankUnit{ID: 1026, Key: 20, Value: []byte("D.Va")}
	u5 := RankUnit{ID: 1027, Key: 30, Value: []byte("Reaper")}
	e := NewRedundantRankEngine(RankEngineConfig{MaxSize: 2, RedundantNodeNum: 1})

	// 通过上报一个较大的Key后再上报一个较小的Key来挤出比自己更高排名的值
	e.Update(u1)
	e.Update(u2)
	e.Update(u4)

	if e.Size() != uint32(2) {
		t.Fatalf("Expect size 2, got: %d", e.Size())
	}

	exist, _, _ := e.Get(u2.ID)
	if !exist {
		t.Errorf("Expect ID %d exist", u2.ID)
	}
	exist, _, _ = e.Get(u4.ID)
	if !exist {
		t.Errorf("Expect ID %d exist", u4.ID)
	}
	exist, _, _ = e.Get(u1.ID)
	if exist {
		t.Errorf("Expect ID %d not exist", u1.ID)
	}

	e.Update(u3)
	exist, _, _ = e.Get(u1.ID)
	if !exist {
		t.Errorf("Expect ID %d exist", u1.ID)
	}

	e.Update(u5)
	// 这时候u1被隐藏，u3的key小于u1，RedundantNodeNum == 1，u3被淘汰
	exist, _, _ = e.Get(u1.ID)
	if exist {
		t.Errorf("Expect ID %d not exist", u1.ID)
	}
	exist, _, _ = e.Get(u3.ID)
	if exist {
		t.Errorf("Expect ID %d not exist", u1.ID)
	}
}

func TestRedundantEngineClear(t *testing.T) {
	e := NewRedundantRankEngine(RankEngineConfig{MaxSize: 2, RedundantNodeNum: 1})
	u1 := RankUnit{ID: 1024, Key: 10, Value: []byte("Soldier76")}
	u2 := RankUnit{ID: 1025, Key: 12, Value: []byte("McCree")}
	u3 := RankUnit{ID: 1026, Key: 14, Value: []byte("Sombra")}

	e.Update(u1)
	e.Update(u2)
	e.Update(u3)

	e.Clear()
	if e.Size() != 0 {
		t.Errorf("Expect empty rank, got: %d", e.Size())
	}
}
