package hardware

import (
	"github.com/faiface/pixel/pixelgl"
)


func (cpu *Cpu) joy1PressButtonA() {
	cpu.Controller |= 1 << 7
}

func (cpu *Cpu) joy1ReleaseButtonA() {
	cpu.Controller &= ^uint8(1 << 7)
}

func (cpu *Cpu) joy1PressButtonB() {
	cpu.Controller |= 1 << 6
}

func (cpu *Cpu) joy1ReleaseButtonB() {
	cpu.Controller &= ^uint8(1 << 6)
}

func (cpu *Cpu) joy1PressButtonSelect() {
	cpu.Controller |= 1 << 5
}

func (cpu *Cpu) joy1ReleaseButtonSelect() {
	cpu.Controller &= ^uint8(1 << 5)
}

func (cpu *Cpu) joy1PressButtonStart() {
	cpu.Controller |= 1 << 4
}

func (cpu *Cpu) joy1ReleaseButtonStart() {
	cpu.Controller &= ^uint8(1 << 4)
}

func (cpu *Cpu) joy1PressButtonUp() {
	cpu.Controller |= 1 << 3
}

func (cpu *Cpu) joy1ReleaseButtonUp() {
	cpu.Controller &= ^uint8(1 << 3)
}

func (cpu *Cpu) joy1PressButtonDown() {
	cpu.Controller |= 1 << 2
}

func (cpu *Cpu) joy1ReleaseButtonDown() {
	cpu.Controller &= ^uint8(1 << 2)
}

func (cpu *Cpu) joy1PressButtonLeft() {
	cpu.Controller |= 1 << 1
}

func (cpu *Cpu) joy1ReleaseButtonLeft() {
	cpu.Controller &= ^uint8(1 << 1)
}

func (cpu *Cpu) joy1PressButtonRight() {
	cpu.Controller |= 1 << 0
}

func (cpu *Cpu) joy1ReleaseButtonRight() {
	cpu.Controller &= ^uint8(1 << 0)
}

func (cpu *Cpu) CheckControllerPresses(win *pixelgl.Window) {
	if win.Pressed(pixelgl.KeyZ) {
		cpu.joy1PressButtonA()
	} else {
		cpu.joy1ReleaseButtonA()
	}

	if win.Pressed(pixelgl.KeyX) {
		cpu.joy1PressButtonB()
	} else {
		cpu.joy1ReleaseButtonB()
	}

	if win.Pressed(pixelgl.KeyRightShift) {
		cpu.joy1PressButtonSelect()
	} else {
		cpu.joy1ReleaseButtonSelect()
	}

	if win.Pressed(pixelgl.KeyS) {
		cpu.joy1PressButtonStart()
	} else {
		cpu.joy1ReleaseButtonStart()
	}

	if win.Pressed(pixelgl.KeyUp) {
		cpu.joy1PressButtonUp()
	} else {
		cpu.joy1ReleaseButtonUp()
	}

	if win.Pressed(pixelgl.KeyDown) {
		cpu.joy1PressButtonDown()
	} else {
		cpu.joy1ReleaseButtonDown()
	}

	if win.Pressed(pixelgl.KeyLeft) {
		cpu.joy1PressButtonLeft()
	} else {
		cpu.joy1ReleaseButtonLeft()
	}

	if win.Pressed(pixelgl.KeyRight) {
		cpu.joy1PressButtonRight()
	} else {
		cpu.joy1ReleaseButtonRight()
	}
}