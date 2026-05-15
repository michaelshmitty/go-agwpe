package main

import "encoding/binary"

type Frame struct {
	Port      byte
	Reserved0 [3]byte
	DataKind  byte
	Reserved1 byte
	PID       byte
	Reserved2 byte
	CallFrom  [10]byte
	CallTo    [10]byte
	DataLen   uint32
	User      [4]byte
}

func (f *Frame) Serialize() []byte {
	buf := make([]byte, 36)

	buf[0] = f.Port
	// Ignore reserved bytes 1 to 3
	buf[4] = f.DataKind
	// Ignore reserved byte 5
	buf[6] = f.PID
	// Ignore reserved byte 7
	copy(buf[8:18], f.CallFrom[:])
	copy(buf[18:28], f.CallTo[:])
	binary.BigEndian.PutUint32(buf[28:32], f.DataLen)
	copy(buf[32:36], f.User[:])

	return buf
}
