package main

import (
	"flag"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/jacobwpeng/sirius/frame"
	"github.com/jacobwpeng/sirius/serverproto"
)

func MustMarshal(msg proto.Message) (data []byte) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func NewFrame(msg proto.Message) (*frame.Frame, error) {
	data := MustMarshal(msg)
	var msgType uint32
	switch msg.(type) {
	case *serverproto.UpdateRequest:
		msgType = uint32(serverproto.MessageType_TypeUpdateRequest)
	}
	return frame.New(msgType, data), nil
}

func ce(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	raddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9427")
	ce(err)
	conn, err := net.DialTCP("tcp", nil, raddr)
	ce(err)

	req := &serverproto.UpdateRequest{
		Rank: proto.Uint32(1),
		Data: &serverproto.RankUnit{
			Id:  proto.Uint64(2191195),
			Key: proto.Uint64(1024),
		},
		Reply: proto.Bool(true),
	}

	f, err := NewFrame(req)
	ce(err)
	_, err = f.WriteTo(conn)
	ce(err)
	var reply frame.Frame
	_, err = reply.ReadFrom(conn)
	ce(err)
	log.Print(reply.ErrCode)
}
