package cpu

import (
	"fmt"

	b "github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/io"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type cpu struct {
	pi        int8
	registers []int8
}

// Instance is the int8erface of the cpu type
type Instance interface {
	get(int8) int8
	set(int8, int8)
	Run(b.Instance)
	executeOrAbort(int8, func(*int8) int8) int8
}

// New returns a new instance of CPU
func New(registers int, word int) Instance {
	return &cpu{
		pi:        0,
		registers: make([]int8, registers),
	}
}

func (cpu *cpu) Run(bus b.Instance) {
	go func() {
		for {
			v := <-bus.ReceiveFrom("cpu")
			if v.Signal == b.WRITE {
				bus.SendTo("memory", "cpu", b.READ, []int8{cpu.pi})
				cpu.pi++
				select {
				case vl := <-bus.ReceiveFrom("cpu"):
					instruction := decode(vl.Payload)

					if instruction.isRegister {
						switch instruction.action {
						case io.Inc:
							cpu.add(instruction.location, []value{value{isRegister: false, value: 1}})
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
			}

		}
	}()
}

func checkLengthOrAbort(params []value, length int, callback func()) {
	if len(params) == length {
		callback()
	} else {
		utils.Abort("Unexpected number of parameters for this action")
	}
}

func (cpu *cpu) extractValue(value value) int8 {
	if value.isRegister {
		return cpu.get(value.value)
	}

	return value.value
}

func (cpu *cpu) mov(register int8, params []value) {
	checkLengthOrAbort(params, 1, func() {
		cpu.set(register, cpu.extractValue(params[0]))
	})
}

func (cpu *cpu) add(register int8, params []value) {
	checkLengthOrAbort(params, 1, func() {
		v := cpu.get(register)
		cpu.set(register, v+cpu.extractValue(params[0]))
	})
}

func (cpu *cpu) imul(register int8, params []value) {
	checkLengthOrAbort(params, 2, func() {
		cpu.set(register, cpu.extractValue(params[0])*cpu.extractValue(params[1]))
	})
}

func (cpu *cpu) set(register int8, value int8) {
	cpu.executeOrAbort(register, func(register *int8) int8 {
		*register = value
		return *register
	})
}

func (cpu *cpu) get(register int8) int8 {
	return cpu.executeOrAbort(register, func(register *int8) int8 {
		return *register
	})
}

func (cpu *cpu) executeOrAbort(register int8, callback func(*int8) int8) int8 {
	if len(cpu.registers)-1 < int(register) || int(register) < 0 {
		utils.Abort("Error while accessing register")
	} else {
		return callback(&cpu.registers[register])
	}

	return -1
}
