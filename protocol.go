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

// ref: chpt6 / pgs 11-47 https://modbus.org/docs/Modbus_Application_Protocol_V1_1b3.pdf
const (
	ReadCoils                      FunctionCode = 0x01
	ReadDiscreteInputs             FunctionCode = 0x02
	ReadHoldingRegisters           FunctionCode = 0x03
	ReadInputRegisters             FunctionCode = 0x04
	WriteSingleCoil                FunctionCode = 0x05
	WriteSingleRegister            FunctionCode = 0x06
	ReadExceptionStatus            FunctionCode = 0x07 // serial only, optional
	Diagnostics                    FunctionCode = 0x08 // serial only
	GetCommEventCounter            FunctionCode = 0x0B // serial only
	GetCommEventLog                FunctionCode = 0x0C // serial only
	ReadDeviceIdentification       FunctionCode = 0x0E
	WriteMultipleCoils             FunctionCode = 0x0F
	WriteMultipleRegisters         FunctionCode = 0x10
	ReportServerID                 FunctionCode = 0x11
	ReadFileRecord                 FunctionCode = 0x14
	WriteFileRecord                FunctionCode = 0x15
	MaskWriteRegister              FunctionCode = 0x16
	ReadWriteMultipleRegisters     FunctionCode = 0x17
	ReadFIFOQueue                  FunctionCode = 0x18
	EncapsulatedInterfaceTransport FunctionCode = 0x2B
)

func (fc FunctionCode) IsValid() bool {
	switch fc {
	case
		ReadCoils,
		ReadDiscreteInputs,
		ReadHoldingRegisters,
		ReadInputRegisters,
		WriteSingleCoil,
		WriteSingleRegister,
		ReadExceptionStatus,
		Diagnostics,
		GetCommEventCounter,
		GetCommEventLog,
		ReadDeviceIdentification,
		WriteMultipleCoils,
		WriteMultipleRegisters,
		ReportServerID,
		ReadFileRecord,
		WriteFileRecord,
		MaskWriteRegister,
		ReadWriteMultipleRegisters,
		EncapsulatedInterfaceTransport:
		return true
	default:
		return false
	}
}
