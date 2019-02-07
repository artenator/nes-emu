package hardware

import (
	"github.com/hajimehoshi/oto"
)

type Apu struct {
	nes *NES

	// audio device
	audioDevice *oto.Player

	// cycle counter
	cyclesPast uint8
	Cyclelimit uint8

	// Status 0x4015
	enableDMC bool
	enableNoise bool
	enableTriangle bool
	enablePulseChannel1 bool
	enablePulseChannel2 bool

	// Frame Counter 0x4017
	sequenceClockCounter uint8
	sequencerMode uint8
	sequenceInterrupt bool
	sequenceCounter uint32

	// pulse base addrs
	pulse1Addr uint16
	pulse2Addr uint16

	// pulse structs
	pulse1 Pulse
	pulse2 Pulse

	// sweep base addrs
	sweep1Addr uint16
	sweep2Addr uint16

	// sweep structs
	sweep1 Sweep
	sweep2 Sweep

	soundOut float64

	// lookup tables
	pulseTable [31]float64
	tndTable [203]float64
}

type Pulse struct {
	apu *Apu
	baseAddr uint16
	curTimer uint16
	curDutyIdx uint8
}

type Sweep struct {
	pulse *Pulse
	baseAddr uint16

	// bit 7
	enabled bool

	// bits 6 5 4
	period uint8

	// bit 3
	negate bool

	// bit 2 1 0
	shiftCounter uint8

	reload bool

	counter uint8
}

var SampleRate = 48000

var dutyCycles = [4][8]uint8{
	{0, 1, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 0, 0, 0},
	{1, 0, 0, 1, 1, 1, 1, 1},
}

func (apu *Apu) InitAPU() {
	apu.enableDMC = false
	apu.enableNoise = false
	apu.enableTriangle = false
	apu.enablePulseChannel1 = false
	apu.enablePulseChannel2 = false

	apu.pulse1Addr = 0x4000
	apu.pulse2Addr = 0x4004

	apu.pulse1 = Pulse{apu, apu.pulse1Addr, 0, 0,}
	apu.pulse2 = Pulse{apu,apu.pulse2Addr, 0, 0}

	apu.sweep1Addr = 0x4001
	apu.sweep2Addr = 0x4005

	apu.sweep1 = Sweep{&apu.pulse1, apu.sweep1Addr,  false, 0, false, 0, false, 0}
	apu.sweep2 = Sweep{&apu.pulse2, apu.sweep2Addr,  false, 0, false, 0, false, 0}

	apu.populatePulseTable()
	apu.populateTNDTable()
}

func (apu *Apu) setFrameCounterValues(frameCounterValue uint8) {
	apu.sequencerMode = (frameCounterValue >> 7) & 0x01
	apu.sequenceInterrupt = ((frameCounterValue >> 6) & 0x01) != 0
	apu.sequenceCounter = 7475

	apu.sweep1.sweepRun()
	apu.sweep2.sweepRun()
}

func (sweep *Sweep) setSweepValues(sweepValue uint8) {
	sweepEnable := (sweepValue >> 7) & 0x01
	sweepPeriod := (sweepValue >> 4) & 0x07
	sweepNegate := (sweepValue >> 3) & 0x01
	sweepShift := sweepValue & 0x07
	sweep.enabled = sweepEnable != 0
	sweep.period = sweepPeriod
	sweep.negate = sweepNegate != 0
	sweep.shiftCounter = sweepShift
	sweep.reload = true
}

func (sweep *Sweep) sweepRun() {
	if sweep.reload {
		sweep.counter = sweep.period
		sweep.reload = false
	} else if sweep.counter > 0 {
		sweep.counter--
	} else {
		sweep.counter = sweep.period
		if (sweep.enabled) {
			if (sweep.negate) {
				sweep.pulse.curTimer -= (sweep.pulse.curTimer >> sweep.shiftCounter)
			} else {
				sweep.pulse.curTimer += (sweep.pulse.curTimer >> sweep.shiftCounter)
			}
		}
	}
}

func (apu *Apu) getPulseTimer(baseAddr uint16) uint16{
	low := uint16(apu.nes.CPU.Memory[baseAddr + 2])
	high := uint16(apu.nes.CPU.Memory[baseAddr + 3] & 0x07)

	return (high << 8) | low
}

func (apu *Apu) GetPulseFrequency(baseAddr uint16) int {
	return cpuSpeed / int(16 * (apu.getPulseTimer(baseAddr) + 1))
}

