package bus

import (
	"sync"
	"time"
)

type bus struct {
	wg       sync.WaitGroup
	channels map[string]chan action
}

type action struct {
	Signal  int
	Origin  string
	Payload []byte
}

// READ and WRITE are the possible signals of an bus operation
const (
	READ = iota
	WRITE
)

// Instance is the interface of the bus type
type Instance interface {
	Wait()
	MakeChannel(string)
	ReceiveFrom(string) chan action
	SendTo(string, string, int, []byte)
}

// New returns a new instance of bus
func New() Instance {
	return &bus{
		channels: make(map[string]chan action),
	}
}

func (bus *bus) MakeChannel(channel string) {
	bus.channels[channel] = make(chan action)
}

func (bus *bus) Wait() {
	bus.wg.Wait()
}

func (bus *bus) SendTo(channel string, origin string, signal int, payload []byte) {
	bus.wg.Add(1)
	go func() {
		time.Sleep(500 * time.Microsecond)
		bus.channels[channel] <- action{Signal: signal, Payload: payload, Origin: origin}
		bus.wg.Done()
	}()
}

func (bus *bus) ReceiveFrom(channel string) chan action {
	return bus.channels[channel]
}
