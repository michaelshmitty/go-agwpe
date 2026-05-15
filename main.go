package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func main() {
	address := "172.16.0.4:8000"
	timeout := 5 * time.Second
	send_buf := new(bytes.Buffer)
	recv_buf := make([]byte, 1024)

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", address)

	rFrame := Frame{
		Port:     0x01,
		DataKind: 0x52,
	}

	err = binary.Write(send_buf, binary.LittleEndian, rFrame)
	if err != nil {
		fmt.Println("Binary write error:", err)
	}

	byteSlice := send_buf.Bytes()

	_, err = conn.Write(byteSlice)
	if err != nil {
		fmt.Println("Error sending frame to server:", err)
		return
	}

	n, err := conn.Read(recv_buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	fmt.Printf("Received %d bytes:\n", n)

	reader := bytes.NewReader(recv_buf)
	var f Frame
	err = binary.Read(reader, binary.LittleEndian, &f)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}
	fmt.Println("Response frame:")
	fmt.Printf("Port: %d\n", f.Port)
	fmt.Printf("DataKind: %c\n", f.DataKind)
	fmt.Printf("PID: 0x%02X\n", f.PID)
	fmt.Printf("CallFrom: %s\n", f.CallFrom)
	fmt.Printf("CallTo: %s\n", f.CallTo)
	fmt.Printf("DataLen: %d\n", f.DataLen)
	major_verion := (uint16(recv_buf[37]) << 8) | uint16(recv_buf[36])
	minor_verion := (uint16(recv_buf[41]) << 8) | uint16(recv_buf[40])
	fmt.Printf("AGWPE version: %d.%d\n", major_verion, minor_verion)
}
