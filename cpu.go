package main

import "fmt"

// cpu speed
var cpuSpeed = 1789773

// cpu registers
var registers = make(map[string]interface{})

func oraI(immediate byte) {
	registers["ACC"] = immediate | byte(registers["ACC"])
}

// CPURunInstr Runs cpu instruction
func CPURunInstr(instr [2]byte) {
	opcode := instr[0]

	fmt.Printf("Received instruction %d\n", int(opcode))

	
	// switch for different instructions
	switch opcode {
		// ORA immediate
	case 0x09:
		oraI(instr[1])
	default:
		fmt.Println("Invalid opcode")
	}
}

func main() {
	// register defaults
	var acc int8
	var sp uint16
	var pc uint16
	var iX int8
	var iY int8
	var status uint8

	registers["ACC"] = acc
	registers["SP"] = sp
	registers["PC"] = pc
	registers["X"] = iX
	registers["Y"] = iY
	registers["S"] = status
	
        CPURunInstr([2]byte{0x01, 0x00})
}


