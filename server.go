package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

// simulated hardware
var hw = NewHardware()

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
		fmt.Errorf("Error parsing frame %s\n", err)
	}
	res, err := handleRequest(frame)
	if err != nil {
		fmt.Errorf("Error handling request %s\n", err)
	}

	n, err = conn.Write([]byte(res))
	if err != nil || n < len(res) {
		// FIXME: figure out how to more gracefully handle this!
		fmt.Errorf("Error writing to client")
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

	pdu := &PDU{
		FunctionCode: FunctionCode(req[7]),
		Data:         req[8 : 8+pduLen-1], // -1 for function code byte
	}

	tcp_frame := &TCP_Frame{
		MBAP: *mbap,
		PDU:  *pdu,
	}
	return tcp_frame, nil
}

func handleRequest(frame *TCP_Frame) ([]byte, error) {
	fmt.Printf("Parsed frame %+v", frame)

	// TODO: consider lifting this and passing as a param
	data := frame.PDU.Data

	switch frame.PDU.FunctionCode {
	// FIXME: implement these!!
	case ReadCoils:
		startAddress := binary.BigEndian.Uint16(data[0:2])
		quantity := binary.BigEndian.Uint16(data[2:4])
		return handleReadCoils(hw, startAddress, quantity)
	// case ReadDiscreteInputs:
	// 	return handleReadDiscreteInputs(frame)
	// case ReadHoldingRegisters:
	// 	return handleReadHoldingRegisters(frame)
	// case ReadInputRegisters:
	// 	return handleReadInputRegisters(frame)
	// case WriteSingleCoil:
	// 	return handleWriteSingleCoil(frame)
	// case WriteSingleRegister:
	// 	return handleWriteSingleRegister(frame)
	// case WriteMultipleCoils:
	// 	return handleWriteMultipleCoils(frame)
	// case WriteMultipleRegisters:
	// 	return handleWriteMultipleRegisters(frame)
	// TODO: implement these later
	// case ReadExceptionStatus:
	// 	return handleReadExceptionStatus(frame)
	// case Diagnostics:
	// 	return handleDiagnostics(frame)
	// case GetCommEventCounter:
	// 	return handleGetCommEventCounter(frame)
	// case GetCommEventLog:
	// 	return handleGetCommEventLog(frame)
	// case ReadDeviceIdentification:
	// 	return handleReadDeviceIdentification(frame)
	// case ReportServerID:
	// 	return handleReportServerID(frame)
	// case ReadFileRecord:
	// 	return handleReadFileRecord(frame)
	// case WriteFileRecord:
	// 	return handleWriteFileRecord(frame)
	// case MaskWriteRegister:
	// 	return handleMaskWriteRegister(frame)
	// case ReadWriteMultipleRegisters:
	// 	return handleReadWriteMultipleRegisters(frame)
	// case EncapsulatedInterfaceTransport:
	// 	return handleEncapsulatedInterfaceTransport(frame)
	default:
		return nil, fmt.Errorf("unsupported function code: 0x%X", frame.FunctionCode)
	}
}

func createErrorResponse(frame *TCP_Frame) []byte {
	return []byte("Error for frame")
}
