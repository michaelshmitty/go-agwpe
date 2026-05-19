package agwpe

import (
	"testing"
	"time"
)

func TestNewPacketEngine(t *testing.T) {
	tests := []struct {
		name            string
		address         string
		options         []Option
		expectedError   bool
		expectedTimeout time.Duration
	}{
		{
			name:            "Default configuration",
			address:         "localhost:1234",
			options:         []Option{},
			expectedError:   false,
			expectedTimeout: 5 * time.Second,
		},
		{
			name:            "Custom timeout",
			address:         "localhost:1234",
			options:         []Option{WithTimeout(30 * time.Second)},
			expectedError:   false,
			expectedTimeout: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewPacketEngine(tt.address, tt.options...)

			if (err != nil) != tt.expectedError {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if !tt.expectedError {
				if engine == nil {
					t.Fatal("Expected non-nil engine")
				}
				if engine.Ready {
					t.Fatal("Expected engine that is not ready")
				}
				if engine.cfg.Timeout != tt.expectedTimeout {
					t.Errorf("Expected timeout %v, got %v", tt.expectedTimeout, engine.cfg.Timeout)
				}
			}

		})
	}
}

func TestConnect(t *testing.T) {
	tests := []struct {
		name          string
		address       string
		options       []Option
		expectedError bool
	}{
		{
			name:          "Successful connection",
			address:       "localhost:1234",
			expectedError: false,
		},
		{
			name:          "Invalid address",
			address:       "invalid-address",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, _ := NewPacketEngine(tt.address, tt.options...)
			err := engine.Connect()

			if tt.expectedError && (err == nil) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}

			if !tt.expectedError && (err != nil) {
				t.Errorf("Did not expect error: %v", err)
			}
		})
	}
}

func TestDisconnect(t *testing.T) {
	engine, _ := NewPacketEngine("localhost:1234")
	engine.Connect()

	err := engine.Disconnect()
	if err != nil {
		t.Errorf("Did not expect error: %v", err)
	}

	// Call to Disconnect() on disconnected packet engine should return an error
	err = engine.Disconnect()
	if err == nil {
		t.Error("Expected an error, but got", nil)
	}
}
