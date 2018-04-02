package io

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type io struct {
	read  chan expr
	write <-chan string
}

// Instance is the interface of the io singleton
type Instance interface {
	Run()
	private()
}

var instance *io
var once sync.Once

// GetInstance returns a new instance of the I/O Module
func GetInstance() Instance {
	once.Do(func() {
		instance = &io{
			read: make(chan expr),
		}
	})

	return instance
}

// Run will start the Cycle from read and write from I/O
func (io *io) Run() {
	go func(ch chan expr) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')

			if err != nil {
				close(ch)
				return
			}

			exprs := parse(s)

			for _, expr := range exprs {
				ch <- expr
			}
		}
		close(ch)
	}(io.read)

stdinloop:
	for {
		select {
		case stdin, ok := <-io.read:
			if !ok {
				break stdinloop
			} else {
				fmt.Println("Read input from stdin:", stdin)
			}
		case <-time.After(1 * time.Second):
			// Do something when there is nothing read from stdin
		}
	}
	fmt.Println("Done, stdin must be closed")
}

func (io *io) private() {}
