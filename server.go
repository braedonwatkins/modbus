package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func RunServer() error {
	listener, err := net.Listen("tcp", "127.0.0.1:502")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Listening on 127.0.0.1:502")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	// TODO: figure out magic number 256 ( i assume this is in the modbus spec )
	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr())
		} else {
			fmt.Printf("Connection read error: %v\n", err)
		}
		return
	}

	req := buf[:n]
	fmt.Printf("Received: %X\n", req)

	frame, err := reqToFrame(req, n)
	if err != nil {
		fmt.Printf("Error parsing frame %s\n", err)
	}
	res := handleRequest(frame)

	n, err = conn.Write([]byte(res))
	if err != nil || n < len(res) {
		// FIXME: figure out how to more gracefully handle this!
		fmt.Printf("Error writing to client")
	}

}

func reqToFrame(req []byte, n int) (*TCP_Frame, error) {
	// REF: page 5 (https://modbus.org/docs/Modbus_Messaging_Implementation_Guide_V1_0b.pdf)

	// validation before mbap header
	if n < MBAP_LEN {
		return nil, errors.New("Buffer is too short to contain MBAP Header")
	}

	mbap := &MBAP{
		TransactionID: binary.BigEndian.Uint16(req[0:2]),
		ProtocolID:    binary.BigEndian.Uint16(req[2:4]),
		Length:        binary.BigEndian.Uint16(req[4:6]),
		UnitID:        req[6],
	}

	// validation before pdu read
	pduLen := int(mbap.Length) - 1 // subtract byte for UnitID
	if pduLen+MBAP_LEN < n {
		return nil, errors.New("Buffer is too short for complete frame")
	}
	functionCode := FunctionCode(req[7])
	if !functionCode.IsValid() {
		return nil, errors.New("Invalid function code")
	}

	pdu := &PDU{
		FunctionCode: functionCode,
		Data:         req[8 : 8+pduLen-1], // -1 for function code byte
	}

	tcp_frame := &TCP_Frame{
		MBAP: *mbap,
		PDU:  *pdu,
	}
	return tcp_frame, nil
}

// FIXME: implement
func handleRequest(frame *TCP_Frame) []byte {
	fmt.Printf("Parsed frame %+v", frame)
	switch frame.FunctionCode {
	case 0x01:
		return []byte{}
	case 0x03:
		return []byte{}
	default:
		return createErrorResponse(frame)
	}
}

func createErrorResponse(frame *TCP_Frame) []byte {
	return []byte("Error for frame")
}
