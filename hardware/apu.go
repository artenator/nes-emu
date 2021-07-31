package hardware

import (
	"github.com/hajimehoshi/oto"
	"log"
)

type Apu struct {
	nes *NES

	// audio device and context
	audioContext *oto.Context
	audioDevice *oto.Player

	// cycle counter
	cyclesPast uint8
	Cyclelimit uint8

	// audioSamples
	audioSamples []float64

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
	cyclesPerSequence uint32

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

	// triangle
	triangle Triangle

	soundOut float64

	// lookup tables
	pulseTable [31]float64
	tndTable [203]float64
}

type Pulse struct {
	apu *Apu
	sweep *Sweep
	baseAddr uint16
	targetTimer uint16
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

type Triangle struct {
	apu *Apu
	curTimer uint16
	countersAddr uint16
	baseAddr uint16
	curTriangleIdx uint8
	linearReload bool
	linearControl bool
	linearCounter uint8
	lengthEnabled bool
	lengthCounter uint8
}

var SampleRate = 48000

var dutyCycles = [4][8]uint8{
	{0, 1, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 0, 0, 0},
	{1, 0, 0, 1, 1, 1, 1, 1},
}

var triangleSequence = [32]uint8{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
}

var lengthTable = [32]uint8{
	10, 254, 20,  2, 40,  4, 80,  6, 160,  8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

func (apu *Apu) InitAPU(initContext bool) {
	// init audio player
	apu.Cyclelimit = 40

	if initContext {
		context, err := oto.NewContext(44100, 1, 1, 4096);
		if err != nil {
			log.Fatal("Audio could not be initialized")
		}

		apu.audioDevice = context.NewPlayer()
	}
	apu.enableDMC = false
	apu.enableNoise = false
	apu.enableTriangle = false
	apu.enablePulseChannel1 = false
	apu.enablePulseChannel2 = false

	apu.cyclesPerSequence = 7457

	apu.audioSamples = make([]float64, apu.Cyclelimit)

	apu.pulse1Addr = 0x4000
	apu.pulse2Addr = 0x4004

	apu.pulse1 = Pulse{apu, &apu.sweep1,apu.pulse1Addr, 0, 0, 0}
	apu.pulse2 = Pulse{apu,&apu.sweep2,apu.pulse2Addr, 0, 0, 0}

	apu.sweep1Addr = 0x4001
	apu.sweep2Addr = 0x4005

	apu.sweep1 = Sweep{&apu.pulse1, apu.sweep1Addr,  false, 0, false, 0, false, 0}
	apu.sweep2 = Sweep{&apu.pulse2, apu.sweep2Addr,  false, 0, false, 0, false, 0}

	apu.triangle = Triangle{apu, 0, 0x4008, 0x400A,0, false, false, 0, false, 0}

	apu.populatePulseTable()
	apu.populateTNDTable()
}

func (apu *Apu) setFrameCounterValues(frameCounterValue uint8) {
	apu.sequencerMode = (frameCounterValue >> 7) & 0x01
	apu.sequenceInterrupt = ((frameCounterValue >> 6) & 0x01) != 0
	apu.sequenceCounter = apu.cyclesPerSequence

	if apu.sequencerMode == 1 {
		apu.quarterFrame()
		apu.halfFrame()
	}
}

func (triangle *Triangle) setLinearCounterValues(linearCounterValue uint8) {
	triangle.linearControl = (linearCounterValue >> 7) & 0x01 == 1
	triangle.lengthEnabled = !triangle.linearControl
}

func (triangle *Triangle) setLengthCounter(value uint8) {
	triangle.lengthCounter = lengthTable[(value >> 3) & 0x1F]
}

func (triangle *Triangle) reloadLinearCounter() {
	triangle.linearCounter = triangle.apu.nes.CPU.Memory[triangle.countersAddr] & 0x7F
}

func (triangle *Triangle) getTriangleTimer() uint16 {
	baseAddr := triangle.baseAddr
	low := uint16(triangle.apu.nes.CPU.Memory[baseAddr])
	high := uint16(triangle.apu.nes.CPU.Memory[baseAddr + 1] & 0x07)

	return (high << 8) | low
}

func (triangle *Triangle) out() uint8 {
	return triangleSequence[triangle.curTriangleIdx]
}

func (triangle *Triangle) triangleRun()  uint8 {
	if triangle.curTimer <= 0 {
		triangle.curTimer = triangle.getTriangleTimer()
		triangle.curTriangleIdx = (triangle.curTriangleIdx + 1) % 32
	} else {
		triangle.curTimer--
	}

	return triangle.out()
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

func (sweep *Sweep) silence() bool {
	targetPeriod := sweep.pulse.curTimer + (sweep.pulse.curTimer >> sweep.shiftCounter)
	if sweep.pulse.curTimer < 8 || (!sweep.negate && targetPeriod > 0x7FF) {
		return true
	} else {
		return false
	}
}

func (sweep *Sweep) sweepRun() {
	if sweep.reload {
		sweep.counter = sweep.period
		sweep.reload = false
	} else if sweep.counter > 0 {
		sweep.counter--
	} else {
		sweep.counter = sweep.period
		if sweep.enabled && !sweep.silence() {
			if sweep.negate {
				sweep.pulse.targetTimer -= (sweep.pulse.targetTimer >> sweep.shiftCounter)
			} else {
				sweep.pulse.targetTimer += (sweep.pulse.targetTimer >> sweep.shiftCounter)
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

func (pulse *Pulse) setTargetTimer() {
	pulse.targetTimer = pulse.getPulTimer()
	pulse.curTimer = pulse.targetTimer
}

func (pulse *Pulse) pulseRun() uint8 {
	if pulse.curTimer <= 0 {
		pulse.curTimer = pulse.targetTimer
		pulse.curDutyIdx = (pulse.curDutyIdx + 1) % 8
	} else {
		pulse.curTimer--
	}

	return pulse.out()
}

func (apu *Apu) APURun() float64 {
	var p1out, p2out, triout uint8

	apu.pulse1.pulseRun()
	apu.pulse2.pulseRun()

	if apu.enablePulseChannel1 && !apu.sweep1.silence() {
		p1out = apu.pulse1.out()
	} else {
		p1out = 0
	}

	if apu.enablePulseChannel2 && !apu.sweep2.silence() {
		p2out = apu.pulse2.out()
	} else {
		p2out = 0
	}

	if apu.enableTriangle {
		triout = apu.triangle.out()
	} else {
		triout = 0
	}

	soundOut := apu.out(p1out, p2out, triout, 0, 0)

	return soundOut
}

func (triangle *Triangle) linearCounterRun() {
	if triangle.linearReload {
		triangle.reloadLinearCounter()
	} else if triangle.linearCounter > 0 {
		triangle.linearCounter--
	}

	if !triangle.linearControl {
		triangle.linearReload = false
	}
}

func (triangle *Triangle) lengthCounterRun() {
	if triangle.lengthEnabled && triangle.lengthCounter > 0 {
		triangle.lengthCounter--
	}
}

func (apu *Apu) halfFrame() {
	apu.sweep1.sweepRun()
	apu.sweep2.sweepRun()
	apu.triangle.lengthCounterRun()
}

func (apu *Apu) quarterFrame() {
	apu.triangle.linearCounterRun()
}

func (apu *Apu) sequenceClockCounterRun() {
	if apu.sequenceCounter > 0 {
		apu.sequenceCounter--
	} else {
		apu.sequenceCounter = apu.cyclesPerSequence
		switch apu.sequencerMode {
		case 1:
			switch apu.sequenceClockCounter {
			case 0:
				apu.quarterFrame()
			case 1:
				apu.quarterFrame()
				apu.halfFrame()
			case 2:
				apu.quarterFrame()
			case 3:
			case 4:
				apu.quarterFrame()
				apu.halfFrame()
			default:
			}
			apu.sequenceClockCounter = (apu.sequenceClockCounter + 1) % 5
		case 0:
			switch apu.sequenceClockCounter {
			case 0:
				apu.quarterFrame()
			case 1:
				apu.quarterFrame()
				apu.halfFrame()
			case 2:
				apu.quarterFrame()
			case 3:
				apu.quarterFrame()
				apu.halfFrame()
			default:
			}
			apu.sequenceClockCounter = (apu.sequenceClockCounter + 1) % 4
		}


	}
}

func (apu *Apu) averageSoundSamples() float64 {
	var sum float64

	for _, sample := range apu.audioSamples{
		sum += sample
	}

	return sum / float64(apu.Cyclelimit)
}

func (apu *Apu) RunAPUCycles(numOfCycles uint16, lastFPS int) {
	for i := uint16(0); i < numOfCycles; i++ {
		if apu.triangle.linearCounter > 0 && apu.triangle.lengthCounter > 0 {
			apu.triangle.triangleRun()
		}

		apu.sequenceClockCounterRun()

		if apu.cyclesPast % 2 == 0 {
			apu.soundOut = apu.APURun()
		}

		if apu.cyclesPast >= apu.Cyclelimit {
			apu.cyclesPast = 0
			apu.audioDevice.Write([]byte{byte(apu.averageSoundSamples() * 0xFF)})
		} else {
			apu.audioSamples[apu.cyclesPast] = apu.soundOut
			apu.cyclesPast++
		}
	}
}