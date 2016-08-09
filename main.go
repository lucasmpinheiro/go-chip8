package main

import (
	"github.com/lucasmpinheiro/go-chip8/system"
)

func main() {
	// Setup render system and register input callbacks.

	// Initialize the system and load the game into memory.
	chip8 := new(system.Chip8)
	chip8.Initialize()
	chip8.LoadGame("games/pong")

	// Emulation loop.
	for {
		// Emulate one cycle.
		chip8.EmulateCycle()

		// If the draw flag is set, update the screen.
		if chip8.DrawFlag() {
			// TODO: Draw graphics.
		}

		// Store key press state.
	}
}
