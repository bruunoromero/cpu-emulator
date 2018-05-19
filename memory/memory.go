package memory

import (
	"fmt"
	"time"

	b "github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/parser"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type memory struct {
	wordLength        int
	lastWritePosition int
	list              [][]byte
	messageQueue      []parser.Msg
	decoder           parser.Decoder
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
	length := maxWords / 4

	if length < 1 {
		utils.Abort("Cannot instanciate ram with this length")
	}

	return &memory{
		lastWritePosition: 0,
		wordLength:        wordLength,
		messageQueue:      make([]parser.Msg, 0),
		list:              make([][]byte, length),
		decoder:           parser.NewDecoder(wordLength),
	}
}

func (memory *memory) Run(bus b.Instance) {
	go func() {
		for {
			memory.getMessage(bus)
			// fmt.Println(memory.getMessage(bus))
			// if value.Signal == b.READ {
			// 	// bus.SendTo(value.Origin, "memory", b.WRITE, memory.read(value.Payload[0]))
			// } else {
			// 	// memory.write(value.Payload)
			// }
			time.Sleep(time.Second / 2)
		}
	}()
}

func (memory *memory) getMessage(bus b.Instance) *parser.Message {
	data := bus.ReceiveFrom("memoryData")
	address := bus.ReceiveFrom("memoryAddress")
	instructions := bus.ReceiveFrom("memoryInstruction")

	msg := make([]parser.Msg, 0)

	if data != nil {
		msg = append(msg, data.Payload...)
	}

	if address != nil {
		msg = append(msg, address.Payload...)
	}

	if instructions != nil {
		msg = append(msg, instructions.Payload...)
	}

	memory.messageQueue = append(memory.messageQueue, msg...)

	groups := memory.decoder.GroupMessages(memory.messageQueue)

	for _, msgs := range groups {
		if memory.decoder.IsMsgComplete(msgs) {
			fmt.Println("ola")
		}
	}

	return nil
}

func (memory *memory) read(payload byte) []byte {
	return memory.list[payload]
}

func (memory *memory) write(payload []byte) {
	if memory.lastWritePosition == len(memory.list) {
		memory.lastWritePosition = 0
	}

	memory.list[memory.lastWritePosition] = payload
	memory.lastWritePosition++
}
