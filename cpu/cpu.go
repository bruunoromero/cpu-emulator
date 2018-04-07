package cpu

import (
	"github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type cpu struct {
	pi        int
	registers []int
	decoder   decoder
}

// Instance is the interface of the cpu type
type Instance interface {
	Get(int) int
	Set(int, int)
	Run(bus.Instance)
	executeOrAbort(int, func(*int) int) int
}

// New returns a new instance of CPU
func New(registers int) Instance {
	return &cpu{
		pi:        0,
		decoder:   newDecoder(),
		registers: make([]int, registers),
	}
}

func (cpu *cpu) Run(bus bus.Instance) {

}

func (cpu *cpu) Set(register int, value int) {
	cpu.executeOrAbort(register, func(register *int) int {
		*register = value
		return *register
	})
}

func (cpu *cpu) Get(register int) int {
	return cpu.executeOrAbort(register, func(register *int) int {
		return *register
	})
}

func (cpu *cpu) executeOrAbort(register int, callback func(*int) int) int {
	if len(cpu.registers)-1 < register || register < 0 {
		utils.Abort("Error while accessing register")
	} else {
		return callback(&cpu.registers[register])
	}

	return -1
}
