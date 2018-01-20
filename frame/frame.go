package frame

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/jacobwpeng/goutil"
)

const (
	FRAME_MAGIC      = 0x19910926
	MAX_PAYLOAD_SIZE = 60000
)

type Frame struct {
	Magic       uint32
	Ctx         uint64
	ErrCode     int32
	PayloadType uint32
	PayloadSize uint32
	Payload     []byte
}

func New(payloadType uint32, payload []byte) *Frame {
	if len(payload) > MAX_PAYLOAD_SIZE {
		return nil
	}
	return &Frame{
		Magic:       FRAME_MAGIC,
		Ctx:         0,
		PayloadType: payloadType,
		PayloadSize: uint32(len(payload)),
		Payload:     payload,
	}
}

func (frame *Frame) WriteTo(w io.Writer) (n int64, err error) {
	cw := goutil.NewCountWriter(w)
	sw := goutil.NewStrickyWriter(cw)
	binary.Write(sw, binary.LittleEndian, frame.Magic)
	binary.Write(sw, binary.LittleEndian, frame.Ctx)
	binary.Write(sw, binary.LittleEndian, frame.ErrCode)
	binary.Write(sw, binary.LittleEndian, frame.PayloadType)
	binary.Write(sw, binary.LittleEndian, frame.PayloadSize)
	sw.Write(frame.Payload)
	return cw.Count(), sw.Err
}

func (frame *Frame) ReadFrom(r io.Reader) (n int64, err error) {
	cr := goutil.NewCountReader(r)
	sr := goutil.NewStrickyReader(cr)
	binary.Read(sr, binary.LittleEndian, &frame.Magic)
	binary.Read(sr, binary.LittleEndian, &frame.Ctx)
	binary.Read(sr, binary.LittleEndian, &frame.ErrCode)
	binary.Read(sr, binary.LittleEndian, &frame.PayloadType)
	binary.Read(sr, binary.LittleEndian, &frame.PayloadSize)
	if sr.Err != nil {
		return cr.Count(), sr.Err
	}
	if err := frame.CheckHeader(); err != nil {
		return cr.Count(), err
	}
	frame.Payload = make([]byte, frame.PayloadSize)
	sr.Read(frame.Payload)
	return cr.Count(), sr.Err
}

func (frame *Frame) CheckHeader() error {
	if frame.Magic != FRAME_MAGIC {
		return fmt.Errorf("Expect magic 0x%X, got: 0x%X", FRAME_MAGIC, frame.Magic)
	}
	if frame.PayloadSize > MAX_PAYLOAD_SIZE {
		return fmt.Errorf("Max payload size %d, got: %d",
			MAX_PAYLOAD_SIZE, frame.PayloadSize)
	}
	return nil
}

func (frame *Frame) Check() error {
	err := frame.CheckHeader()
	if err != nil {
		return err
	}
	if len(frame.Payload) != int(frame.PayloadSize) {
		return fmt.Errorf("Expect payload size: %d, got: %d",
			frame.PayloadSize, len(frame.Payload))
	}
	return nil
}
