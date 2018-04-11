package io

import (
	"bufio"
	"os"

	b "github.com/bruunoromero/cpu-emulator/bus"
)

type io struct {
	encoder encoder
	read    chan []byte
	write   <-chan string
}

// Instance is the interface of the io type
type Instance interface {
	Run(b.Instance)
}

// New returns a new instance of the I/O Module
func New(registers []string, word int) Instance {
	return &io{
		read:    make(chan []byte),
		encoder: newEncoder(registers, word),
	}
}

// Run will start the Cycle from read and write from I/O
func (io *io) Run(bus b.Instance) {
	values := 0
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')

			if err != nil {
				close(io.read)
				return
			}

			exprs := io.encoder.parse(s)

			for _, expr := range exprs {
				values++
				io.read <- expr
			}

		}
	}()

	for {
		select {
		case stdin, ok := <-io.read:
			if !ok {
				break
			} else {
				bus.SendTo("memory", "io", b.WRITE, stdin)
				bus.SendTo("cpu", "io", b.WRITE, []byte{})
			}
		}
	}
}
