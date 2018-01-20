package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/jacobwpeng/sirius/frame"
	"github.com/jacobwpeng/sirius/serverproto"
)

const (
	MAX_BUFFERED_JOB_RESULT = 128
	CONN_READ_TIMEOUT       = time.Millisecond * 100
	CONN_WRITE_TIMEOUT      = time.Millisecond * 100
)

type TCPClient struct {
	dispatcher     *Dispatcher
	conn           *net.TCPConn
	wg             sync.WaitGroup
	doneChan       chan struct{}
	errChan        chan error
	jobResultQueue chan JobResult
}

func NewTCPClient(dispatcher *Dispatcher, conn *net.TCPConn) *TCPClient {
	return &TCPClient{
		dispatcher:     dispatcher,
		conn:           conn,
		doneChan:       make(chan struct{}, 2),
		errChan:        make(chan error, 2),
		jobResultQueue: make(chan JobResult, MAX_BUFFERED_JOB_RESULT),
	}
}

func MustMarshal(pb proto.Message) []byte {
	data, err := proto.Marshal(pb)
	if err != nil {
		glog.Fatal(err)
	}
	return data
}

func (c *TCPClient) Run() {
	defer c.conn.Close()
	c.wg.Add(1)
	go c.StartReading()
	c.wg.Add(1)
	go c.StartWriting()
	select {
	case err := <-c.errChan:
		serr, ok := err.(*Error)
		if ok && serr.PrevErr == io.EOF {
			glog.V(2).Infof("TCPClient %s disconnected", c.conn.RemoteAddr())
		} else {
			glog.Warningf("TCPClient %s error: %s", c.conn.RemoteAddr(), err)
		}
		c.Stop()
	case <-c.doneChan:
	}
	c.conn.Close()
	c.wg.Wait()
	glog.V(2).Infof("TCPClient %s done", c.conn.RemoteAddr())
}

func (c *TCPClient) StartReading() {
	defer c.wg.Done()
	for {
		var frame frame.Frame
		br := bufio.NewReader(c.conn)
		if _, err := frame.ReadFrom(br); err != nil {
			select {
			case <-c.doneChan:
				return
			default:
			}
			c.errChan <- NewError("Read frame", err)
			break
		}
		glog.V(2).Infof("New frame from %s", c.conn.RemoteAddr())
		if err := frame.Check(); err != nil {
			c.errChan <- NewError("Check frame", err)
			break
		}
		job, err := c.CreateJob(&frame)
		if err != nil {
			c.errChan <- NewError("Create job", err)
			break
		}
		c.dispatcher.jobQueue <- job
	}
}

func (c *TCPClient) CreateJob(frame *frame.Frame) (job Job, err error) {
	msgType := serverproto.MessageType(frame.PayloadType)
	var msg proto.Message
	switch msgType {
	case serverproto.MessageType_TypeGetRequest:
		msg = &serverproto.GetRequest{}
	case serverproto.MessageType_TypeGetByRankRequest:
		msg = &serverproto.GetByRankRequest{}
	case serverproto.MessageType_TypeGetRangeRequest:
		msg = &serverproto.GetRangeRequest{}
	case serverproto.MessageType_TypeUpdateRequest:
		msg = &serverproto.UpdateRequest{}
	case serverproto.MessageType_TypeDeleteRequest:
		msg = &serverproto.DeleteRequest{}
	default:
		return job, fmt.Errorf("Unexpected type: %d", msgType)
	}
	if err = proto.Unmarshal(frame.Payload, msg); err != nil {
		return job, err
	}
	glog.V(2).Info("New message: ", msg)
	switch m := msg.(type) {
	case *serverproto.GetRequest:
		job.RankID = m.GetRank()
	case *serverproto.GetByRankRequest:
		job.RankID = m.GetRank()
	case *serverproto.GetRangeRequest:
		job.RankID = m.GetRank()
	case *serverproto.UpdateRequest:
		job.RankID = m.GetRank()
	case *serverproto.DeleteRequest:
		job.RankID = m.GetRank()
	default:
		glog.Warning("Unexpected message type")
	}
	job.Frame = frame
	job.Msg = msg
	job.resultChan = c.jobResultQueue
	return job, nil
}

func (c *TCPClient) StartWriting() {
	defer c.wg.Done()
	for {
		select {
		case <-c.doneChan:
			return
		case jobResult := <-c.jobResultQueue:
			var payload []byte
			if jobResult.Msg != nil {
				payload = MustMarshal(jobResult.Msg)
			}
			replyFrame := frame.New(jobResult.FramePayloadType, payload)
			replyFrame.ErrCode = jobResult.ErrCode
			replyFrame.Ctx = jobResult.FrameCtx
			if _, err := replyFrame.WriteTo(c.conn); err != nil {
				select {
				case <-c.doneChan:
					return
				default:
				}
				c.errChan <- NewError("Write frame", err)
				return
			}
			glog.V(2).Infof("Write frame to %s, %v", c.conn.RemoteAddr(), replyFrame)
		}
	}
}

func (c *TCPClient) Stop() {
	select {
	case <-c.doneChan:
		return
	default:
		close(c.doneChan)
	}
}

func (c *TCPClient) StopAndWait() {
	c.Stop()
	c.wg.Wait()
}
