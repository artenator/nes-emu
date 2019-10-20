package hardware

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type expectedState struct {
	PC uint16
	A uint8
	X uint8
	Y uint8
	P uint8
}

func extractHex(s, repattern string) uint64 {
	re := regexp.MustCompile(repattern)
	reMatch := strings.Split(re.FindString(s), ":")
	result, _ := strconv.ParseUint(reMatch[len(reMatch) - 1], 16, 64)
	return result
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]expectedState, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var expStates []expectedState
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		curStr := scanner.Text()
		ix := extractHex(curStr, "X:[0-9A-F]{2}")
		iy := extractHex(curStr, "Y:([0-9A-F]{2})")
		ia := extractHex(curStr, "A:[0-9A-F]{2}")
		ip := extractHex(curStr, "P:[0-9A-F]{2}")
		ipc := extractHex(curStr, "[0-9A-F]{4}")

		curState := expectedState{
			PC: uint16(ipc),
			X: uint8(ix),
			Y: uint8(iy),
			A: uint8(ia),
			P: uint8(ip),
		}
		expStates = append(expStates, curState)
	}


	return expStates, scanner.Err()
}

func TestCpu(t *testing.T) {
	// create new nes
	nes := NewNES()

	cart, err := CreateCartridge("nestest.nes")

	expected, _ := readLines("nestest.log")

	firstInstruction := uint16(0xC000)

	// Set the PC to be at the address
	nes.CPU.PC = firstInstruction

	nes.CPU.setCpuInitialState()

	if err != nil {
		log.Println(err)
	} else {
		nes.LoadCartridge(cart)
	}

	// number of instructions ran
	var numOfInstructions uint = 0

	opcode := nes.CPU.Read8(nes.CPU.PC)

	for opcode != 0x00 {
		if nes.CPU.PC != expected[numOfInstructions].PC {
			t.Errorf("Wrong PC. Expected %02x but got %02x\n %+v\nPC:%02x", expected[numOfInstructions].PC, nes.CPU.PC, Instructions[opcode], nes.CPU.PC)
		}

		if nes.CPU.A != expected[numOfInstructions].A {
			t.Errorf("Wrong Acc. Expected %02x but got %02x\n %+v\nPC:%02x", expected[numOfInstructions].A, nes.CPU.A, Instructions[opcode], nes.CPU.PC)
		}

		if nes.CPU.X != expected[numOfInstructions].X {
			t.Errorf("Wrong X. Expected %02x but got %02x\n %+v\nPC:%02x", expected[numOfInstructions].X, nes.CPU.X, Instructions[opcode], nes.CPU.PC)
		}

		if nes.CPU.Y != expected[numOfInstructions].Y {
			t.Errorf("Wrong Y. Expected %02x but got %02x\n %+v\nPC:%02x", expected[numOfInstructions].Y, nes.CPU.Y, Instructions[opcode], nes.CPU.PC)
		}

		if nes.CPU.P != expected[numOfInstructions].P {
			t.Errorf("Wrong P. Expected %02x but got %02x\n %+v\nPC:%02x", expected[numOfInstructions].P, nes.CPU.P, Instructions[opcode], nes.CPU.PC)
		}

		nes.CPU.RunInstruction(Instructions[opcode], false)

		numOfInstructions++

		opcode = nes.CPU.Read8(nes.CPU.PC)
	}
}

