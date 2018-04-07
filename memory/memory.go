package memory

import (
	"fmt"

	"github.com/bruunoromero/cpu-emulator/bus"
)

type memory struct{}

// Instance is the interface for the memory type
type Instance interface {
	Read() []int
	Write([]int)
	Run(bus bus.Instance)
}

// New returns a new instance of Memory
func New() Instance {
	return &memory{}
}

func (memory *memory) Run(bus bus.Instance) {
	go func() {
		for {
			value := bus.ReceiveAtRAM()
			memory.Write(value)
		}
	}()
}

func (memory *memory) Read() []int {
	return nil
}

func (memory *memory) Write(payload []int) {
	fmt.Println(payload)
}
