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
	address  string
	cfg      Config
	Ready    bool
	conn     net.Conn
	mu       sync.Mutex
	Messages chan []byte
	version  string
}

func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.Timeout = d
	}
}

func NewPacketEngine(address string, opts ...Option) *PacketEngine {
	defaultCfg := Config{
		Timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(&defaultCfg)
	}

	pe := PacketEngine{
		address:  address,
		cfg:      defaultCfg,
		Ready:    false,
		Messages: make(chan []byte),
	}

	return &pe
}

func (pe *PacketEngine) Connect() error {
	conn, err := net.DialTimeout("tcp", pe.address, pe.cfg.Timeout)
	if err != nil {
		return err
	}
	pe.conn = conn
	fmt.Println("Connected to", pe.address)
	go pe.receiveLoop()

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

func (pe *PacketEngine) Send(data []byte) error {
	_, err := pe.conn.Write(data)

	return err
}

func (pe *PacketEngine) receiveLoop() {
	defer close(pe.Messages)

	for {
		buf := make([]byte, 1024)
		n, err := pe.conn.Read(buf)
		if err != nil {
			// FIXME: handle error
			panic(err)
		}

		data := buf[:n]

		pe.Messages <- data
	}
}

func (pe *PacketEngine) Version() string {
	if pe.conn == nil {
		return ""
	}
	if pe.version != "" {
		return pe.version
	}

	return fmt.Sprintf("%d.%d")
}
