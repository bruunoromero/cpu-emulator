package memory

import (
	b "github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type memory struct {
	wordLength        int
	lastWritePosition int
	list              [][]int8
}

// Instance is the interface for the memory type
type Instance interface {
	write([]int8)
	read(int8) []int8
	Run(bus b.Instance)
}

// New returns a new instance of Memory
func New(size int, wordLength int) Instance {
	wordLengthByte := wordLength / 8

	maxWords := size / wordLengthByte

	if maxWords/4 < 1 {
		utils.Abort("Cannot instanciate ram with this length")
	}

	return &memory{
		lastWritePosition: 0,
		wordLength:        wordLength,
		list:              make([][]int8, maxWords),
	}
}

func (memory *memory) Run(bus b.Instance) {
	go func() {
		for {
			select {
			case value := <-bus.ReceiveFrom("memory"):
				if value.Signal == b.READ {
					bus.SendTo(value.Origin, "memory", b.WRITE, memory.read(value.Payload[0]))
				} else {
					memory.write(value.Payload)
				}
			}
		}
	}()
}

func (memory *memory) read(payload int8) []int8 {
	return memory.list[payload]
}

func (memory *memory) write(payload []int8) {
	if memory.lastWritePosition > len(memory.list) {
		memory.lastWritePosition = 0
	}

	memory.list[memory.lastWritePosition] = payload
	memory.lastWritePosition++

	// fmt.Println(memory.list)
}
