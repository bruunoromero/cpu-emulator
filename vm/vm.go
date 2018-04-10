package vm

import (
	"sync"

	"github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/cpu"
	"github.com/bruunoromero/cpu-emulator/io"
	"github.com/bruunoromero/cpu-emulator/memory"
)

var once sync.Once

// Start initiates the Von Neumann loop
func Start(registers []string, wordLength int, memoryLength int) {
	once.Do(func() {
		bus := bus.New()
		io := io.New(registers, wordLength)
		cpu := cpu.New(len(registers), wordLength)
		memory := memory.New(memoryLength, wordLength)

		bus.MakeChannel("io")
		bus.MakeChannel("cpu")
		bus.MakeChannel("memory")

		memory.Run(bus)
		cpu.Run(bus)
		io.Run(bus)
	})
}
