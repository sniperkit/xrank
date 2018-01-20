package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/jacobwpeng/sirius/engine"
	"github.com/jacobwpeng/sirius/serverproto"
)

const (
	MAX_BUFFERED_JOB = 128
)

type RankHandler struct {
	primaryRankID uint32
	primaryRank   engine.RankEngine
	snapshotRanks map[uint32]engine.RankEngine
	done          chan struct{}
	jobQueue      chan Job
}

func NewRankHandler(rankID uint32, rank engine.RankEngine) *RankHandler {
	return &RankHandler{
		primaryRankID: rankID,
		primaryRank:   rank,
		done:          make(chan struct{}),
		jobQueue:      make(chan Job, MAX_BUFFERED_JOB),
		snapshotRanks: make(map[uint32]engine.RankEngine),
	}
}

func (h *RankHandler) AddSnapshotRank(rankID uint32,
	rank engine.RankEngine) error {
	if rank.Config().PrimaryRankID != h.primaryRankID {
		return fmt.Errorf("Expect primary rank id %d, got: %d",
			h.primaryRankID, rank.Config().PrimaryRankID)
	}
	_, exist := h.snapshotRanks[rankID]
	if exist {
		return fmt.Errorf("Snapshot rank %d already exist", rankID)
	}
	h.snapshotRanks[rankID] = rank
	return nil
}

func RankUnitToProto(u engine.RankUnit) *serverproto.RankUnit {
	return &serverproto.RankUnit{
		Id:  proto.Uint64(u.ID),
		Key: proto.Uint64(u.Key),
	}
}

func RankUnitFromProto(u *serverproto.RankUnit) engine.RankUnit {
	return engine.RankUnit{
		ID:    u.GetId(),
		Key:   u.GetKey(),
		Value: u.GetValue(),
	}
}

func (h *RankHandler) Start(wg *sync.WaitGroup) {
	const CRON_CHECK_INTERVAL time.Duration = time.Millisecond * 500
	c := time.NewTicker(CRON_CHECK_INTERVAL).C
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case job := <-h.jobQueue:
				h.HandleJob(job)
			case <-h.done:
				glog.Infof("RankHandler %d exit", h.primaryRankID)
				return
			case now := <-c:
				h.CronCheckAllRanks(now)
			}
		}
	}()
}

func (h *RankHandler) Stop() {
	close(h.done)
}

func (h *RankHandler) CronCheckAllRanks(now time.Time) {
	for rankID, rank := range h.snapshotRanks {
		h.MaybeClearRank(rankID, rank, now)
		h.MaybeSnapshotRank(rankID, rank, now)
	}
	h.MaybeClearRank(h.primaryRankID, h.primaryRank, now)
}

func (h *RankHandler) FindRank(rankID uint32) engine.RankEngine {
	if rankID == h.primaryRankID {
		return h.primaryRank
	}
	rank, _ := h.snapshotRanks[rankID]
	return rank
}

func (h *RankHandler) HandleJob(job Job) {
	glog.V(2).Infof("Rank: %d, Ctx: %d", job.RankID, job.Frame.Ctx)
	rank := h.FindRank(job.RankID)
	if rank == nil {
		glog.Fatalf("Rank %d not found!")
	}
	now := time.Now()
	if job.RankID == h.primaryRankID {
		h.MaybeSnapshotPrimaryRank(now)
	} else {
		h.MaybeSnapshotRank(job.RankID, rank, now)
	}
	h.MaybeClearRank(job.RankID, rank, now)
	var jobResult JobResult
	switch msg := job.Msg.(type) {
	case *serverproto.GetRequest:
		jobResult = h.HandleGet(job, rank, msg)
	case *serverproto.GetByRankRequest:
		jobResult = h.HandleGetByRank(job, rank, msg)
	case *serverproto.GetRangeRequest:
		jobResult = h.HandleGetRange(job, rank, msg)
	case *serverproto.UpdateRequest:
		jobResult = h.HandleUpdate(job, rank, msg, now)
		if !msg.GetReply() {
			return
		}
	case *serverproto.DeleteRequest:
		jobResult = h.HandleDelete(job, rank, msg)
		if !msg.GetReply() {
			return
		}
	default:
		glog.Infof("Unexpected msg type %d", job.Frame.PayloadType)
	}
	glog.V(2).Infof("Write job result, FrameCtx: %d", job.Frame.Ctx)
	job.resultChan <- jobResult
}

func (h *RankHandler) MaybeSnapshotPrimaryRank(now time.Time) {
	for rankID, rank := range h.snapshotRanks {
		h.MaybeSnapshotRank(rankID, rank, now)
	}
}

func (h *RankHandler) MaybeSnapshotRank(rankID uint32, rank engine.RankEngine,
	now time.Time) bool {
	if rank.Config().SnapshotPeriod.Empty() {
		return false
	}
	lastTime := rank.LastSnapshotTime()
	nextTime := rank.Config().SnapshotPeriod.NextTime(lastTime)
	if now.Before(nextTime) {
		return false
	}
	rank.CopyFrom(h.primaryRank)
	rank.SetLastSnapshotTime(now)
	glog.Infof("Snapshot primary rank %d to rank %d", h.primaryRankID, rankID)
	return true
}

