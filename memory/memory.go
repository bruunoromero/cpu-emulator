package memory

import (
	"time"

	b "github.com/bruunoromero/cpu-emulator/bus"
	"github.com/bruunoromero/cpu-emulator/parser"
	"github.com/bruunoromero/cpu-emulator/utils"
)

type memory struct {
	wordLength        int
	lastWritePosition int
	list              [][]parser.Msg
	messageQueue      []parser.Msg
	decoder           parser.Decoder
}

// Instance is the interface for the memory type
type Instance interface {
	Run(bus b.Instance)
	read(byte) []parser.Msg
	write(int, []parser.Msg) int
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
		list:              make([][]parser.Msg, length),
		decoder:           parser.NewDecoder(wordLength),
	}
}

func (memory *memory) Run(bus b.Instance) {
	go func() {
		for {
			data := bus.ReceiveFrom("memoryData")
			address := bus.ReceiveFrom("memoryAddress")
			instructions := bus.ReceiveFrom("memoryInstruction")

			messages := memory.decoder.GetMessagesWithQueue(address.Payload, data.Payload, instructions.Payload, &memory.messageQueue)

			for _, message := range messages {
				if len(message) > 0 {
					msg := message[0]

					if msg.Signal == b.WRITE {
						if msg.Origin == "io" {
							position := memory.write(0, message)
							bus.SendTo("cpu", "memory", b.READ, []parser.Msg{parser.Msg{Key: position, Index: 0, Lenght: 0, Type: parser.REGISTER, Value: byte(position)}})
						} else if msg.Origin == "cpu" {
							memory.write(len(memory.list)/2, message)
						}
					} else if msg.Signal == b.READ {
						v := memory.read(msg.Value)
						bus.SendTo(msg.Origin, "memory", b.WRITE, v)
					}
				}
			}

			time.Sleep(time.Second / 2)
		}
	}()
}

func (memory *memory) read(payload byte) []parser.Msg {
	return memory.list[payload]
}

func (memory *memory) write(offset int, payload []parser.Msg) int {
	if memory.lastWritePosition == len(memory.list) {
		memory.lastWritePosition = 0
	}

	position := memory.lastWritePosition + offset
	memory.list[position] = payload
	memory.lastWritePosition++

	return position
}
