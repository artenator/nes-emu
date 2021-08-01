package hardware

type instruction struct {
	// Assembly Language Form
	assemblyCode string

	// Assembly Code Enum Form
	code uint8

	// Opcode of the instruction
	opcode uint8

	// Size in Bytes of Instruction
	bytes uint8

	// Num of cycles
	Cycles uint8

	// Addressing Mode
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
	indX = iota

	// Indirect, Y
	// operand is zeropage address
	// effective address is word in (LL, LL + 1) inc with carry
	// C.w($00LL) + Y
	// OPC ($LL), Y
	indY = iota

	// Relative
	// Branch target is PC + signed offset BB
	// Branch offsets are a signed 8 bit value (-128, 127) in 2s complement
	// Page transitions may occur and add an extra cycle to execution
	// OPC $BB
	rel = iota

	// Zero Page
	// operand is the zero page address (high byte 0 == $00LL)
	// OPC $LL
	zpg = iota

	// Zero Page, X
	// operand is zero page address
	// effective address is address incremented by X without carry
	// OPC $LL, X
	zpgX = iota

	// Zero Page, Y
	// operand is zero page address
	// effective address is address incremented by Y without carry
	// OPC $LL, Y
	zpgY = iota

)

const (
	ADC = iota
	AND
	ASL
	BCC
	BCS
	BEQ
	BIT
	BMI
	BNE
	BPL
	BRK
	BVC
	BVS
	CLC
	CLD
	CLI
	CLV
	CMP
	CPX
	CPY
	DCP
	DEC
	DEX
	DEY
	EOR
	INC
	INX
	INY
	ISC
	JMP
	JSR
	LAX
	LDA
	LDX
	LDY
	LSR
	NOP
	ORA
	PHA
	PHP
	PLA
	PLP
	RLA
	ROL
	ROR
	RRA
	RTI
	RTS
	SAX
	SBC
	SEC
	SED
	SEI
	SLO
	SRE
	STA
	STX
	STY
	TAX
	TAY
	TSX
	TXA
	TXS
	TYA

	// not used or not implemented?
	KIL
	ANC
	ALR
	ARR
	XAA
	AXA
	TAS
	SHY
	SHX
	LAS
	AXS
)