func (h *RankHandler) MaybeClearRank(rankID uint32, rank engine.RankEngine,
	now time.Time) bool {
	if rank.Config().ClearPeriod.Empty() {
		return false
	}
	lastTime := rank.LastClearTime()
	nextTime := rank.Config().ClearPeriod.NextTime(lastTime)
	if now.Before(nextTime) {
		return false
	}
	rank.Clear()
	rank.SetLastClearTime(now)
	glog.Infof("Clear rank %d", rankID)
	return true
}

func (h *RankHandler) HandleGet(job Job, rank engine.RankEngine,
	msg *serverproto.GetRequest) JobResult {
	_, pos, value := rank.Get(msg.GetId())

	resp := &serverproto.GetResponse{
		Rank: proto.Uint32(job.RankID),
		Pos:  proto.Uint32(pos),
		Data: RankUnitToProto(value),
	}
	return JobResult{
		FrameCtx:         job.Frame.Ctx,
		FramePayloadType: uint32(serverproto.MessageType_TypeGetResponse),
		Msg:              resp,
	}
}

func (h *RankHandler) HandleGetByRank(job Job, rank engine.RankEngine,
	msg *serverproto.GetByRankRequest) JobResult {
	exist, value := rank.GetByRank(msg.GetPos())
	var pos uint32
	if exist {
		pos = msg.GetPos()
	}

	resp := &serverproto.GetByRankResponse{
		Rank: proto.Uint32(job.RankID),
		Pos:  proto.Uint32(pos),
		Data: RankUnitToProto(value),
	}
	return JobResult{
		FrameCtx:         job.Frame.Ctx,
		FramePayloadType: uint32(serverproto.MessageType_TypeGetByRankResponse),
		Msg:              resp,
	}
}

func (h *RankHandler) HandleGetRange(job Job, rank engine.RankEngine,
	msg *serverproto.GetRangeRequest) JobResult {
	values := rank.GetRange(msg.GetStart(), msg.GetNum())
	data := make([]*serverproto.RankUnit, len(values))
	for i, u := range values {
		data[i] = RankUnitToProto(u)
	}

	resp := &serverproto.GetRangeResponse{
		Rank:  proto.Uint32(job.RankID),
		Total: proto.Uint32(rank.Size()),
		Data:  data,
	}
	return JobResult{
		FrameCtx:         job.Frame.Ctx,
		FramePayloadType: uint32(serverproto.MessageType_TypeGetRangeResponse),
		Msg:              resp,
	}
}

func (h *RankHandler) HandleUpdate(job Job, rank engine.RankEngine,
	msg *serverproto.UpdateRequest, now time.Time) (res JobResult) {

	ts := now.Unix()
	begin := msg.ServerTimeRange.GetBegin()
	end := msg.ServerTimeRange.GetEnd()
	if (begin != 0 || end != 0) && (ts < begin || ts >= end) {
		glog.Infof("Drop update request: expect time range [%d, %d), now %d",
			begin, end, ts)
		return JobResult{
			FrameCtx: job.Frame.Ctx,
			ErrCode:  ErrServerTimeRange,
		}
	}

	if !msg.GetBypassNoUpdate() &&
		!rank.Config().NoUpdatePeriod.Empty() &&
		rank.Config().NoUpdatePeriod.Contains(now) {
		glog.Infof("Drop update request: no update time period, now %d",
			ts)
		return JobResult{
			FrameCtx: job.Frame.Ctx,
			ErrCode:  ErrNoUpdateTimePeriod,
		}
	}

	_, lastPos, lastData := rank.Update(RankUnitFromProto(msg.Data))
	_, pos, _ := rank.Get(msg.Data.GetId())
	if !msg.GetReply() {
		return res
	}
	resp := &serverproto.UpdateResponse{
		Rank:    proto.Uint32(job.RankID),
		LastPos: proto.Uint32(lastPos),
		Pos:     proto.Uint32(pos),
	}
	if msg.GetLastData() {
		resp.Data = RankUnitToProto(lastData)
	}
	return JobResult{
		FrameCtx:         job.Frame.Ctx,
		FramePayloadType: uint32(serverproto.MessageType_TypeUpdateResponse),
		Msg:              resp,
	}
}

func (h *RankHandler) HandleDelete(job Job, rank engine.RankEngine,
	msg *serverproto.DeleteRequest) (res JobResult) {

	_, lastPos, lastData := rank.Delete(msg.GetId())
	if !msg.GetReply() {
		return res
	}
	resp := &serverproto.DeleteResponse{
		Rank:    proto.Uint32(job.RankID),
		LastPos: proto.Uint32(lastPos),
	}
	if msg.GetLastData() {
		resp.Data = RankUnitToProto(lastData)
	}
	return JobResult{
		FrameCtx:         job.Frame.Ctx,
		FramePayloadType: uint32(serverproto.MessageType_TypeDeleteResponse),
		Msg:              resp,
	}
}
