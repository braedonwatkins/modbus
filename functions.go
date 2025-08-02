package main

import "fmt"

const MAX_READ_COILS = 0x7D0

func handleReadCoils(hw *Hardware, startAddress uint16, quantity uint16) ([]byte, byte, error) {
	hw.mu.RLock()
	defer hw.mu.RUnlock()

	if quantity == 0 || quantity > MAX_READ_COILS {
		return nil, 0, fmt.Errorf("Number of coils should be between 0 and %d, got %d", MAX_READ_COILS, quantity)
	}

	endAddress := startAddress + quantity
	bits := hw.coils[startAddress:endAddress]
	res, numBytes := PackBits(bits)

	return res, numBytes, nil
}

func handleReadDiscreteInputs()             {}
func handleReadHoldingRegisters()           {}
func handleReadInputRegisters()             {}
func handleWriteSingleCoil()                {}
func handleWriteSingleRegister()            {}
func handleReadExceptionStatus()            {}
func handleDiagnostics()                    {}
func handleGetCommEventCounter()            {}
func handleGetCommEventLog()                {}
func handleReadDeviceIdentification()       {}
func handleWriteMultipleCoils()             {}
func handleWriteMultipleRegisters()         {}
func handleReportServerID()                 {}
func handleReadFileRecord()                 {}
func handleWriteFileRecord()                {}
func handleMaskWriteRegister()              {}
func handleReadWriteMultipleRegisters()     {}
func handleEncapsulatedInterfaceTransport() {}
