package main

type RTU_Frame struct {
	slaveID uint8 // ID of slave device
	PDU
	CRC
}

type CRC uint32 // checksum

// TCP frame for Modbus (ref: page 5 of https://modbus.org/docs/Modbus_Messaging_Implementation_Guide_V1_0b.pdf)
type TCP_Frame struct {
	Header
	PDU
}

type Header struct {
	TransactionID uint16
	ProtocolID    uint16 // should always be 0 as this implies Modbus
	Length        uint16 // number of following bytes (I believe for the PDU + UnitID)
	UnitID        uint8  // ID of slave device; I'm told this is often ignored in practice but I will still use it
}

type PDU struct {
	FunctionCode FunctionCode // known list of codes
	Data         []byte
}

// dont like that this could be any arbitrary uint8 instead of constrained to the list but turns out enums in golang are... silly
type FunctionCode uint8

const (
	CoilRead     FunctionCode = 0x01
	RegisterRead FunctionCode = 0x03
)
