package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/shirou/gopsutil/cpu"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"log"
	"nes-emu/hardware"
	"time"
)

var imd = imdraw.New(nil)
var cpuInfo, _ = cpu.Info()

type ColorAtPixel struct {
	color hardware.Color
	x uint8
	y uint8
}

func drawPixel(c hardware.Color, x, y float64) {
	colorRGBA := pixel.RGB(float64(c.R)/255, float64(c.G)/255, float64(c.B)/255)
	imd.Color = colorRGBA
	imd.Push(pixel.V(x+1, y+1))
	imd.Push(pixel.V(x, y))
	imd.Rectangle(0)
}

func runNESInstruction(nes hardware.NES, numOfInstructions *uint) {
	opcode := nes.CPU.Read8(nes.CPU.PC)
	instr := hardware.Instructions[opcode]
	nes.CPU.RunInstruction(instr, false)
	nes.PPU.RunPPUCycles(uint16(3 * instr.Cycles))

	inVBlank := (nes.CPU.Memory[0x2002]>>7)&1 == 1
	NMIEnabled := (nes.CPU.Memory[0x2000]>>7)&1 == 1

	*numOfInstructions++

	if inVBlank && NMIEnabled {
		nes.PPU.ClearVBlank()
		nes.CPU.HandleNMI()
	}
}

func runNESTurbo(nes hardware.NES, numOfInstructions *uint) {
	for true {
		runNESInstruction(nes, numOfInstructions)
	}
}

func runNEStoFrame(nes hardware.NES, numOfInstructions *uint) {
	for !nes.PPU.FrameReady {
		runNESInstruction(nes, numOfInstructions)
	}

	nes.PPU.FrameReady = false
}

func run() {
	// all of our code will be fired up from here
	cfg := pixelgl.WindowConfig{
		Title:  "Arte's NES Emulator",
		Bounds: pixel.R(0, 0, 256, 240),
		VSync:  false,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	nes := hardware.NewNES()

	cart, err := hardware.CreateCartridge("balloon-fight.nes")

	if err != nil {
		log.Println(err)
	} else {
		nes.LoadCartridge(cart)
		nes.CPU.Reset()
	}

	var (
		numOfInstructions uint = 0
		frames = 0
		us = time.Tick(16666 * time.Microsecond)
		second = time.Tick(time.Second)
	)

	// main drawing loop
	for !win.Closed() {
		var img *image.RGBA = image.NewRGBA(image.Rect(0, 0, 256, 240))

		if win.Pressed(pixelgl.KeyZ) {
			nes.CPU.Joy1PressButtonA()
		} else {
			nes.CPU.Joy1ReleaseButtonA()
		}

		if win.Pressed(pixelgl.KeyX) {
			nes.CPU.Joy1PressButtonB()
		} else {
			nes.CPU.Joy1ReleaseButtonB()
		}

		if win.Pressed(pixelgl.KeyRightShift) {
			nes.CPU.Joy1PressButtonSelect()
		} else {
			nes.CPU.Joy1ReleaseButtonSelect()
		}

		if win.Pressed(pixelgl.KeyS) {
			nes.CPU.Joy1PressButtonStart()
		} else {
			nes.CPU.Joy1ReleaseButtonStart()
		}

		if win.Pressed(pixelgl.KeyUp) {
			nes.CPU.Joy1PressButtonUp()
		} else {
			nes.CPU.Joy1ReleaseButtonUp()
		}

		if win.Pressed(pixelgl.KeyDown) {
			nes.CPU.Joy1PressButtonDown()
		} else {
			nes.CPU.Joy1ReleaseButtonDown()
		}

		if win.Pressed(pixelgl.KeyLeft) {
			nes.CPU.Joy1PressButtonLeft()
		} else {
			nes.CPU.Joy1ReleaseButtonLeft()
		}

		if win.Pressed(pixelgl.KeyRight) {
			nes.CPU.Joy1PressButtonRight()
		} else {
			nes.CPU.Joy1ReleaseButtonRight()
		}

		runNEStoFrame(nes, &numOfInstructions)



		imd.Clear()

		for y := 0; y < 240; y++ {
			for x := 0; x < 256; x++ {
				c := nes.PPU.GetColorAtPixel(uint8(x), uint8(y))
				img.SetRGBA(x, y, color.RGBA{c.R, c.G, c.B, uint8(c.A)})
				//drawPixel(c, float64(x), float64(239 - y))
			}
		}

		pic := pixel.PictureDataFromImage(img)
		sprite := pixel.NewSprite(pic, pic.Bounds())

		win.Clear(colornames.Black)
		//imd.Draw(win)
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()

		<-us
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("FPS: %d %s", frames, cfg.Title))
			frames = 0
		default:
		}
	}
}

func main() {
	log.Println("Arte's NES Emu")
	pixelgl.Run(run)
}
