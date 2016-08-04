package core

import (
	"io/ioutil"
	"log"
)

// Chip8 defines the system's main struct.
type Chip8 struct {
	opcode uint16
	memory [0x1000]uint8 // 4kb
	V      [0x10]uint8   // 16 registers.
	I      uint16        // Index register.
	pc     uint16        // Program counter.

	gfx      [64 * 32]uint8 // 64x32 display.
	drawFlag bool           // Indicates if the display should be redrawn or not.

	delayTimer uint8
	soundTimer uint8

	stack [0x10]uint16 // Memory stack.
	sp    uint8        // Stack pointer.

	key [0x10]bool // Key state array.
}

// Initialize sets up memory and registers.
func (chip8 *Chip8) Initialize() {
	chip8.pc = 0x200 // Program counter starts at 0x200.
	chip8.opcode = 0 // Reset current opcode.
	chip8.I = 0      // Reset index register.
	chip8.sp = 0     // Reset stack pointer.

	// Clear display.
	for i := range chip8.gfx {
		chip8.gfx[i] = 0
	}

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
	chip8.opcode = uint16(chip8.memory[chip8.pc])<<8 | uint16(chip8.memory[chip8.pc+1])

	// Decode opcode.
	chip8.DecodeOpcode()

	// Update timers.
	if chip8.delayTimer > 0 {
		chip8.delayTimer--
	}

	if chip8.soundTimer > 0 {
		if chip8.soundTimer == 1 {
			log.Println("BEEP!")
		}
		chip8.soundTimer--
	}
}

