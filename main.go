package main

import (
	"github.com/faiface/pixel/pixelgl"
	"nes-emu/cmd"
)

func run() {
	cmd.Execute()
}

func main() {
	pixelgl.Run(run)
}
