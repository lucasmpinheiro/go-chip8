package core

import (
	"io/ioutil"
	"log"
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

// DecodeOpcode decodes the current opcode and performs the right operation.
func (chip8 *Chip8) DecodeOpcode() {
	switch chip8.opcode & 0xF000 {

	case 0x0000:
		switch chip8.opcode & 0x0FFF {
		case 0x00E0: // 0x00E0: Clears the screen.
			// TODO: implement this
		case 0x00EE: // 0x00EE: Returns from a subroutine
			// TODO: implement this
		default:
			log.Fatalln("Unknown opcode [0x0000]: %X", chip8.opcode)
		}

	case 0x1000: // 0x1NNN: Jumps to address NNN.
		chip8.pc = chip8.opcode & 0x0FFF

	case 0x2000: // 0x2NNN: Calls subroutine at NNN.
		chip8.stack[chip8.sp] = chip8.pc
		chip8.sp++
		chip8.pc = chip8.opcode & 0x0FFF

	case 0x3000: // 0x3XNN: Skips the next instruction if VX equals NN.
		x := (chip8.opcode & 0x0F00) / 100
		if chip8.V[x] == uint8(chip8.opcode&0x00FF) {
			chip8.pc += 4
		} else {
			chip8.pc += 2
		}

	case 0x4000: // 0x4XNN: Skips the next instruction if VX doesn't equal NN.
		x := (chip8.opcode & 0x0F00) / 100
		if chip8.V[x] != uint8(chip8.opcode&0x00FF) {
			chip8.pc += 4
		} else {
			chip8.pc += 2
		}

	case 0x5000: // 0x5XY0: Skips the next instruction if VX equals VY.
		x := (chip8.opcode & 0x0F00) / 100
		y := (chip8.opcode & 0x00F0) / 10
		if chip8.V[x] == chip8.V[y] {
			chip8.pc += 4
		} else {
			chip8.pc += 2
		}

	case 0x6000: // 0x6XNN:	Sets VX to NN.
		x := (chip8.opcode & 0x0F00) / 100
		chip8.V[x] = uint8(chip8.opcode & 0x00FF)
		chip8.pc += 2

	case 0x7000: // 0x7XNN:	Adds NN to VX.
		x := (chip8.opcode & 0x0F00) / 100
		chip8.V[x] += uint8(chip8.opcode & 0x00FF)
		chip8.pc += 2

	case 0x8000:
		switch chip8.opcode & 0x000F {
		case 0x0000:
		case 0x0001:
		case 0x0002:
		case 0x0003:
		case 0x0004:
		case 0x0005:
		case 0x0006:
		case 0x0007:
		case 0x000E:
		default:
			log.Fatalln("Unknown opcode [0x8000]: %X", chip8.opcode)
		}

	case 0x9000: // 0x9XY0: Skips the next instruction if VX doesn't equal VY.
		x := (chip8.opcode & 0x0F00) / 100
		y := (chip8.opcode & 0x00F0) / 10
		if chip8.V[x] != chip8.V[y] {
			chip8.pc += 4
		} else {
			chip8.pc += 2
		}

	case 0xA000: // ANNN: sets I to address NNN
		chip8.I = chip8.opcode & 0x0FFF
		chip8.pc += 2

	case 0xB000: // 0xBNNN: Jumps to the address NNN plus V0.
		chip8.pc = (chip8.opcode & 0x0FFF) + uint16(chip8.V[0])

	case 0xC000:

	case 0xD000:

	case 0xE000:
		switch chip8.opcode & 0x00FF {
		case 0x009E:
		case 0x00A1:
		default:
			log.Fatalln("Unknown opcode [0xE000]: %X", chip8.opcode)
		}

	case 0xF000:
		switch chip8.opcode & 0x00FF {
		case 0x0007: // 0xFX07: Sets VX to the value of the delay timer.
			x := (chip8.opcode & 0x0F00) / 100
			chip8.V[x] = chip8.delayTimer
			chip8.pc += 2

		case 0x000A:
		case 0x0015: // 0xFX15: Sets the delay timer to VX.
			x := (chip8.opcode & 0x0F00) / 100
			chip8.delayTimer = chip8.V[x]
			chip8.pc += 2
		case 0x0018: // 0xFX15: Sets the sound timer to VX.
			x := (chip8.opcode & 0x0F00) / 100
			chip8.soundTimer = chip8.V[x]
			chip8.pc += 2
		case 0x001E: // 0xFX1E: Adds VX to I.
			x := (chip8.opcode & 0x0F00) / 100
			if chip8.I+uint16(chip8.V[x]) > 0xFFF { // Check for range overflow.
				chip8.V[0xF] = 1
			} else {
				chip8.V[0xF] = 0
			}
			chip8.I += uint16(chip8.V[x])
			chip8.pc += 2
		case 0x0029:
		case 0x0033:
		case 0x0055: // 0xFX55: Stores V0 to VX (including VX) in memory starting at address I.
			x := (chip8.opcode & 0x0F00) / 100
			if chip8.I+x > uint16(len(chip8.memory)) {
				log.Fatal("Error: end of memory.")
			} else {
				for i := uint16(0); i <= x; i++ {
					chip8.memory[chip8.I+i] = chip8.V[i]
				}
			}

		case 0x0065: // 0xFX65: Fills V0 to VX (including VX) with values from memory starting at address I.
			x := (chip8.opcode & 0x0F00) / 100
			if chip8.I+x > uint16(len(chip8.memory)) {
				log.Fatal("Error: end of memory.")
			} else {
				for i := uint16(0); i <= x; i++ {
					chip8.V[i] = chip8.memory[chip8.I+i]
				}
			}

		default:
			log.Fatalln("Unknown opcode [0xF000]: %X", chip8.opcode)
		}

	default:
		log.Fatalln("Unknown opcode: %X", chip8.opcode)
	}
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
