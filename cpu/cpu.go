package cpu

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bruunoromero/cpu-emulator/memory"

	"github.com/bradfitz/slice"
	b "github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/parser"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type cacheable struct {
	value  int
	access int
}

type cpu struct {
	pi                      int
	memoryOffest            int
	executionIndex          int
	wordLenth               int
	frequency               int
	lastConditionalIndex    int
	loopExecuting           int
	isLooping               bool
	isBuildingLoop          bool
	isWaitingForConditional bool
	registers               []int
	loops                   map[int]loop
	messageQueue            []parser.Msg
	encoder                 parser.Encoder
	decoder                 parser.Decoder
	memory                  memory.Instance
	cache                   map[int]cacheable
	executionMap            map[int][]parser.Msg
}

// Instance is the interface of the cpu type
type Instance interface {
	Run(b.Instance)
	get(parser.Parameter) int
	set(parser.Parameter, int)
	executeOrAbort(int, func(*int) int) int
}

type loop struct {
	condition    parser.Action
	ifTrue       parser.Action
	ifFalse      parser.Action
	instructions []parser.Action
}

// New returns a new instance of CPU
func New(registers int, word int, memory int, frequency int, memoryI memory.I, encoder parser.Encoder) Instance {
	return &cpu{
		lastConditionalIndex:    0,
		executionIndex:          0,
		loopExecuting:           -1,
		pi:                      -1,
		isBuildingLoop:          false,
		isLooping:               false,
		wordLenth:               word,
		isWaitingForConditional: false,
		memory:                  memoryI,
		encoder:                 encoder,
		frequency:               frequency,
		memoryOffest:            (memory / 2) - 1,
		loops:                   make(map[int]loop),
		messageQueue:            make([]parser.Msg, 0),
		registers:               make([]int, registers),
		cache:                   make(map[int]cacheable),
		decoder:                 parser.NewDecoder(word),
		executionMap:            make(map[int][]parser.Msg),
	}
}

func mapSlice(vs []parser.Msg, f func(parser.Msg) byte) []byte {
	vsm := make([]byte, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func getValue(msg parser.Msg) byte {
	return msg.Value
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

						instruction := cpu.decoder.Decode(cpu.executionMap[cpu.executionIndex])
						cpu.executeInstruction(instruction)
						cpu.executionIndex++
					}
				}
			}

			time.Sleep(time.Second / (time.Duration(cpu.frequency) * 4))
		}
	}()
}

func (cpu *cpu) shouldExecute(instruction parser.Action) bool {
	if isConditional(instruction.Action) {
		loop := cpu.loops[cpu.loopExecuting]
		loop.condition = instruction
		cpu.isWaitingForConditional = true
		cpu.lastConditionalIndex = instruction.Key
		cpu.loops[cpu.loopExecuting] = loop
		return false
	} else if cpu.isWaitingForConditional {
		if instruction.Key == cpu.lastConditionalIndex+1 {
			loop := cpu.loops[cpu.loopExecuting]
			loop.ifTrue = instruction
			cpu.loops[cpu.loopExecuting] = loop
			return false
		} else if instruction.Key == cpu.lastConditionalIndex+2 {
			loop := cpu.loops[cpu.loopExecuting]
			cpu.isBuildingLoop = false
			cpu.isWaitingForConditional = false
			loop.ifFalse = instruction
			cpu.loops[cpu.loopExecuting] = loop
			cpu.branchLoop()
			return false
		}
	}

	return true
}

func (cpu *cpu) cacheInstructions(instruction parser.Action) {
	if cpu.isBuildingLoop {
		cached := parser.Action{
			Key:        instruction.Key,
			Signal:     instruction.Signal,
			Origin:     instruction.Origin,
			Action:     instruction.Action,
			Location:   instruction.Location,
			Parameters: make([]parser.Parameter, 0),
		}

		for _, parameter := range instruction.Parameters {
			cached.Parameters = append(cached.Parameters, parser.Parameter{
				Type:  parameter.Type,
				Value: parameter.Value,
			})
		}

		loop := cpu.loops[cpu.loopExecuting]
		loop.instructions = append(loop.instructions, cached)
		cpu.loops[cpu.loopExecuting] = loop
	}
}

func (cpu *cpu) syncCache() {
	cpu.lfu()
}

func (cpu *cpu) lfu() {
	for position, cache := range cpu.cache {
		if cache.access >= 5 {
			cpu.writeToMemory(position, cache.value)
			cache.access = 0
			cpu.cache[position] = cache
		}
	}
}

func (cpu *cpu) executeInstruction(instruction parser.Action) {
	if cpu.shouldExecute(instruction) {
		cpu.cacheInstructions(instruction)

		for index, parameter := range instruction.Parameters {
			if parameter.Type == parser.MEMORY {
				instruction.Parameters[index] = cpu.resolveParameter(parameter)
			}
		}

		switch instruction.Location.Type {
		case parser.MEMORY:
			if cpu.isLooping {
				cpu.executeOnCache(instruction)
			} else {
				cpu.executeOnMemory(instruction)
			}
		default:
			cpu.executeOnRegister(instruction)
		}
	}
}

