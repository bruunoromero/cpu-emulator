package io

import (
	"bufio"
	"os"

	"github.com/bruunoromero/cpu-emulator/bus"
)

type io struct {
	encoder encoder
	read    chan []int
	write   <-chan string
}

// Instance is the interface of the io type
type Instance interface {
	Run(bus.Instance)
}

// New returns a new instance of the I/O Module
func New(registers []string) Instance {
	return &io{
		read:    make(chan []int),
		encoder: newEncoder(registers),
	}
}

// Run will start the Cycle from read and write from I/O
func (io *io) Run(bus bus.Instance) {
	go func(ch chan []int) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')

			if err != nil {
				close(ch)
				return
			}

			exprs := io.encoder.parse(s)

			for _, expr := range exprs {
				ch <- expr
			}
		}
	}(io.read)

	for {
		stdin, ok := <-io.read
		if !ok {
			break
		} else {
			bus.SendToCPU(stdin)
		}
	}
}
