package core

import (
	"io/ioutil"
)

// Chip8 defines the system's main struct.
type Chip8 struct {
	opcode uint16
	memory [0x1000]uint8 // 4kb
	V      [0x10]uint8   // 16 registers
	I      uint16        // Index register
	pc     uint16        // Program counter

	gfx [64 * 32]bool // 64x32 display

	delayTimer uint8
	soundTimer uint8

	stack [0x10]uint16 // Memory stack
	sp    uint8        // Stack pointer

	key [0x10]bool // Key state array
}

// Initialize sets up memory and registers.
func (chip8 *Chip8) Initialize() {
	chip8.pc = 0x200 // Program counter starts at 0x200
	chip8.opcode = 0 // Reset current opcode
	chip8.I = 0      // Reset index register
	chip8.sp = 0     // Reset stack pointer

	// TODO: Clear display.

	// Clear stack.
	for i := range chip8.stack {
		chip8.stack[i] = 0
	}

	// Clear registers V0-VF.
	for i := range chip8.V {
		chip8.V[i] = 0
	}

	// Clear memory.
	for i := range chip8.memory {
		chip8.memory[i] = 0
	}

	// Load fontset.
	for i, font := range chip8Fontset {
		chip8.memory[i] = font
	}

	// Reset timers.
	chip8.delayTimer = 0
	chip8.soundTimer = 0
}

// EmulateCycle executes all operations required during a cycle.
func (chip8 *Chip8) EmulateCycle() {
	// Fetch opcode.
	// Decode opcode.
	// Execute opcode.

	// Update timers.
}

// LoadGame loads the program stored on a file to the memory.
func (chip8 *Chip8) LoadGame(filename string) {
	// Read file data.
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	// Load program into memory.
	for i, value := range data {
		if i+0x200 < 0x1000 {
			chip8.memory[i+0x200] = value
		} else {
			log.Fatal("Error: end of memory.")
		}
	}
}

var chip8Fontset [80]uint8 = [80]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80} // F