func (cpu *cpu) executeOnRegister(instruction parser.Action) {
	switch instruction.Action {
	case parser.Inc:
		cpu.add(instruction.Location, []parser.Parameter{parser.Parameter{Type: parser.LITERAL, Value: 1}})
		fmt.Println("registers", cpu.registers)
	case parser.Add:
		cpu.add(instruction.Location, instruction.Parameters)
		fmt.Println("registers", cpu.registers)
	case parser.Mov:
		cpu.mov(instruction.Location, instruction.Parameters)
		fmt.Println("registers", cpu.registers)
	case parser.Imul:
		cpu.imul(instruction.Location, instruction.Parameters)
		fmt.Println("registers", cpu.registers)
	case parser.Label:
		cpu.label(instruction.Location)
	case parser.Jump:
		cpu.jump(instruction.Location)
	case parser.NULL:
		cpu.null()
	}
}

func (cpu *cpu) executeOnCache(instruction parser.Action) {
	switch instruction.Action {
	case parser.Inc:
		cpu.incOnCache(instruction.Location)
	case parser.Add:
		cpu.addOnCache(instruction.Location, instruction.Parameters)
	case parser.Mov:
		cpu.movOnCache(instruction.Location, instruction.Parameters)
	case parser.Imul:
		cpu.imulOnCache(instruction.Location, instruction.Parameters)
	}
}

func isConditional(action int) bool {
	return action == parser.EQ ||
		action == parser.GT ||
		action == parser.LT ||
		action == parser.GTEQ ||
		action == parser.LTEQ
}

func (cpu *cpu) branchLoop() {
	loop := cpu.loops[cpu.loopExecuting]
	isTrue := cpu.executeCondition()

	cpu.isLooping = false

	if isTrue {
		cpu.executeInstruction(loop.ifTrue)
	} else {
		cpu.executeInstruction(loop.ifFalse)
	}
}

func (cpu *cpu) executeLoop() {
	loop := cpu.loops[cpu.loopExecuting]

	for _, instruction := range loop.instructions {
		cpu.syncCache()
		cpu.executeInstruction(instruction)
	}

	cpu.branchLoop()
}

func (cpu *cpu) executeCondition() bool {
	condition := cpu.loops[cpu.loopExecuting].condition
	left := condition.Location
	right := condition.Parameters[0]

	leftValue := cpu.extractValue(left)
	rightValue := cpu.extractValue(right)

	if condition.Action == parser.EQ {
		return leftValue == rightValue
	} else if condition.Action == parser.GT {
		return leftValue > rightValue
	} else if condition.Action == parser.LT {
		return leftValue < rightValue
	} else if condition.Action == parser.GTEQ {
		return leftValue >= rightValue
	} else if condition.Action == parser.LTEQ {
		return leftValue <= rightValue
	}

	return false
}

func (cpu *cpu) executeOnMemory(instruction parser.Action) {
	switch instruction.Action {
	case parser.Inc:
		cpu.incOnMemory(instruction.Location)
	case parser.Add:
		cpu.addOnMemory(instruction.Location, instruction.Parameters)
	case parser.Mov:
		cpu.movOnMemory(instruction.Location, instruction.Parameters)
	case parser.Imul:
		cpu.imulOnMemory(instruction.Location, instruction.Parameters)
	}
}

func (cpu *cpu) getFirstChunkOfType(msgs []parser.Msg, t int) []parser.Msg {
	chunks := cpu.decoder.MakeChunks(msgs)

	for _, chunk := range chunks {
		if chunk[0].Type == t {
			return chunk
		}
	}

	return nil
}

func findFirstUnresolved(params []parser.Parameter) *parser.Parameter {
	for _, param := range params {
		if param.Type == parser.MEMORY {
			return &param
		}
	}

	return nil
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
	} else if value.Type == parser.MEMORY {
		if cpu.isLooping {
			if cacheValue, ok := cpu.cache[value.Value]; ok {
				cacheValue.access++
				cpu.cache[value.Value] = cacheValue

				return cacheValue.value
			}

			v := value.Value

			cpu.cache[v] = cacheable{
				access: 1,
				value:  v,
			}

			return v
		}

		return cpu.resolveParameter(value).Value
	}

	utils.Abort("undefined type")
	return 0
}

func (cpu *cpu) jump(label parser.Parameter) {
	cpu.isLooping = true
	cpu.isBuildingLoop = false
	cpu.loopExecuting = cpu.extractValue(label)
	cpu.executeLoop()
}

func (cpu *cpu) null() {
	cpu.isLooping = false
	cpu.loopExecuting = -1
	cpu.syncCache()
}

