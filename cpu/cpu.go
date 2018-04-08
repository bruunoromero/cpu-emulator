package cpu

import (
	"fmt"

	"github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/io"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type cpu struct {
	pi        int
	registers []int
}

// Instance is the interface of the cpu type
type Instance interface {
	get(int) int
	set(int, int)
	Run(bus.Instance)
	executeOrAbort(int, func(*int) int) int
}

// New returns a new instance of CPU
func New(registers int) Instance {
	return &cpu{
		pi:        0,
		registers: make([]int, registers),
	}
}

func (cpu *cpu) Run(bus bus.Instance) {
	go func() {
		for {
			value := bus.ReceiveAtCPU()
			instruction := decode(value)

			if instruction.isRegister {
				switch instruction.action {
				case io.Inc:
					cpu.add(instruction.location, []int{1})
				case io.Add:
					cpu.add(instruction.location, instruction.params)
				case io.Mov:
					cpu.mov(instruction.location, instruction.params)
				case io.Imul:
					cpu.imul(instruction.location, instruction.params)
				}
			}

			fmt.Println(cpu.registers)

		}
	}()
}

func checkLengthOrAbort(params []int, length int, callback func()) {
	if len(params) == length {
		callback()
	} else {
		utils.Abort("Unexpected number of parameters for this action")
	}
}

func (cpu *cpu) mov(register int, params []int) {
	checkLengthOrAbort(params, 1, func() {
		cpu.set(register, params[0])
	})
}

func (cpu *cpu) add(register int, params []int) {
	checkLengthOrAbort(params, 1, func() {
		v := cpu.get(register)
		cpu.set(register, v+params[0])
	})
}

func (cpu *cpu) imul(register int, params []int) {
	checkLengthOrAbort(params, 2, func() {
		cpu.set(register, params[0]*params[1])
	})
}

func (cpu *cpu) set(register int, value int) {
	cpu.executeOrAbort(register, func(register *int) int {
		*register = value
		return *register
	})
}

func (cpu *cpu) get(register int) int {
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
