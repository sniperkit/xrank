package server

import "github.com/golang/protobuf/proto"

type JobResult struct {
	FrameCtx         uint64
	FramePayloadType uint32
	ErrCode          int32
	Msg              proto.Message
}
