package cpu

import (
	"fmt"

	b "github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/io"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type cpu struct {
	pi        int
	maxPi     int
	registers []int
	decoder   decoder
}

// Instance is the interface of the cpu type
type Instance interface {
	get(int) int
	set(int, int)
	Run(b.Instance)
	executeOrAbort(int, func(*int) int) int
}

// New returns a new instance of CPU
func New(registers int, word int, memory int) Instance {
	return &cpu{
		pi:        0,
		maxPi:     int(memory),
		decoder:   newDecoder(word),
		registers: make([]int, registers),
	}
}

func (cpu *cpu) Run(bus b.Instance) {
	go func() {
		for {
			v := <-bus.ReceiveFrom("cpu")
			if v.Signal == b.WRITE {
				bus.SendTo("memory", "cpu", b.READ, []byte{byte(cpu.pi)})

				if cpu.pi == cpu.maxPi {
					cpu.pi = 0
				} else {
					cpu.pi++
				}

				select {
				case vl := <-bus.ReceiveFrom("cpu"):
					instruction := cpu.decoder.decode(vl.Payload)

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

func (cpu *cpu) extractValue(value value) int {
	if value.isRegister {
		return cpu.get(value.value)
	}

	return value.value
}

func (cpu *cpu) mov(register int, params []value) {
	checkLengthOrAbort(params, 1, func() {
		cpu.set(register, cpu.extractValue(params[0]))
	})
}

func (cpu *cpu) add(register int, params []value) {
	checkLengthOrAbort(params, 1, func() {
		v := cpu.get(register)
		cpu.set(register, v+cpu.extractValue(params[0]))
	})
}

func (cpu *cpu) imul(register int, params []value) {
	checkLengthOrAbort(params, 2, func() {
		cpu.set(register, cpu.extractValue(params[0])*cpu.extractValue(params[1]))
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
