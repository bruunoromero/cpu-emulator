package memory

import (
	"github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type memory struct {
	wordLength        int
	lastWritePosition int
	list              []int8
}

// Instance is the interface for the memory type
type Instance interface {
	write([]int8)
	read() []int8
	Run(bus bus.Instance)
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
		list:              make([]int8, maxWords),
	}
}

func (memory *memory) Run(bus bus.Instance) {
	go func() {
		for {
			value := bus.ReceiveFrom("memory")
			memory.write(value)
		}
	}()
}

func (memory *memory) read() []int8 {
	return nil
}

func (memory *memory) write(payload []int8) {
	if memory.lastWritePosition+4 > len(memory.list) {
		memory.lastWritePosition = 0
	}

	for i, v := range payload {
		memory.list[memory.lastWritePosition+1] = v
		memory.lastWritePosition += i
	}
}
