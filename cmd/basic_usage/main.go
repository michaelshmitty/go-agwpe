package main

import (
	"fmt"
	"os"
)

import "git.smith.eu/m/go-agwpe"

func main() {
	engine := agwpe.NewPacketEngine("172.16.0.4:8000")

	err := engine.Connect()
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
	}
	defer engine.Disconnect()

	version := engine.Version()
	fmt.Printf("AGWPE version: %s\n", version)

	frame := agwpe.NewVersionFrame()
	frameData, err := frame.Serialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Serialization error: %v\n", err)
	}

	err = engine.Send(frameData)
	if err != nil {
		fmt.Println("Error:", err)
	}

	for msg := range engine.Messages {
		major_verion := (uint16(msg[37]) << 8) | uint16(msg[36])
		minor_verion := (uint16(msg[41]) << 8) | uint16(msg[40])
		fmt.Printf("AGWPE version: %d.%d\n", major_verion, minor_verion)
	}
}
