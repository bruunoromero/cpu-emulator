package memory

import (
	b "github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type memory struct {
	wordLength        int
	lastWritePosition int
	list              [][]byte
}

// Instance is the interface for the memory type
type Instance interface {
	write([]byte)
	read(byte) []byte
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
		list:              make([][]byte, maxWords),
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

func (memory *memory) read(payload byte) []byte {
	return memory.list[payload]
}

func (memory *memory) write(payload []byte) {
	if memory.lastWritePosition > len(memory.list) {
		memory.lastWritePosition = 0
	}

	memory.list[memory.lastWritePosition] = payload
	memory.lastWritePosition++
}
