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
