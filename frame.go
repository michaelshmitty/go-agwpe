package agwpe

import (
	"encoding/binary"
	"fmt"
)

type Header struct {
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

func serializeHeader(h Header) []byte {
	data := make([]byte, 36)

	data[0] = h.Port
	// Ignore reserved bytes 1 to 3
	data[4] = h.DataKind
	// Ignore reserved byte 5
	data[6] = h.PID
	// Ignore reserved byte 7
	copy(data[8:18], h.CallFrom[:])
	copy(data[18:28], h.CallTo[:])
	binary.LittleEndian.PutUint32(data[28:32], h.DataLen)
	copy(data[32:36], h.User[:])

	return data
}

type Frame interface {
	Header() Header
}

type VersionFrame struct {
	header       Header
	major, minor uint16
}

func NewVersionFrame() *VersionFrame {
	return &VersionFrame{
		header: Header{
			Port:     0x01,
			DataKind: 0x52},
	}
}

func (f *VersionFrame) Header() Header {
	return f.header
}

func (f *VersionFrame) Serialize() ([]byte, error) {
	data := serializeHeader(f.header)

	return data, nil
}

func (f *VersionFrame) parseVersion(data []byte) {
	f.major = (uint16(data[37]) << 8) | uint16(data[36])
	f.minor = (uint16(data[41]) << 8) | uint16(data[40])
	fmt.Printf("AGWPE version: %d.%d\n", f.major, f.minor)
}