var Instructions = [256]instruction{
	// BRK - (Implied)
	instruction{
		"BRK",
		BRK,
		0x00,
		1,
		7,
		impl,},
	// ORA - (Indirect, X)
	instruction{
		"ORA",
		ORA,
		0x01,
		2,
		6,
		indX,},
	// KIL - Implied
	instruction{
		"KIL",
		KIL,
		0x02,
		1,
		0,
		impl,},
	// SLO - (indirect, X)
	instruction{
		"SLO",
		SLO,
		0x03,
		2,
		8,
		indX,},
	// NOP - Zero Page
	instruction{
		"NOP",
		NOP,
		0x04,
		2,
		3,
		zpg,},
	// ORA - Zero Page
	instruction{
		"ORA",
		ORA,
		0x05,
		2,
		3,
		zpg,},
	// ASL - Zero Page
	instruction{
		"ASL",
		ASL,
		0x06,
		2,
		5,
		zpg,},
	// SLO - Zero Page
	instruction{
		"SLO",
		SLO,
		0x07,
		2,
		5,
		zpg,},
	// PHP - Zero Page
	instruction{
		"PHP",
		PHP,
		0x08,
		1,
		3,
		impl,},
	// ORA - Immediate
	instruction{
		"ORA",
		ORA,
		0x09,
		2,
		2,
		imm,},
	// ASL - Accumulator
	instruction{
		"ASL",
		ASL,
		0x0A,
		1,
		2,
		A,},
	// ANC - Immediate
	instruction{
		"ANC",
		ANC,
		0x0B,
		2,
		2,
		imm,},
	// NOP - Absolute
	instruction{
		"NOP",
		NOP,
		0x0C,
		3,
		4,
		abs,},
	// ORA - Absolute
	instruction{
		"ORA",
		ORA,
		0x0D,
		3,
		3,
		abs,},
	// ASL - Absolute
	instruction{
		"ASL",
		ASL,
		0x0E,
		3,
		6,
		abs,},
	// SLO - Absolute
	instruction{
		"SLO",
		SLO,
		0x0F,
		3,
		6,
		abs,},
	// BPL - Relative
	instruction{
		"BPL",
		BPL,
		0x10,
		2,
		2,
		rel,},
	// ORA - (Indirect), Y
	instruction{
		"ORA",
		ORA,
		0x11,
		2,
		5,
		indY,},
	// KIL - Implied
	instruction{
		"KIL",
		KIL,
		0x12,
		1,
		0,
		impl,},
	// SLO - (Indirect), Y
	instruction{
		"SLO",
		SLO,
		0x13,
		2,
		8,
		indY,},
	// NOP - Zero Page
	instruction{
		"NOP",
		NOP,
		0x14,
		2,
		4,
		zpg,},
	// ORA - Zero Page X
	instruction{
		"ORA",
		ORA,
		0x15,
		2,
		4,
		zpgX,},
	// ASL - Zero Page X
	instruction{
		"ASL",
		ASL,
		0x16,
		2,
		6,
		zpgX,},
	// SLO - Zero Page X
	instruction{
		"SLO",
		SLO,
		0x17,
		2,
		6,
		zpgX,},
	// CLC - Implied
	instruction{
		"CLC",
		CLC,
		0x18,
		1,
		2,
		impl,},
	// ORA - Absolute, Y
	instruction{
		"ORA",
		ORA,
		0x19,
		3,
		4,
		absY,},
	// NOP - Implied
	instruction{
		"NOP",
		NOP,
		0x1A,
		1,
		2,
		impl,},
	// SLO - Absolute, Y
	instruction{
		"SLO",
		SLO,
		0x1B,
		3,
		7,
		absY,},
	// NOP - Absolute, X
	instruction{
		"NOP",
		NOP,
		0x1C,
		3,
		4,
		absX,},
	// ORA - Absolute, X
	instruction{
		"ORA",
		ORA,
		0x1D,
		3,
		4,
		absX,},
	// ASL - Absolute, X
	instruction{
		"ASL",
		ASL,
		0x1E,
		3,
		7,
		absX,},
	// SLO - Absolute, X
	instruction{
		"SLO",
		SLO,
		0x1F,
		3,
		7,
		absX,},
	// JSR - Absolute
	instruction{
		"JSR",
		JSR,
		0x20,
		3,
		6,
		abs,},
	// AND - (Indirect, X)
	instruction{
		"AND",
		AND,
		0x21,
		2,
		6,
		indX,},
	// KIL - Implied
	instruction{
		"KIL",
		KIL,
		0x22,
		1,
		0,
		impl,},
	// RLA - (Indirect, X)
	instruction{
		"RLA",
		RLA,
		0x23,
		2,
		8,
		indX,},
	// BIT - Zero Page
	instruction{
		"BIT",
		BIT,
		0x24,
		2,
		3,
		zpg,},
	// AND - Zero Page
	instruction{
		"AND",
		AND,
		0x25,
		2,
		3,
		zpg,},
	// ROL - Zero Page
	instruction{
		"ROL",
		ROL,
		0x26,
		2,
		5,
		zpg,},
	// RLA - Zero Page
	instruction{
		"RLA",
		RLA,
		0x27,
		2,
		5,
		zpg,},
	// PLP - Implied
	instruction{
		"PLP",
		PLP,
		0x28,
		1,
		4,
		impl,},
	// AND - Immediate
	instruction{
		"AND",
		AND,
		0x29,
		2,
		2,
		imm,},
	// ROL - Accumulator
	instruction{
		"ROL",
		ROL,
		0x2A,
		1,
		2,
		A,},
	// ANC - Immediate
	instruction{
		"ANC",
		ANC,
		0x2B,
		2,
		2,
		imm,},
	// BIT - Absolute
	instruction{
		"BIT",
		BIT,
		0x2C,
		3,
		4,
		abs,},
	// AND - Absolute
	instruction{
		"AND",
		AND,
		0x2D,
		3,
		6,
		abs,},
	// ROL - Absolute
	instruction{
		"ROL",
		ROL,
		0x2E,
		3,
		6,
		abs,},
	// RLA - Absolute
	instruction{
		"RLA",
		RLA,
		0x2F,
		3,
		6,
		abs,},
	// BMI - Relative
	instruction{
		"BMI",
		BMI,
		0x30,
		2,
		2,
		rel,},
	// AND - (Indirect), Y
	instruction{
		"AND",
		AND,
		0x31,
		2,
		5,
		indY,},
	// KIL - Implied
	instruction{
		"KIL",
		KIL,
		0x32,
		1,
		0,
		impl,},
	// RLA - (Indirect), Y
	instruction{
		"RLA",
		RLA,
		0x33,
		2,
		8,
		indY,},
	// NOP - Zero Page X
	instruction{
		"NOP",
		NOP,
		0x34,
		2,
		4,
		zpgX,},
	// AND - Zero Page X
	instruction{
		"AND",
		AND,
		0x35,
		2,
		4,
		zpgX,},
	// ROL - Zero Page X
	instruction{
		"ROL",
		ROL,
		0x36,
		2,
		6,
		zpgX,},
	// RLA - (Indrect), Y
	instruction{
		"RLA",
		RLA,
		0x37,
		2,
		6,
		indY,},
	// SEC - Implied
	instruction{
		"SEC",
		SEC,
		0x38,
		1,
		2,
		impl,},
	// AND - Absolute Y
	instruction{
		"AND",
		AND,
		0x39,
		3,
		4,
		absY,},
	// NOP - Implied
	instruction{
		"NOP",
		NOP,
		0x3A,
		1,
		2,
		impl,},
	// RLA - Absolute Y
	instruction{
		"RLA",
		RLA,
		0x3B,
		3,
		7,
		absY,},
	// NOP - Implied
	instruction{
		"NOP",
		NOP,
		0x3C,
		3,
		2,
		impl,},
	// AND - Absolute X
	instruction{
		"AND",
		AND,
		0x3D,
		3,
		4,
		absX,},
	// ROL - Absolute X
	instruction{
		"ROL",
		ROL,
		0x3E,
		3,
		7,
		absX,},
	// RLA - Absolute X
	instruction{
		"RLA",
		RLA,
		0x3F,
		3,
		7,
		absX,},
	// RTI - Implied
	instruction{
		"RTI",
		RTI,
		0x40,
		1,
		6,
		impl,},
	// EOR - (Indirect, X)
	instruction{
		"EOR",
		EOR,
		0x41,
		2,
		5,
		indX,},
	// EOR - Implied
	instruction{
		"EOR",
		EOR,
		0x42,
		1,
		0,
		impl,},
	// SRE - (Indirect, X)
	instruction{
		"SRE",
		SRE,
		0x43,
		2,
		8,
		indX,},
	// NOP - Zero Page
	instruction{
		"NOP",
		NOP,
		0x44,
		2,
		3,
		zpg,},
	// EOR - Zero Page
	instruction{
		"EOR",
		EOR,
		0x45,
		2,
		3,
		zpg,},
	// LSR - Zero Page
	instruction{
		"LSR",
		LSR,
		0x46,
		2,
		5,
		zpg,},
	// SRE - Zero Page
	instruction{
		"SRE",
		SRE,
		0x47,
		2,
		5,
		zpg,},
	// PHA - Implied
	instruction{
		"PHA",
		PHA,
		0x48,
		1,
		3,
		impl,},
	// EOR - Immediate
	instruction{
		"EOR",
		EOR,
		0x49,
		2,
		2,
		imm,},
	// LSR - Accumulator
	instruction{
		"LSR",
		LSR,
		0x4A,
		1,
		2,
		A,},
	// ALR - Immediate
	instruction{
		"ALR",
		ALR,
		0x4B,
		2,
		2,
		imm,},
	// JMP - Absolute
	instruction{
		"JMP",
		JMP,
		0x4C,
		3,
		3,
		abs,},
	// EOR - Absolute 
	instruction{
		"EOR",
		EOR,
		0x4D,
		3,
		4,
		abs,},
	// LSR - Absolute 
	instruction{
		"LSR",
		LSR,
		0x4E,
		3,
		6,
		abs,},
	// SRE - Absolute 
	instruction{
		"SRE",
		SRE,
		0x4F,
		3,
		6,
		abs,},
	// BVC - Relative
	instruction{
		"BVC",
		BVC,
		0x50,
		2,
		2,
		rel,},
	// EOR - (Indirect), Y
	instruction{
		"EOR",
		EOR,
		0x51,
		2,
		5,
		indY,},
	// KIL - Implied 
	instruction{
		"KIL",
		KIL,
		0x52,
		1,
		0,
		impl,},
	// SRE - (Indirect), Y
	instruction{
		"SRE",
		SRE,
		0x53,
		2,
		8,
		indY,},
	// NOP - Zero Page X 
	instruction{
		"NOP",
		NOP,
		0x54,
		2,
		4,
		zpgX,},
	// EOR - Zero Page X 
	instruction{
		"EOR",
		EOR,
		0x55,
		2,
		4,
		zpgX,},
	// LSR - Zero Page X 
	instruction{
		"LSR",
		LSR,
		0x56,
		2,
		6,
		zpgX,},
	// SRE - Zero Page X 
	instruction{
		"SRE",
		SRE,
		0x57,
		2,
		6,
		zpgX,},
	// CLI - Implied
	instruction{
		"CLI",
		CLI,
		0x58,
		1,
		2,
		impl,},
	// EOR - Absolute Y
	instruction{
		"EOR",
		EOR,
		0x59,
		3,
		4,
		absY,},
	// NOP - Implied 
	instruction{
		"NOP",
		NOP,
		0x5A,
		1,
		2,
		impl,},
	// SRE - Absolute Y
	instruction{
		"SRE",
		SRE,
		0x5B,
		3,
		7,
		absY,},
	// NOP - Absolute X
	instruction{
		"NOP",
		NOP,
		0x5C,
		3,
		4,
		absX,},
	// EOR - Absolute X
	instruction{
		"EOR",
		EOR,
		0x5D,
		3,
		4,
		absX,},
	// LSR - Absolute X
	instruction{
		"LSR",
		LSR,
		0x5E,
		3,
		7,
		absX,},
	// SRE - Absolute X
	instruction{
		"SRE",
		SRE,
		0x5F,
		3,
		7,
		absX,},
	// RTS - Implied 
	instruction{
		"RTS",
		RTS,
		0x60,
		1,
		6,
		impl,},
	// ADC - (Indirect, X) 
	instruction{
		"ADC",
		ADC,
		0x61,
		2,
		6,
		indX,},
	// KIL - Implied 
	instruction{
		"KIL",
		KIL,
		0x62,
		1,
		0,
		impl,},
	// RRA - (Indirect, X) 
	instruction{
		"RRA",
		RRA,
		0x63,
		2,
		8,
		indX,},
	// NOP - Zero Page 
	instruction{
		"NOP",
		NOP,
		0x64,
		2,
		3,
		zpg,},
	// ADC - Zero Page
	instruction{
		"ADC",
		ADC,
		0x65,
		2,
		3,
		zpg,},
	// ROR - Zero Page 
	instruction{
		"ROR",
		ROR,
		0x66,
		2,
		5,
		zpg,},
	// RRA - Zero Page 
	instruction{
		"RRA",
		RRA,
		0x67,
		2,
		5,
		zpg,},
	// PLA - Implied 
	instruction{
		"PLA",
		PLA,
		0x68,
		1,
		4,
		impl,},
	// ADC - Immediate 
	instruction{
		"ADC",
		ADC,
		0x69,
		2,
		2,
		imm,},
	// ROR - Accumulator 
	instruction{
		"ROR",
		ROR,
		0x6A,
		1,
		2,
		A,},
	// ARR - Immediate 
	instruction{
		"ARR",
		ARR,
		0x6B,
		2,
		2,
		imm,},
	// JMP - Indirect 
	instruction{
		"JMP",
		JMP,
		0x6C,
		3,
		5,
		ind,},
	// ADC - Absolute 
	instruction{
		"ADC",
		ADC,
		0x6D,
		3,
		4,
		abs,},
	// ROR - Absolute 
	instruction{
		"ROR",
		ROR,
		0x6E,
		3,
		6,
		abs,},
	// RRA - Absolute 
	instruction{
		"RRA",
		RRA,
		0x6F,
		3,
		6,
		abs,},
	// BVS - Relative 
	instruction{
		"BVS",
		BVS,
		0x70,
		2,
		2,
		rel,},
	// ADC - (Indirect), Y
	instruction{
		"ADC",
		ADC,
		0x71,
		2,
		5,
		indY,},
	// KIL - Implied 
	instruction{
		"KIL",
		KIL,
		0x72,
		1,
		9,
		impl,},
	// RRA - (Indirect), Y
	instruction{
		"RRA",
		RRA,
		0x73,
		2,
		8,
		indY,},
	// NOP - Zero Page X
	instruction{
		"NOP",
		NOP,
		0x74,
		2,
		4,
		zpgX,},
	// ADC - Zero Page 
	instruction{
		"ADC",
		ADC,
		0x75,
		2,
		3,
		zpgX,},
	// ROR - Zero Page X
	instruction{
		"ROR",
		ROR,
		0x76,
		2,
		6,
		zpgX,},
	// RRA - Zero Page X 
	instruction{
		"RRA",
		RRA,
		0x77,
		2,
		6,
		zpgX,},
	// SEI - Implied 
	instruction{
		"SEI",
		SEI,
		0x78,
		1,
		2,
		impl,},
	// ADC - Absolute Y
	instruction{
		"ADC",
		ADC,
		0x79,
		3,
		4,
		absY,},
	// NOP - Implied 
	instruction{
		"NOP",
		NOP,
		0x7A,
		1,
		2,
		impl,},
	// RRA - Absolute Y
	instruction{
		"RRA",
		RRA,
		0x7B,
		3,
		7,
		absY,},
	// NOP - Absolute X
	instruction{
		"NOP",
		NOP,
		0x7C,
		3,
		4,
		absX,},
	// ADC - Absolute X
	instruction{
		"ADC",
		ADC,
		0x7D,
		3,
		4,
		absX,},
	// ROR - Absolute X
	instruction{
		"ROR",
		ROR,
		0x7E,
		3,
		7,
		absX,},
	// RRA - Absolute X
	instruction{
		"RRA",
		RRA,
		0x7F,
		3,
		7,
		absX,},
	// NOP - Immediate 
	instruction{
		"NOP",
		NOP,
		0x80,
		2,
		2,
		imm,},
	// STA - (Indirect, X) 
	instruction{
		"STA",
		STA,
		0x81,
		2,
		6,
		indX,},
	// NOP - Immediate 
	instruction{
		"NOP",
		NOP,
		0x82,
		2,
		2,
		imm,},
	// SAX - (Indirect, X) 
	instruction{
		"SAX",
		SAX,
		0x83,
		2,
		6,
		indX,},
	// STY - Zero Page 
	instruction{
		"STY",
		STY,
		0x84,
		2,
		3,
		zpg,},
	// STA - Zero Page 
	instruction{
		"STA",
		STA,
		0x85,
		2,
		3,
		zpg,},
	// STX - Zero Page 
	instruction{
		"STX",
		STX,
		0x86,
		2,
		3,
		zpg,},
	// SAX - Zero Page 
	instruction{
		"SAX",
		SAX,
		0x87,
		2,
		3,
		zpg,},
	// DEY - Implied 
	instruction{
		"DEY",
		DEY,
		0x88,
		1,
		2,
		impl,},
	// NOP - Immediate 
	instruction{
		"NOP",
		NOP,
		0x89,
		2,
		2,
		imm,},
	// TXA - Implied 
	instruction{
		"TXA",
		TXA,
		0x8A,
		1,
		2,
		impl,},
	// XAA - Immediate 
	instruction{
		"XAA",
		XAA,
		0x8B,
		2,
		2,
		imm,},
	// STY - Absolute 
	instruction{
		"STY",
		STY,
		0x8C,
		3,
		4,
		abs,},
	// STA - Absolute 
	instruction{
		"STA",
		STA,
		0x8D,
		3,
		4,
		abs,},
	// STX - Absolute 
	instruction{
		"STX",
		STX,
		0x8E,
		3,
		4,
		abs,},
	// SAX - Absolute 
	instruction{
		"SAX",
		SAX,
		0x8F,
		3,
		4,
		abs,},
	// BCC - Relative 
	instruction{
		"BCC",
		BCC,
		0x90,
		2,
		2,
		rel,},
	// STA - (Indirect), Y 
	instruction{
		"STA",
		STA,
		0x91,
		2,
		6,
		indY,},
	// KIL - Implied 
	instruction{
		"KIL",
		KIL,
		0x92,
		1,
		0,
		impl,},
	// AXA - (Indirect), Y 
	instruction{
		"AXA",
		AXA,
		0x93,
		2,
		6,
		abs,},
	// STY - Zero Page X 
	instruction{
		"STY",
		STY,
		0x94,
		2,
		4,
		zpgX,},
	// STA - Zero Page X 
	instruction{
		"STA",
		STA,
		0x95,
		2,
		4,
		zpgX,},
	// STX - Zero Page Y 
	instruction{
		"STX",
		STX,
		0x96,
		2,
		4,
		zpgY,},
	// SAX - Zero Page Y 
	instruction{
		"SAX",
		SAX,
		0x97,
		2,
		4,
		zpgY,},
	// TYA - Implied 
	instruction{
		"TYA",
		TYA,
		0x98,
		1,
		2,
		impl,},
	// STA - Absolute Y 
	instruction{
		"STA",
		STA,
		0x99,
		3,
		5,
		absY,},
	// TXS - Implied 
	instruction{
		"TXS",
		TXS,
		0x9A,
		1,
		2,
		impl,},
	// TAS - Absolute Y
	instruction{
		"TAS",
		TAS,
		0x9B,
		3,
		5,
		absY,},
	// SHY - Absolute X
	instruction{
		"SHY",
		SHY,
		0x9C,
		3,
		5,
		absX,},
	// STA - Absolute X
	instruction{
		"STA",
		STA,
		0x9D,
		3,
		5,
		absX,},
	// SHX - Absolute Y
	instruction{
		"SHX",
		SHX,
		0x9E,
		3,
		5,
		absY,},
	// EOR - Absolute 
	instruction{
		"EOR",
		EOR,
		0x9F,
		3,
		4,
		abs,},
	// LDY - Immediate 
	instruction{
		"LDY",
		LDY,
		0xA0,
		2,
		3,
		imm,},
	// LDA - (Indirect, X) 
	instruction{
		"LDA",
		LDA,
		0xA1,
		2,
		6,
		indX,},
	// LDX - Immediate 
	instruction{
		"LDX",
		LDX,
		0xA2,
		2,
		2,
		imm,},
	// LAX - (Indirect, X) 
	instruction{
		"LAX",
		LAX,
		0xA3,
		2,
		6,
		indX,},
	// LDY - Zero Page 
	instruction{
		"LDY",
		LDY,
		0xA4,
		2,
		3,
		zpg,},
	// LDA - Zero Page 
	instruction{
		"LDA",
		LDA,
		0xA5,
		2,
		3,
		zpg,},
	// LDX - Zero Page 
	instruction{
		"LDX",
		LDX,
		0xA6,
		2,
		3,
		zpg,},
	// LAX - Zero Page 
	instruction{
		"LAX",
		LAX,
		0xA7,
		2,
		3,
		zpg,},
	// TAY - Implied 
	instruction{
		"TAY",
		TAY,
		0xA8,
		1,
		2,
		impl,},
	// LDA - Immediate 
	instruction{
		"LDA",
		LDA,
		0xA9,
		2,
		2,
		imm,},
	// TAX - Implied 
	instruction{
		"TAX",
		TAX,
		0xAA,
		1,
		2,
		impl,},
	// LAX - Immediate 
	instruction{
		"LAX",
		LAX,
		0xAB,
		2,
		2,
		imm,},
	// LDY - Absolute 
	instruction{
		"LDY",
		LDY,
		0xAC,
		3,
		4,
		abs,},
	// LDA - Absolute 
	instruction{
		"LDA",
		LDA,
		0xAD,
		3,
		4,
		abs,},
	// LDX - Absolute 
	instruction{
		"LDX",
		LDX,
		0xAE,
		3,
		4,
		abs,},
	// LAX - Absolute 
	instruction{
		"LAX",
		LAX,
		0xAF,
		3,
		4,
		abs,},
	// BCS - Relative 
	instruction{
		"BCS",
		BCS,
		0xB0,
		2,
		2,
		rel,},
	// LDA - (Indirect), Y 
	instruction{
		"LDA",
		LDA,
		0xB1,
		2,
		5,
		indY,},
	// KIL - Implied 
	instruction{
		"KIL",
		KIL,
		0xB2,
		1,
		0,
		impl,},
	// LAX - (Indirect), Y 
	instruction{
		"LAX",
		LAX,
		0xB3,
		2,
		5,
		indY,},
	// LDY - Zero Page X 
	instruction{
		"LDY",
		LDY,
		0xB4,
		2,
		4,
		zpgX,},
	// LDA - Zero Page X 
	instruction{
		"LDA",
		LDA,
		0xB5,
		2,
		4,
		zpgX,},
	// LDX - Zero Page Y 
	instruction{
		"LDX",
		LDX,
		0xB6,
		2,
		4,
		zpgY,},
	// LAX - Zero Page Y 
	instruction{
		"LAX",
		LAX,
		0xB7,
		2,
		4,
		zpgY,},
	// CLV - Implied 
	instruction{
		"CLV",
		CLV,
		0xB8,
		1,
		2,
		impl,},
	// LDA - Absolute Y
	instruction{
		"LDA",
		LDA,
		0xB9,
		3,
		4,
		absY,},
	// TSX - Implied 
	instruction{
		"TSX",
		TSX,
		0xBA,
		1,
		2,
		impl,},
	// LAS - Absolute Y
	instruction{
		"LAS",
		LAS,
		0xBB,
		3,
		4,
		absY,},
	// LDY - Absolute X
	instruction{
		"LDY",
		LDY,
		0xBC,
		3,
		4,
		absX,},
	// LDA - Absolute X
	instruction{
		"LDA",
		LDA,
		0xBD,
		3,
		4,
		absX,},
	// LDX - Absolute Y
	instruction{
		"LDX",
		LDX,
		0xBE,
		3,
		4,
		absY},
	// LAX - Absolute Y
	instruction{
		"LAX",
		LAX,
		0xBF,
		3,
		4,
		absY,},
	// CPY - Immediate 
	instruction{
		"CPY",
		CPY,
		0xC0,
		2,
		2,
		imm,},
	// CMP - (Indirect, X) 
	instruction{
		"CMP",
		CMP,
		0xC1,
		2,
		6,
		indX,},
	// NOP - Immediate 
	instruction{
		"NOP",
		NOP,
		0xC2,
		2,
		2,
		imm,},
	// DCP - (Indirect, X) 
	instruction{
		"DCP",
		DCP,
		0xC3,
		2,
		8,
		indX,},
	// CPY - Zero Page 
	instruction{
		"CPY",
		CPY,
		0xC4,
		2,
		3,
		zpg,},
	// CMP - Zero Page 
	instruction{
		"CMP",
		CMP,
		0xC5,
		2,
		3,
		zpg,},
	// DEC - Zero Page 
	instruction{
		"DEC",
		DEC,
		0xC6,
		2,
		5,
		zpg,},
	// DCP - Zero Page 
	instruction{
		"DCP",
		DCP,
		0xC7,
		2,
		5,
		zpg,},
	// INY - Implied 
	instruction{
		"INY",
		INY,
		0xC8,
		1,
		2,
		impl,},
	// CPY - Immediate 
	instruction{
		"CMP",
		CMP,
		0xC9,
		2,
		2,
		imm,},
	// DEX - Implied 
	instruction{
		"DEX",
		DEX,
		0xCA,
		1,
		2,
		impl,},
	// AXS - Immediate 
	instruction{
		"AXS",
		AXS,
		0xCB,
		2,
		2,
		imm,},
	// CPY - Absolute 
	instruction{
		"CPY",
		CPY,
		0xCC,
		3,
		4,
		abs,},
	// CMP - Absolute 
	instruction{
		"CMP",
		CMP,
		0xCD,
		3,
		4,
		abs,},
	// DEC - Absolute 
	instruction{
		"DEC",
		DEC,
		0xCE,
		3,
		3,
		abs,},
	// DCP - Absolute 
	instruction{
		"DCP",
		DCP,
		0xCF,
		3,
		6,
		abs,},
	// BNE - Relative 
	instruction{
		"BNE",
		BNE,
		0xD0,
		2,
		2,
		rel,},
	// CMP - (Indirect), Y
	instruction{
		"CMP",
		CMP,
		0xD1,
		2,
		5,
		indY,},
	// KIL - Implied 
	instruction{
		"KIL",
		KIL,
		0xD2,
		1,
		0,
		impl,},
	// DCP - (Indirect), Y 
	instruction{
		"DCP",
		DCP,
		0xD3,
		2,
		8,
		indY,},
	// NOP - Zero Page X 
	instruction{
		"NOP",
		NOP,
		0xD4,
		2,
		4,
		zpgX,},
	// CMP - Zero Page X 
	instruction{
		"CMP",
		CMP,
		0xD5,
		2,
		4,
		zpgX,},
	// DEC - Zero Page X 
	instruction{
		"DEC",
		DEC,
		0xD6,
		2,
		6,
		zpgX,},
	// DCP - Zero Page X 
	instruction{
		"DCP",
		DCP,
		0xD7,
		2,
		6,
		zpgX,},
	// CLD - Implied 
	instruction{
		"CLD",
		CLD,
		0xD8,
		1,
		2,
		impl,},
	// CMP - Absolute Y
	instruction{
		"CMP",
		CMP,
		0xD9,
		3,
		4,
		absY,},
	// NOP - Implied 
	instruction{
		"NOP",
		NOP,
		0xDA,
		1,
		2,
		impl,},
	// DCP - Absolute Y
	instruction{
		"DCP",
		DCP,
		0xDB,
		3,
		7,
		absY,},
	// NOP - Absolute X
	instruction{
		"NOP",
		NOP,
		0xDC,
		3,
		4,
		absX,},
	// CMP - Absolute X
	instruction{
		"CMP",
		CMP,
		0xDD,
		3,
		4,
		absX,},
	// DEC - Absolute X
	instruction{
		"DEC",
		DEC,
		0xDE,
		3,
		7,
		absX,},
	// DCP - Absolute X
	instruction{
		"DCP",
		DCP,
		0xDF,
		3,
		7,
		absX,},
	// CPX - Immediate 
	instruction{
		"CPX",
		CPX,
		0xE0,
		2,
		2,
		imm,},
	// SBC - (Indirect, X) 
	instruction{
		"SBC",
		SBC,
		0xE1,
		2,
		6,
		indX,},
	// NOP - Immediate 
	instruction{
		"NOP",
		NOP,
		0xE2,
		2,
		2,
		imm,},
	// ISC - (Indirect, X) 
	instruction{
		"ISC",
		ISC,
		0xE3,
		2,
		8,
		indX,},
	// CPX - Zero Page 
	instruction{
		"CPX",
		CPX,
		0xE4,
		2,
		3,
		zpg,},
	// SBC - Zero Page 
	instruction{
		"SBC",
		SBC,
		0xE5,
		2,
		3,
		zpg,},
	// INC - Zero Page 
	instruction{
		"INC",
		INC,
		0xE6,
		2,
		5,
		zpg,},
	// ISC - Zero Page 
	instruction{
		"ISC",
		ISC,
		0xE7,
		2,
		5,
		zpg,},
	// INX - Implied 
	instruction{
		"INX",
		INX,
		0xE8,
		1,
		2,
		impl,},
	// SBC - Immediate 
	instruction{
		"SBC",
		SBC,
		0xE9,
		2,
		2,
		imm,},
	// NOP - Implied
	instruction{
		"NOP",
		NOP,
		0xEA,
		1,
		2,
		impl,},
	// SBC - Immediate 
	instruction{
		"SBC",
		SBC,
		0xEB,
		2,
		2,
		imm,},
	// CPX - Absolute 
	instruction{
		"CPX",
		CPX,
		0xEC,
		3,
		4,
		abs,},
	// SBC - Absolute 
	instruction{
		"SBC",
		SBC,
		0xED,
		3,
		4,
		abs,},
	// INC - Absolute 
	instruction{
		"INC",
		INC,
		0xEE,
		3,
		6,
		abs,},
	// ISC - Absolute 
	instruction{
		"ISC",
		ISC,
		0xEF,
		3,
		6,
		abs,},
	// BEQ - Relative 
	instruction{
		"BEQ",
		BEQ,
		0xF0,
		2,
		2,
		rel,},
	// SBC - (Indirect), Y
	instruction{
		"SBC",
		SBC,
		0xF1,
		2,
		5,
		indY,},
	// KIL - Implied 
	instruction{
		"KIL",
		KIL,
		0xF2,
		1,
		0,
		impl,},
	// ISC - (Indirect), Y 
	instruction{
		"ISC",
		ISC,
		0xF3,
		2,
		8,
		indY,},
	// NOP - Zero Page X 
	instruction{
		"NOP",
		NOP,
		0xF4,
		2,
		4,
		zpgX,},
	// SBC - Zero Page X 
	instruction{
		"SBC",
		SBC,
		0xF5,
		2,
		4,
		zpgX,},
	// INC - Zero Page X 
	instruction{
		"INC",
		INC,
		0xF6,
		2,
		6,
		zpgX,},
	// ISC - Zero Page X 
	instruction{
		"ISC",
		ISC,
		0xF7,
		2,
		6,
		zpgX,},
	// SED - Implied 
	instruction{
		"SED",
		SED,
		0xF8,
		1,
		2,
		impl,},
	// SBC - Absolute Y
	instruction{
		"SBC",
		SBC,
		0xF9,
		3,
		4,
		absY,},
	// NOP - Implied 
	instruction{
		"NOP",
		NOP,
		0xFA,
		1,
		2,
		impl,},
	// ISC - Absolute Y
	instruction{
		"ISC",
		ISC,
		0xFB,
		3,
		7,
		absY,},
	// NOP - Absolute X
	instruction{
		"NOP",
		NOP,
		0xFC,
		3,
		4,
		absX,},
	// SBC - Absolute X
	instruction{
		"SBC",
		SBC,
		0xFD,
		3,
		4,
		absX,},
	// INC - Absolute X 
	instruction{
		"INC",
		INC,
		0xFE,
		3,
		7,
		absX,},
	// ISC - Absolute X
	instruction{
		"ISC",
		ISC,
		0xFF,
		3,
		7,
		absX,},
}
