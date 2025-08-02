package main

import "sync"

type RegisterType uint8

// ref: page 11 https://www.modbus.org/docs/Modbus_Application_Protocol_V1_1b3.pdf
type Hardware struct {
	mu               sync.RWMutex
	coils            map[uint16]bool
	discreteInputs   map[uint16]bool
	inputRegisters   map[uint16]uint16
	holdingRegisters map[uint16]uint16
}

func (h *Hardware) simulate_data() {
	// some llm generated data

	// Motor control outputs (coils 0-15)
	h.coils[0] = false  // Motor 1 Start
	h.coils[1] = false  // Motor 1 Stop
	h.coils[2] = true   // Motor 2 Start (running)
	h.coils[3] = false  // Motor 2 Stop
	h.coils[10] = false // Alarm Reset
	h.coils[11] = true  // System Enable

	// Digital sensor inputs (discrete inputs 0-31)
	h.discreteInputs[0] = true   // Emergency Stop (pressed)
	h.discreteInputs[1] = false  // Door Open
	h.discreteInputs[2] = true   // Pressure OK
	h.discreteInputs[3] = false  // Level High
	h.discreteInputs[4] = true   // Motor 1 Running Feedback
	h.discreteInputs[5] = false  // Motor 1 Fault
	h.discreteInputs[20] = true  // System Ready
	h.discreteInputs[21] = false // Maintenance Mode

	// Analog sensor readings (input registers 0-99)
	h.inputRegisters[0] = 2347   // Temperature (23.47°C * 100)
	h.inputRegisters[1] = 1523   // Pressure (15.23 PSI * 100)
	h.inputRegisters[2] = 875    // Flow Rate (8.75 GPM * 100)
	h.inputRegisters[3] = 4521   // Voltage (45.21V * 100)
	h.inputRegisters[4] = 1234   // Current (12.34A * 100)
	h.inputRegisters[10] = 65535 // Counter (max value)
	h.inputRegisters[50] = 42    // Random sensor

	// Configuration and setpoints (holding registers 0-199)
	h.holdingRegisters[0] = 2500    // Temperature Setpoint (25.00°C * 100)
	h.holdingRegisters[1] = 1500    // Pressure Setpoint (15.00 PSI * 100)
	h.holdingRegisters[2] = 900     // Flow Setpoint (9.00 GPM * 100)
	h.holdingRegisters[10] = 30     // Alarm Delay (seconds)
	h.holdingRegisters[11] = 5      // Retry Count
	h.holdingRegisters[12] = 9600   // Baud Rate
	h.holdingRegisters[100] = 1     // System Mode (1=Auto, 0=Manual)
	h.holdingRegisters[101] = 0     // Error Code
	h.holdingRegisters[102] = 12345 // Serial Number

	// Sparse addressing example - equipment at different stations
	h.coils[1000] = false         // Station 2 Motor
	h.coils[2000] = true          // Station 3 Valve
	h.inputRegisters[5000] = 1876 // Remote Temperature
	h.holdingRegisters[4000] = 50 // Remote Setpoint
}

func NewHardware() *Hardware {
	hardware := &Hardware{
		coils:            make(map[uint16]bool),
		discreteInputs:   make(map[uint16]bool),
		inputRegisters:   make(map[uint16]uint16),
		holdingRegisters: make(map[uint16]uint16),
	}

	hardware.simulate_data()

	return hardware
}
