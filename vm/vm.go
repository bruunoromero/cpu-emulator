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
func Start(registers []string) {
	once.Do(func() {
		bus := bus.New()
		memory := memory.New()
		io := io.New(registers)
		cpu := cpu.New(len(registers))

		memory.Run(bus)
		cpu.Run(bus)
		io.Run(bus)
	})
}
