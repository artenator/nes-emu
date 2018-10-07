package hardware

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
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
	//cpu := hardware.Cpu{}
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

	for opcode != 0x00 && numOfInstructions < 1200 {
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

		time.Sleep(500 * time.Nanosecond)

		numOfInstructions++



		if numOfInstructions % 100 == 0 {
			if (nes.CPU.Memory[0x2002] >> 7) & 1 == 0 {
				//nes.PPU.setVBlank()
			} else {
				//nes.PPU.clearVBlank()
			}
		}

		opcode = nes.CPU.Read8(nes.CPU.PC)
	}
}