package io

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/bruunoromero/cpu-emulator/parser"
	"github.com/bruunoromero/cpu-emulator/utils"

	b "github.com/bruunoromero/cpu-emulator/bus"
)

type io struct {
	read    chan []parser.Msg
	write   <-chan string
	encoder parser.Encoder
}

// Instance is the interface of the io type
type Instance interface {
	Run(b.Instance)
}

// New returns a new instance of the I/O Module
func New(encoder parser.Encoder) Instance {
	return &io{
		encoder: encoder,
		read:    make(chan []parser.Msg),
	}
}

// Run will start the Cycle from read and write from I/O
func (io *io) Run(bus b.Instance) {
	go func() {
		path, fileErr := filepath.Abs("./code.s")

		if fileErr != nil {
			utils.Abort("Could not get path of the file")
		}

		inFile, openErr := os.Open(path)

		if openErr != nil {
			utils.Abort("Could not open the file")
		}

		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)

		codeIndex := 0
		for scanner.Scan() {
			instructions := io.encoder.ExpandInstruction(scanner.Text())
			exprs := make([][]parser.Msg, 0)
			for _, instruction := range instructions {
				exprs = append(exprs, io.encoder.Parse(codeIndex, instruction)...)
				codeIndex++
			}

			for _, expr := range exprs {
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
			}
		}
	}
}