//square_table [n] = 95.52 / (8128.0 / n + 100)
func (apu *Apu) populatePulseTable() {
	for i := 0; i < 31; i++ {
		apu.pulseTable[i] = 95.52 / (8128.0 / float64(i) + 100)
	}
}

//square_out = square_table [square1 + square2]
func (apu *Apu) pulseOut(s1, s2 uint8) float64 {
	//apu.pulseTable[]
	return apu.pulseTable[s1 + s2]
}

//tnd_table [n] = 163.67 / (24329.0 / n + 100)
func (apu *Apu) populateTNDTable() {
	for i := 0; i < 203; i++ {
		apu.tndTable[i] = 163.67 / (24329.0 / float64(i) + 100)
	}
}

//tnd_out = tnd_table [3 * triangle + 2 * noise + dmc]
func (apu *Apu) tndOut(t, n, d uint8) float64 {
	return apu.tndTable[3 * t + 2 * n + d]
}

func (apu *Apu) out(s1, s2, t, n, d uint8) float64 {
	return apu.pulseOut(s1, s2) + apu.tndOut(t, n, d)
}

func (pulse *Pulse) getDuty() uint8 {
	duty := uint8((pulse.apu.nes.CPU.Memory[pulse.baseAddr] >> 6) & 3)
	return duty
}

func (pulse *Pulse) getVolume() uint8 {
	volume := uint8((pulse.apu.nes.CPU.Memory[pulse.baseAddr] >> 0) & 0xF)
	return volume
}

func (pulse *Pulse) out() uint8 {
	curDutyPattern := dutyCycles[pulse.getDuty()]
	curDutyValue := curDutyPattern[pulse.curDutyIdx % 8]

	return curDutyValue * pulse.getVolume()
}

func (pulse *Pulse) getPulTimer() uint16{
	baseAddr := pulse.baseAddr
	low := uint16(pulse.apu.nes.CPU.Memory[baseAddr + 2])
	high := uint16(pulse.apu.nes.CPU.Memory[baseAddr + 3] & 0x07)

	return (high << 8) | low
}

func (pulse *Pulse) pulseRun() uint8 {
	if pulse.curTimer <= 0 {
		pulse.curTimer = pulse.getPulTimer()
		pulse.curDutyIdx = (pulse.curDutyIdx + 1) % 8
	} else {
		pulse.curTimer--
	}

	return pulse.out()
}

func (apu *Apu) APURun() float64 {
	var p1out, p2out uint8

	if apu.enablePulseChannel1 {
		p1out = apu.pulse1.pulseRun()
	} else {
		p1out = 0
	}

	if apu.enablePulseChannel2 {
		p2out = apu.pulse2.pulseRun()
	} else {
		p2out = 0
	}

	soundOut := apu.out(p1out, p2out, 0, 0, 0)

	return soundOut
}

func (apu *Apu) sequenceClockCounterRun() {
	if apu.sequenceCounter > 0 {
		apu.sequenceCounter--
	} else {
		switch apu.sequencerMode {
		case 1:
			switch apu.sequenceClockCounter {
			case 0:
			case 1:
				apu.sweep1.sweepRun()
				apu.sweep2.sweepRun()
			case 2:
			case 3:
			case 4:
				apu.sweep1.sweepRun()
				apu.sweep2.sweepRun()
			default:
			}
			apu.sequenceClockCounter = (apu.sequenceClockCounter + 1) % 5
		case 0:
			switch apu.sequenceClockCounter {
			case 0:
			case 1:
				apu.sweep1.sweepRun()
				apu.sweep2.sweepRun()
			case 2:
			case 3:
				apu.sweep1.sweepRun()
				apu.sweep2.sweepRun()
			default:
			}
			apu.sequenceClockCounter = (apu.sequenceClockCounter + 1) % 4
		}


	}
}

func (apu *Apu) RunAPUCycles(numOfCycles uint16, lastFPS int) {
	for i := uint16(0); i < numOfCycles; i++ {

		apu.cyclesPast++

		apu.sequenceClockCounterRun()

		if apu.cyclesPast % 2 == 0 {

			apu.soundOut = apu.APURun()
			if apu.cyclesPast >= apu.Cyclelimit {
				apu.cyclesPast = 0
				apu.audioDevice.Write([]byte{byte(apu.soundOut * 0xFF)})
			}
		}
	}
}