package agwpe

import (
	"testing"
	"time"
)

// NOTE(m): Use more table-driven-tests? At the moment this feels less abstract.

func TestNewPacketEngine(t *testing.T) {
	// Default configuration
	engine, err := NewPacketEngine("localhost:1234")
	if err != nil {
		t.Errorf("Did not expect error: %v", err)
	}
	if engine.cfg.Timeout != 5*time.Second {
		t.Errorf("Expected 5 second timeout, got %v", engine.cfg.Timeout)
	}
	if engine.Ready {
		t.Error("Expected engine that is not ready")
	}

	// Custom timeout
	engine, err = NewPacketEngine("localhost:1234", WithTimeout(30*time.Second))
	if err != nil {
		t.Errorf("Did not expect error: %v", err)
	}
	if engine.cfg.Timeout != 30*time.Second {
		t.Errorf("Expected 30 second timeout, got %v", engine.cfg.Timeout)
	}
}

func TestNewPacketEngine_Connect(t *testing.T) {
	// Successful connect
	engine, _ := NewPacketEngine("localhost:1234")
	err := engine.Connect()
	if err != nil {
		t.Errorf("Did not expect error: %v", err)
	}

	// Invalid address
	engine, _ = NewPacketEngine("invalid-address")
	err = engine.Connect()
	if err == nil {
		t.Error("Expected an error")
	}
}

func TestNewPacketEngine_Disconnect(t *testing.T) {
	engine, _ := NewPacketEngine("localhost:1234")
	engine.Connect()

	// Disconnect connected engine
	err := engine.Disconnect()
	if err != nil {
		t.Errorf("Did not expect error: %v", err)
	}

	// Call to Disconnect() on disconnected packet engine should return an error
	err = engine.Disconnect()
	if err == nil {
		t.Error("Expected an error")
	}
}
