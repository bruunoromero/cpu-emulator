package io

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type io struct {
	read  chan string
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
			read: make(chan string),
		}
	})

	return instance
}

func (io *io) Run() {
	go func(ch chan string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')

			if err != nil {
				close(ch)
				return
			}

			ch <- s
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
