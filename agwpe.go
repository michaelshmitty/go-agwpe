package agwpe

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Config struct {
	Timeout time.Duration
}

type Option func(*Config)

type PacketEngine struct {
	address string
	cfg     Config
	Ready   bool
	conn    net.Conn
	mu      sync.Mutex
}

func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.Timeout = d
	}
}

func NewPacketEngine(address string, opts ...Option) (*PacketEngine, error) {
	defaultCfg := Config{
		Timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(&defaultCfg)
	}

	return &PacketEngine{
		address: address,
		cfg:     defaultCfg,
		Ready:   false,
	}, nil
}

func (pe *PacketEngine) Connect() error {
	conn, err := net.DialTimeout("tcp", pe.address, pe.cfg.Timeout)
	if err != nil {
		return err
	}
	pe.conn = conn
	fmt.Println("Connected to", pe.address)

	return nil
}

func (pe *PacketEngine) Disconnect() error {
	if pe.conn == nil {
		return fmt.Errorf("Not connected!")
	}
	pe.conn.Close()
	pe.conn = nil
	fmt.Println("Disonnected from", pe.address)

	return nil
}

// import "encoding/binary"

// type Frame struct {
// 	Port      byte
// 	Reserved0 [3]byte
// 	DataKind  byte
// 	Reserved1 byte
// 	PID       byte
// 	Reserved2 byte
// 	CallFrom  [10]byte
// 	CallTo    [10]byte
// 	DataLen   uint32
// 	User      [4]byte
// }

// func (f *Frame) Serialize() []byte {
// 	buf := make([]byte, 36)

// 	buf[0] = f.Port
// 	// Ignore reserved bytes 1 to 3
// 	buf[4] = f.DataKind
// 	// Ignore reserved byte 5
// 	buf[6] = f.PID
// 	// Ignore reserved byte 7
// 	copy(buf[8:18], f.CallFrom[:])
// 	copy(buf[18:28], f.CallTo[:])
// 	binary.LittleEndian.PutUint32(buf[28:32], f.DataLen)
// 	copy(buf[32:36], f.User[:])

// 	return buf
// }
