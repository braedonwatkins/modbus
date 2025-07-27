package main

type RTU_Frame struct {
	slaveID uint8 // ID of slave device
	PDU
	CRC
}

type CRC uint32 // checksum

// TCP frame for Modbus (ref: page 5 of https://modbus.org/docs/Modbus_Messaging_Implementation_Guide_V1_0b.pdf)
type TCP_Frame struct {
	MBAP
	PDU
}

const MBAP_LEN = 7

type MBAP struct {
	TransactionID uint16
	ProtocolID    uint16 // should always be 0 as this implies Modbus
	Length        uint16 // number of following bytes (I believe for the PDU + UnitID)
	UnitID        uint8  // ID of slave device; I'm told this is often ignored in practice but I will still use it
}

type PDU struct {
	FunctionCode FunctionCode // known list of codes
	Data         []byte
}

// turns out enums in golang aren't really a thing! this will do for now
type FunctionCode uint8

const (
	CoilRead     FunctionCode = 0x01
	RegisterRead FunctionCode = 0x03
	//.... FIXME: add the rest!
)

// neat use of method receiver pattern with go
func (fc FunctionCode) IsValid() bool {
	switch fc {
	case CoilRead, RegisterRead:
		return true
	default:
		return false
	}
}
