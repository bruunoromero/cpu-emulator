package vm

import (
	"sync"

	"github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/cpu"
	"github.com/bruunoromero/cpu-emulator/io"
	"github.com/bruunoromero/cpu-emulator/memory"
	"github.com/bruunoromero/cpu-emulator/parser"
)

var once sync.Once

// Start initiates the Von Neumann loop
func Start(registers []string, busLength int, wordLength int, memoryLength int, frequency int) {
	once.Do(func() {
		encoder := parser.NewEncoder(registers, wordLength)

		io := io.New(encoder)
		bus := bus.New(frequency, busLength)
		memory := memory.New(memoryLength, wordLength, frequency)
		cpu := cpu.New(len(registers), wordLength, (memoryLength/(wordLength/8))/4, frequency, memory, encoder)

		bus.MakeChannel("cpu")
		bus.MakeChannel("memory")

		bus.Run()
		memory.Run(bus)
		cpu.Run(bus)
		io.Run(bus)
	})
}