func runBlarggTest(filename string) string {
	nes := NewNES()

	cart, err := CreateCartridge(filename)

	if err != nil {
		log.Println(err)
	} else {
		nes.LoadCartridge(cart)
	}

	nes.CPU.Reset();

	nes.APU.InitAPU(false)

	opcode := nes.CPU.Read8(nes.CPU.PC)

	log.Printf("%+v", nes.CPU.PC)

	testStarted := false

	for !testStarted {
		nes.CPU.RunInstruction(Instructions[opcode], false)
		opcode = nes.CPU.Read8(nes.CPU.PC)

		if nes.CPU.Memory[0x6000] == 0x80 {
			testStarted = true
		}

		if nes.CPU.Memory[0x6000] == 0x81 {
			nes.CPU.Reset()
		}
	}

	for nes.CPU.Memory[0x6000] > 0x7F {
		nes.CPU.RunInstruction(Instructions[opcode], false)
		opcode = nes.CPU.Read8(nes.CPU.PC)

		if nes.CPU.Memory[0x6000] == 0x81 {
			nes.CPU.Reset()
		}
	}

	var testMsgByteArr []byte
	curIdx := 0x6004
	curByte := nes.CPU.Memory[curIdx]

	for curByte != 0 {
		testMsgByteArr = append(testMsgByteArr, curByte)
		curIdx++
		curByte = nes.CPU.Memory[curIdx]
	}

	return string(testMsgByteArr)
}

func TestBlarggCpu01(t *testing.T) {
	// create new nes

	resultMsg := runBlarggTest("./blargg_cpu_singles/01-basics.nes")

	if !strings.Contains(strings.ToUpper(resultMsg), "PASSED") {
		t.Errorf("Blargg Test did not pass\nMESSAGE: %s", resultMsg)
	}
}

func TestBlarggCpu02(t *testing.T) {

	resultMsg := runBlarggTest("./blargg_cpu_singles/02-implied.nes")

	if !strings.Contains(strings.ToUpper(resultMsg), "PASSED") {
		t.Errorf("Blargg Test did not pass\nMESSAGE: %s", resultMsg)
	}
}

func TestBlarggCpu03(t *testing.T) {

	resultMsg := runBlarggTest("./blargg_cpu_singles/03-immediate.nes")

	if !strings.Contains(strings.ToUpper(resultMsg), "PASSED") {
		t.Errorf("Blargg Test did not pass\nMESSAGE: %s", resultMsg)
	}
}

func TestBlarggCpu04(t *testing.T) {

	resultMsg := runBlarggTest("./blargg_cpu_singles/04-zero_page.nes")

	if !strings.Contains(strings.ToUpper(resultMsg), "PASSED") {
		t.Errorf("Blargg Test did not pass\nMESSAGE: %s", resultMsg)
	}
}

func TestBlarggCpu05(t *testing.T) {

	resultMsg := runBlarggTest("./blargg_cpu_singles/05-zp_xy.nes")

	if !strings.Contains(strings.ToUpper(resultMsg), "PASSED") {
		t.Errorf("Blargg Test did not pass\nMESSAGE: %s", resultMsg)
	}
}

func TestBlarggCpu06(t *testing.T) {

	resultMsg := runBlarggTest("./blargg_cpu_singles/06-absolute.nes")

	if !strings.Contains(strings.ToUpper(resultMsg), "PASSED") {
		t.Errorf("Blargg Test did not pass\nMESSAGE: %s", resultMsg)
	}
}

func TestBlarggCpu07(t *testing.T) {

	resultMsg := runBlarggTest("./blargg_cpu_singles/07-abs_xy.nes")

	if !strings.Contains(strings.ToUpper(resultMsg), "PASSED") {
		t.Errorf("Blargg Test did not pass\nMESSAGE: %s", resultMsg)
	}
}

func TestBlarggCpu08(t *testing.T) {

	resultMsg := runBlarggTest("./blargg_cpu_singles/08-ind_x.nes")

	if !strings.Contains(strings.ToUpper(resultMsg), "PASSED") {
		t.Errorf("Blargg Test did not pass\nMESSAGE: %s", resultMsg)
	}
}

func TestBlarggCpu09(t *testing.T) {

	resultMsg := runBlarggTest("./blargg_cpu_singles/09-ind_y.nes")

	if !strings.Contains(strings.ToUpper(resultMsg), "PASSED") {
		t.Errorf("Blargg Test did not pass\nMESSAGE: %s", resultMsg)
	}
}