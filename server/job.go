package server

import (
	"github.com/golang/protobuf/proto"
	"github.com/jacobwpeng/sirius/frame"
)

type Job struct {
	Frame      *frame.Frame
	RankID     uint32
	Msg        proto.Message
	resultChan chan<- JobResult
}
