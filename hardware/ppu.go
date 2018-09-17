package hardware


type Ppu struct {
	CPU *Cpu
	Memory [0x4000]byte
}