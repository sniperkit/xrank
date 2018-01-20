package server

import (
	"fmt"
	"sync"

	"github.com/golang/glog"
	"github.com/jacobwpeng/sirius/engine"
)

type Dispatcher struct {
	wg             sync.WaitGroup
	doneChan       chan struct{}
	rankHandlers   []*RankHandler
	mappedHandlers map[uint32]*RankHandler
	jobQueue       chan Job
}

func NewDispatcher(ranks map[uint32]engine.RankEngine) (*Dispatcher, error) {
	rankHandlers := make([]*RankHandler, 0)
	mappedHandlers := make(map[uint32]*RankHandler)
	for rankID, rank := range ranks {
		primaryRankID := rank.Config().PrimaryRankID
		var rankHandler *RankHandler
		if primaryRankID != 0 {
			continue
		}
		rankHandler = NewRankHandler(rankID, rank)
		rankHandlers = append(rankHandlers, rankHandler)
		mappedHandlers[rankID] = rankHandler
		glog.Infof("New rank handler for rank %d", rankID)
	}

	for rankID, rank := range ranks {
		primaryRankID := rank.Config().PrimaryRankID
		if primaryRankID == 0 {
			continue
		}
		// 处理所有的非Primary Rank
		rankHandler, _ := mappedHandlers[primaryRankID]
		if rankHandler == nil {
			return nil, fmt.Errorf("Primary rank %d not found", primaryRankID)
		}
		if err := rankHandler.AddSnapshotRank(rankID, rank); err != nil {
			return nil, err
		}
		mappedHandlers[rankID] = rankHandler
		glog.Infof("Add snapshot rank %d to %d", rankID, primaryRankID)
	}

	return &Dispatcher{
		doneChan:       make(chan struct{}),
		rankHandlers:   rankHandlers,
		mappedHandlers: mappedHandlers,
		jobQueue:       make(chan Job, MAX_BUFFERED_JOB),
	}, nil
}

func (d *Dispatcher) AddRank(rankID uint32, rank engine.RankEngine) bool {
	_, exist := d.mappedHandlers[rankID]
	if exist {
		return false
	}
	d.mappedHandlers[rankID] = NewRankHandler(rankID, rank)
	return true
}

func (d *Dispatcher) Start() {
	for _, handler := range d.rankHandlers {
		handler.Start(&d.wg)
	}
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for {
			select {
			case <-d.doneChan:
				glog.V(2).Info("Dispatcher exit")
				return
			case job := <-d.jobQueue:
				glog.V(2).Info("New job in dispatcher")
				rankHandler, exist := d.mappedHandlers[job.RankID]
				if !exist {
					glog.Infof("Rank %d not exist", job.RankID)
					job.resultChan <- JobResult{
						FrameCtx: job.Frame.Ctx,
						ErrCode:  ErrRankNotFound,
					}
					continue
				}
				rankHandler.jobQueue <- job
			}
		}
	}()
}

func (d *Dispatcher) Stop() {
	close(d.doneChan)
	for _, handler := range d.rankHandlers {
		handler.Stop()
	}
	d.wg.Wait()
}
