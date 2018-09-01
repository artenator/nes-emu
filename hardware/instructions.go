package hardware

type instruction struct {
	// Assembly Language Form
	assemblyCode string

	// Opcode of the instruction
	opcode uint8

	// Size in Bytes of Instruction
	bytes uint8

	// Num of cycles
	cycles uint8

	//Addressing Mode
	mode uint8
}

const (
	// Accumulator
	// Operates on data in the accumulator.
	// No operands needed
	// OPC A
	A = iota

	// Absolute
	// Takes an absolute address 16 bit
	// Any address in memory $LLHH
	// OPC $LLHH
	abs = iota

	// Absolute, X
	// operand is the address
	// effective address is address incremented by X with carry
	// OPC $LLHH,X
	absX = iota

	// Absolute, Y
	// operand is the address
	// effective address is address incremented by Y with carry
	// OPC $LLHH,Y
	absY = iota

	// Immediate mode
	// operand is given by the instruction
	// one byte immediate $BB
	// OPC #$BB
	imm = iota

	// Implied
	// No operand addresses required
	// implied by the instruction itself
	// OPC
	impl = iota

	// Indirect
	// operand is the address
	// effective address is contents of word at address $LLHH
	// OPC ($LLHH)
	ind = iota

	// X, Indirect
	// operand is zeropage address
	// effective address is word in (LL + X, LL + X + 1), inc without carry
	// C.w($00LL + X)
	// OPC ($LL, X)
	Xind = iota

	// Indirect, Y
	// operand is zeropage address
	// effective address is word in (LL, LL + 1) inc with carry
	// C.w($00LL) + Y
	// OPC ($LL), Y
	indY = iota

	// Relative

)

var Instructions = [256]instruction{
	// BRK
	instruction{"BRK",
				0x00,
				1,
				7,
				impl,},
	// ORA - (Indirect, X)
	instruction{"ORA",
				0x01,
				2,
				6,
				impl,},
	// KIL
	instruction{"KIL",
				0x02,
				1,
				0,},
	// SLO
	instruction{"SLO",
				0x03,
				2,
				8,},
	// NOP
	instruction{"NOP",
				0x04,
				2,
				3,},
	// ORA - Zero Page
	instruction{"ORA",
				0x05,
				2,
				3,},
	// ASL - Zero Page
	instruction{"ASL",
				0x06,
				2,
				5,},
	// SLO - Zero Page
	instruction{"SLO",
				0x07,
				2,
				5,},
	// PHP - Zero Page
	instruction{"PHP",
				0x08,
				1,
				3,},
}