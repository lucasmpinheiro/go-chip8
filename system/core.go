package chip8

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
}

// EmulateCycle executes all operations required during a cycle.
func (chip8 *Chip8) EmulateCycle() {
	// Fetch opcode.
	// Decode opcode.
	// Execute opcode.

	// Update timers.
}
