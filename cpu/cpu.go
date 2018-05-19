package cpu

import (
	b "github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/io"
	"github.com/bruunoromero/cpu-emulator/parser"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type cpu struct {
	pi        int
	maxPi     int
	registers []int
	decoder   parser.Decoder
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
		decoder:   parser.NewDecoder(word),
		registers: make([]int, registers),
	}
}

func (cpu *cpu) Run(bus b.Instance) {
	go func() {
		for {
			v := bus.ReceiveFrom("cpu")
			if v.Signal == b.WRITE {
				bus.SendTo("memory", "cpu", b.READ, []byte{byte(cpu.pi)})

				if cpu.pi == cpu.maxPi-1 {
					cpu.pi = 0
				} else {
					cpu.pi++
				}

				vl := bus.ReceiveFrom("cpu")
				instruction := cpu.decoder.Decode(vl.Payload)

				if instruction.IsRegister {
					switch instruction.Action {
					case io.Inc:
						cpu.add(instruction.location, []parser.Value{value{isRegister: false, value: 1}})
					case io.Add:
						cpu.add(instruction.location, instruction.params)
					case io.Mov:
						cpu.mov(instruction.location, instruction.params)
					case io.Imul:
						cpu.imul(instruction.location, instruction.params)
					}
				}
			}

		}
	}()
}

func checkLengthOrAbort(params []parser.Value, length int, callback func()) {
	if len(params) == length {
		callback()
	} else {
		utils.Abort("Unexpected number of parameters for this action")
	}
}

func (cpu *cpu) extractValue(value parser.Value) int {
	if value.IsRegister {
		return cpu.get(value.Value)
	}

	return value.Value
}

func (cpu *cpu) mov(register int, params []parser.Value) {
	checkLengthOrAbort(params, 1, func() {
		cpu.set(register, cpu.extractValue(params[0]))
	})
}

func (cpu *cpu) add(register int, params []parser.Value) {
	checkLengthOrAbort(params, 1, func() {
		v := cpu.get(register)
		cpu.set(register, v+cpu.extractValue(params[0]))
	})
}

func (cpu *cpu) imul(register int, params []parser.Value) {
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
