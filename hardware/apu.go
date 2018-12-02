package hardware

type Apu struct {
	nes *NES

	// Status 0x4015
	enableDMC bool
	enableNoise bool
	enableTriangle bool
	enablePulseChannel1 bool
	enablePulseChannel2 bool

	// pulse base addrs
	pulse1Addr uint16
	pulse2Addr uint16

	// lookup tables
	pulseTable [31]float32
	tndTable [203]float32
}

func (apu *Apu) InitAPU() {
	apu.enableDMC = false
	apu.enableNoise = false
	apu.enableTriangle = false
	apu.enablePulseChannel1 = false
	apu.enablePulseChannel2 = false

	apu.pulse1Addr = 0x4000
	apu.pulse2Addr = 0x4004

	apu.populatePulseTable()
}

func (apu *Apu) getPulseTimer(baseAddr uint16) uint16{
	low := uint16(apu.nes.CPU.Memory[baseAddr + 2])
	high := uint16(apu.nes.CPU.Memory[baseAddr + 3] & 0x07)

	return (high << 11) | low
}

func (apu *Apu) getPulseFrequency(baseAddr uint16) int {
	return cpuSpeed / int(16 * (apu.getPulseTimer(baseAddr) + 1))
}

//square_table [n] = 95.52 / (8128.0 / n + 100)
func (apu *Apu) populatePulseTable() {
	for i := 0; i < 31; i++ {
		apu.pulseTable[i] = 95.52 / (8128.0 / float32(i) + 100)
	}
}

//square_out = square_table [square1 + square2]
func (apu *Apu) pulseOut() float32 {
	//apu.pulseTable[]
	return 1.2213
}