package cpu

import (
	"fmt"
	"time"

	"github.com/bradfitz/slice"

	b "github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/parser"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type cpu struct {
	pi             int
	executionIndex int
	registers      []int
	messageQueue   []parser.Msg
	decoder        parser.Decoder
	executionMap   map[int][]parser.Msg
}

// Instance is the interface of the cpu type
type Instance interface {
	Run(b.Instance)
	get(parser.Parameter) int
	set(parser.Parameter, int)
	executeOrAbort(int, func(*int) int) int
}

// New returns a new instance of CPU
func New(registers int, word int, memory int) Instance {
	return &cpu{
		pi:             -1,
		executionIndex: 0,
		messageQueue:   make([]parser.Msg, 0),
		registers:      make([]int, registers),
		decoder:        parser.NewDecoder(word),
		executionMap:   make(map[int][]parser.Msg),
	}
}

func (cpu *cpu) Run(bus b.Instance) {
	go func() {
		for {
			data := bus.ReceiveFrom("cpuData")
			address := bus.ReceiveFrom("cpuAddress")
			instructions := bus.ReceiveFrom("cpuInstruction")

			messages := cpu.decoder.GetMessagesWithQueue(address.Payload, data.Payload, instructions.Payload, &cpu.messageQueue)

			slice.Sort(messages, func(i int, j int) bool {
				return messages[i][0].Key < messages[j][0].Key
			})

			for _, message := range messages {
				if len(message) > 0 {
					msg := message[0]

					if msg.Signal == b.READ {
						if msg.Key != cpu.pi+1 {
							cpu.messageQueue = append(message, cpu.messageQueue...)
							continue
						}
						cpu.pi++
						bus.SendTo("memory", "cpu", b.READ, []parser.Msg{parser.Msg{Key: int(msg.Value), Index: 0, Lenght: 0, Type: parser.MEMORY, Value: msg.Value}})
					} else {
						cpu.executionMap[msg.Key] = message
						if cpu.executionMap[cpu.executionIndex] == nil {
							continue
						}

						cpu.executionIndex++
						instruction := cpu.decoder.Decode(message)

						// 1. CHECK IF ALL VALUES ARE NOT FROM MEMORY, OTHERWISE GET VELUES FROM MEMORY
						// AND DECREMENT EXECUTIONINDEX TO REEXECUTE THE SAME INSTRUCTION WITH THE RAW VALUE
						// PICKED FROM MEMORY

						// 2. IF LOCATION IS MEMORY SEND A MESSAGE TO MEMORY, IF NOT THEN USE THE CODE BELOW

						switch instruction.Action {
						case parser.Inc:
							cpu.add(instruction.Location, []parser.Parameter{parser.Parameter{Type: parser.LITERAL, Value: 1}})
						case parser.Add:
							cpu.add(instruction.Location, instruction.Parameters)
						case parser.Mov:
							cpu.mov(instruction.Location, instruction.Parameters)
						case parser.Imul:
							cpu.imul(instruction.Location, instruction.Parameters)
						}
					}
				}
			}

			fmt.Println(cpu.registers)
			time.Sleep(time.Second / 2)
		}
	}()
}

func checkLengthOrAbort(params []parser.Parameter, length int, callback func()) {
	if len(params) == length {
		callback()
	} else {
		utils.Abort("Unexpected number of parameters for this action")
	}
}

func (cpu *cpu) extractValue(value parser.Parameter) int {
	if value.Type == parser.REGISTER {
		return cpu.get(value)
	} else if value.Type == parser.LITERAL {
		return value.Value
	}

	fmt.Println(value.Type, value.Value)
	utils.Abort("undefined type")
	return 0
}

func (cpu *cpu) extractValueFromMemory(value parser.Parameter) {

}

func (cpu *cpu) mov(register parser.Parameter, params []parser.Parameter) {
	checkLengthOrAbort(params, 1, func() {
		value := params[0]
		cpu.set(register, cpu.extractValue(value))

	})
}

func (cpu *cpu) add(register parser.Parameter, params []parser.Parameter) {
	checkLengthOrAbort(params, 1, func() {
		v := cpu.get(register)
		value := params[0]
		cpu.set(register, v+cpu.extractValue(value))
	})
}

func (cpu *cpu) imul(register parser.Parameter, params []parser.Parameter) {
	checkLengthOrAbort(params, 2, func() {
		value1 := params[0]
		value2 := params[1]
		cpu.set(register, cpu.extractValue(value1)*cpu.extractValue(value2))
	})
}

func (cpu *cpu) set(location parser.Parameter, value int) {
	cpu.executeOrAbort(location.Value, func(register *int) int {
		*register = value
		return *register
	})
}

func (cpu *cpu) get(location parser.Parameter) int {
	return cpu.executeOrAbort(location.Value, func(register *int) int {
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