func (cpu *cpu) label(label parser.Parameter) {
	loopIndex := cpu.extractValue(label)
	cpu.isLooping = true
	cpu.isBuildingLoop = true
	cpu.loopExecuting = loopIndex
	cpu.loops[loopIndex] = loop{
		instructions: make([]parser.Action, 0),
	}
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

func (cpu *cpu) movOnMemory(memory parser.Parameter, params []parser.Parameter) {
	checkLengthOrAbort(params, 1, func() {
		value := cpu.extractValue(params[0])
		message := cpu.encoder.MapParams([]string{strconv.Itoa(value)})

		fmt.Println("mov on memory position: ", memory.Value+cpu.memoryOffest, "; value: ", message)
		cpu.memory.Write(byte(memory.Value+cpu.memoryOffest), message)
	})
}

func (cpu *cpu) addOnMemory(memory parser.Parameter, params []parser.Parameter) {
	checkLengthOrAbort(params, 1, func() {
		memoryValue := cpu.resolveParameter(memory)
		v := cpu.extractValue(params[0])
		value := v + memoryValue.Value
		message := cpu.encoder.MapParams([]string{strconv.Itoa(value)})

		fmt.Println("add on memory position: ", memory.Value+cpu.memoryOffest, "; value: ", message)
		cpu.memory.Write(byte(memory.Value+cpu.memoryOffest), message)
	})
}

func (cpu *cpu) imulOnMemory(memory parser.Parameter, params []parser.Parameter) {
	checkLengthOrAbort(params, 2, func() {
		v0 := cpu.extractValue(params[0])
		v1 := cpu.extractValue(params[1])
		value := v0 * v1
		message := cpu.encoder.MapParams([]string{strconv.Itoa(value)})

		fmt.Println("imul on memory position: ", memory.Value+cpu.memoryOffest, "; value: ", message)
		cpu.memory.Write(byte(memory.Value+cpu.memoryOffest), message)
	})
}

func (cpu *cpu) incOnMemory(memory parser.Parameter) {
	memoryValue := cpu.resolveParameter(memory)
	value := memoryValue.Value + 1
	message := cpu.encoder.MapParams([]string{strconv.Itoa(value)})
	fmt.Println("inc on memory position: ", memory.Value+cpu.memoryOffest, "; value: ", message)
	cpu.memory.Write(byte(memory.Value+cpu.memoryOffest), message)
}

func (cpu *cpu) movOnCache(cache parser.Parameter, params []parser.Parameter) {
	checkLengthOrAbort(params, 1, func() {
		cpu.extractValue(cache)

		value := cpu.extractValue(params[0])

		cacheEl := cpu.cache[cache.Value]
		cacheEl.value = value
		cpu.cache[cache.Value] = cacheEl
		fmt.Println("mov on cache position: ", cache.Value, "; value: ", cacheEl.value)
	})
}

func (cpu *cpu) addOnCache(cache parser.Parameter, params []parser.Parameter) {
	checkLengthOrAbort(params, 1, func() {
		value := cpu.extractValue(cache) + cpu.extractValue(params[0])

		cacheEl := cpu.cache[cache.Value]
		cacheEl.value = value
		cpu.cache[cache.Value] = cacheEl

		fmt.Println("add on cache position: ", cache.Value, "; value: ", cacheEl.value)
	})
}

func (cpu *cpu) imulOnCache(cache parser.Parameter, params []parser.Parameter) {
	checkLengthOrAbort(params, 2, func() {
		cpu.extractValue(cache)

		fmt.Println(params[0].Type == parser.LITERAL)
		v0 := cpu.extractValue(params[0])
		v1 := cpu.extractValue(params[1])
		value := v0 * v1

		cacheEl := cpu.cache[cache.Value]
		cacheEl.value = value
		cpu.cache[cache.Value] = cacheEl
		// fmt.Println("imul on cache position: ", cache.Value, "; value: ", cacheEl.value)
	})
}

func (cpu *cpu) incOnCache(cache parser.Parameter) {
	value := cpu.extractValue(cache) + 1

	cacheEl := cpu.cache[cache.Value]
	cacheEl.value = value
	cpu.cache[cache.Value] = cacheEl

	fmt.Println("inc on cache position: ", cache.Value, "; value: ", cacheEl.value)
}

func (cpu *cpu) writeToMemory(position int, value int) {
	message := cpu.encoder.MapParams([]string{strconv.Itoa(value)})
	cpu.memory.Write(byte(position+cpu.memoryOffest), message)

	fmt.Println("write on memory position: ", position+cpu.memoryOffest, "; value: ", message)
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

func (cpu *cpu) resolveParameter(parameter parser.Parameter) parser.Parameter {
	msg := cpu.memory.Read(byte(parameter.Value + cpu.memoryOffest))
	bytes := mapSlice(msg, getValue)
	value := utils.FromBytes(cpu.wordLenth, bytes)

	return parser.Parameter{
		Value: value,
		Type:  parser.LITERAL,
	}
}
