# Go Chip8 Emulator

This is a project to create a simple Chip 8 emulator written in [Go](https://golang.org/).

It's largely based on the [tutorial by Laurance Muller](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/).

## TODO

### Opcodes
- [ ] 0x0000
    - [x] 0x00E0
    - [ ] 0x00EE
- [x] 0x1NNN
- [x] 0x2NNN
- [x] 0x3XNN
- [x] 0x4XNN
- [x] 0x5XY0
- [x] 0x6XNN
- [x] 0x7XNN
- [ ] 0x8000
    - [ ] 0x8XY0
    - [ ] 0x8XY1
    - [ ] 0x8XY2
    - [ ] 0x8XY3
    - [x] 0x8XY4
    - [ ] 0x8XY5
    - [ ] 0x8XY6
    - [ ] 0x8XY7
    - [ ] 0x8XYE
- [x] 0x9XY0
- [x] 0xANNN
- [x] 0xBNNN
- [ ] 0xCXNN
- [x] 0xDXYN
- [ ] 0xE000
    - [ ] 0xEX9E
    - [ ] 0xEXA1
- [ ] 0xF000:
    - [x] 0xFX07
    - [ ] 0xFX0A
    - [x] 0xFX15
    - [x] 0xFX18
    - [x] 0xFX1E
    - [x] 0xFX29
    - [x] 0xFX33
    - [x] 0xFX55
    - [x] 0xFX65