// DecodeOpcode decodes the current opcode and performs the right operation.
func (chip8 *Chip8) DecodeOpcode() {
	switch chip8.opcode & 0xF000 {

	case 0x0000:
		switch chip8.opcode & 0x0FFF {

		// 00E0: Clears the screen.
		case 0x00E0:
			for i := range chip8.gfx {
				chip8.gfx[i] = 0
			}

			chip8.drawFlag = true
			chip8.pc += 2

		// 00EE: Returns from a subroutine
		case 0x00EE:
			// TODO: implement this
		default:
			log.Fatalln("Unknown opcode [0x0000]: %X", chip8.opcode)
		}

	// 1NNN: Jumps to address NNN.
	case 0x1000:
		chip8.pc = chip8.opcode & 0x0FFF

	// 2NNN: Calls subroutine at NNN.
	case 0x2000:
		chip8.stack[chip8.sp] = chip8.pc
		chip8.sp++
		chip8.pc = chip8.opcode & 0x0FFF

	// 3XNN: Skips the next instruction if VX equals NN.
	case 0x3000:
		x := (chip8.opcode & 0x0F00) >> 8
		if chip8.V[x] == uint8(chip8.opcode&0x00FF) {
			chip8.pc += 4
		} else {
			chip8.pc += 2
		}

	// 4XNN: Skips the next instruction if VX doesn't equal NN.
	case 0x4000:
		x := (chip8.opcode & 0x0F00) >> 8
		if chip8.V[x] != uint8(chip8.opcode&0x00FF) {
			chip8.pc += 4
		} else {
			chip8.pc += 2
		}

	// 5XY0: Skips the next instruction if VX equals VY.
	case 0x5000:
		x := (chip8.opcode & 0x0F00) >> 8
		y := (chip8.opcode & 0x00F0) >> 4
		if chip8.V[x] == chip8.V[y] {
			chip8.pc += 4
		} else {
			chip8.pc += 2
		}

	// 6XNN:	Sets VX to NN.
	case 0x6000:
		x := (chip8.opcode & 0x0F00) >> 8
		chip8.V[x] = uint8(chip8.opcode & 0x00FF)
		chip8.pc += 2

	// 7XNN:	Adds NN to VX.
	case 0x7000:
		x := (chip8.opcode & 0x0F00) >> 8
		chip8.V[x] += uint8(chip8.opcode & 0x00FF)
		chip8.pc += 2

	case 0x8000:
		switch chip8.opcode & 0x000F {
		case 0x0000:
		case 0x0001:
		case 0x0002:
		case 0x0003:

		// 8XY4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
		case 0x0004:
			x := (chip8.opcode & 0x0F00) >> 8
			y := (chip8.opcode & 0x00F0) >> 4
			if chip8.V[x]+chip8.V[y] > 0xFF {
				chip8.V[0xF] = 1 // carry
			} else {
				chip8.V[0xF] = 0
			}
			chip8.V[x] += chip8.V[y]
			chip8.pc += 2

		case 0x0005:
		case 0x0006:
		case 0x0007:
		case 0x000E:
		default:
			log.Fatalln("Unknown opcode [0x8000]: %X", chip8.opcode)
		}

	// 9XY0: Skips the next instruction if VX doesn't equal VY.
	case 0x9000:
		x := (chip8.opcode & 0x0F00) >> 8
		y := (chip8.opcode & 0x00F0) >> 4
		if chip8.V[x] != chip8.V[y] {
			chip8.pc += 4
		} else {
			chip8.pc += 2
		}

	// ANNN: sets I to address NNN
	case 0xA000:
		chip8.I = chip8.opcode & 0x0FFF
		chip8.pc += 2

	// BNNN: Jumps to the address NNN plus V0.
	case 0xB000:
		chip8.pc = (chip8.opcode & 0x0FFF) + uint16(chip8.V[0])

	case 0xC000:

	// DXYN: Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels. Each row of 8 pixels is read as bit-coded starting from memory location I; I value doesn’t change after the execution of this instruction. As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that doesn’t happen.
	case 0xD000:
		x := (chip8.opcode & 0x0F00) >> 8
		y := (chip8.opcode & 0x00F0) >> 4
		n := chip8.opcode & 0x000F // == sprite height

		var pixels uint8 // Stores a sprite line.

		chip8.V[0xF] = 0 // Reset VF.

		for row := uint16(0); row < n; row++ { // Loop through the sprite rows.
			pixels = chip8.memory[chip8.I+row] // Fetch sprite line from memory.

			for col := uint16(0); col < 8; col++ { // Loop through the sprite columns.
				// Check if the current evaluated pixel is set to 1.
				if (pixels & (0x80 >> col)) != 0 {
					pixelIdx := x + col + ((y + row) * 64) // Get the pixel index.

					// Check if the pixel on the display is set to 1.
					if chip8.gfx[pixelIdx] == 1 {
						// If so, set VF.
						chip8.V[0xF] = 1
					}

					chip8.gfx[pixelIdx] ^= 1 // Update the pixel.
				}
			}
		}

		chip8.drawFlag = true
		chip8.pc += 2

	case 0xE000:
		switch chip8.opcode & 0x00FF {
		case 0x009E:
		case 0x00A1:
		default:
			log.Fatalln("Unknown opcode [0xE000]: %X", chip8.opcode)
		}

	case 0xF000:
		switch chip8.opcode & 0x00FF {

		// FX07: Sets VX to the value of the delay timer.
		case 0x0007:
			x := (chip8.opcode & 0x0F00) >> 8
			chip8.V[x] = chip8.delayTimer
			chip8.pc += 2

		case 0x000A:

		// FX15: Sets the delay timer to VX.
		case 0x0015:
			x := (chip8.opcode & 0x0F00) >> 8
			chip8.delayTimer = chip8.V[x]
			chip8.pc += 2

		// FX15: Sets the sound timer to VX.
		case 0x0018:
			x := (chip8.opcode & 0x0F00) >> 8
			chip8.soundTimer = chip8.V[x]
			chip8.pc += 2

		// FX1E: Adds VX to I.
		case 0x001E:
			x := (chip8.opcode & 0x0F00) >> 8
			if chip8.I+uint16(chip8.V[x]) > 0xFFF { // Check for range overflow.
				chip8.V[0xF] = 1
			} else {
				chip8.V[0xF] = 0
			}
			chip8.I += uint16(chip8.V[x])
			chip8.pc += 2

		// FX29: Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
		case 0x0029:
			x := (chip8.opcode & 0x0F00) >> 8
			chip8.I = uint16(chip8.V[x] * 0x5)
			chip8.pc += 2

		// FX33: Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.)
		case 0x0033:
			x := (chip8.opcode & 0x0F00) >> 8
			chip8.memory[chip8.I] = chip8.V[x] / 100
			chip8.memory[chip8.I+1] = (chip8.V[x] / 10) % 10
			chip8.memory[chip8.I+2] = chip8.V[x] % 10
			chip8.pc += 2

		// FX55: Stores V0 to VX (including VX) in memory starting at address I.
		case 0x0055:
			x := (chip8.opcode & 0x0F00) >> 8
			if chip8.I+x > uint16(len(chip8.memory)) {
				log.Fatal("Error: end of memory.")
			} else {
				for i := uint16(0); i <= x; i++ {
					chip8.memory[chip8.I+i] = chip8.V[i]
				}
			}

		// FX65: Fills V0 to VX (including VX) with values from memory starting at address I.
		case 0x0065:
			x := (chip8.opcode & 0x0F00) >> 8
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
