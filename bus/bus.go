package bus

import (
	"container/list"
	"sync"

	"github.com/bruunoromero/cpu-emulator/parser"
)

type bus struct {
	frequency int
	buffer    *list.List
	wg        sync.WaitGroup
	channels  map[string]chan action
}

type action struct {
	Signal  int
	Origin  string
	Payload parser.Msg
}

type msg struct {
	channel string
	action  action
}

// READ and WRITE are the possible signals of an bus operation
const (
	READ = iota
	WRITE
)

// Instance is the interface of the bus type
type Instance interface {
	Wait()
	Start()
	MakeChannel(string)
	ReceiveFrom(string) chan action
	SendTo(string, string, int, parser.Msg)
}

// New returns a new instance of bus
func New(frequency int) Instance {
	return &bus{
		frequency: frequency,
		buffer:    list.New(),
		channels:  make(map[string]chan action),
	}
}

func (bus *bus) MakeChannel(channel string) {
	bus.channels[channel] = make(chan action)
}

func (bus *bus) Wait() {
	bus.wg.Wait()
}

func (bus *bus) SendTo(channel string, origin string, signal int, payload parser.Msg) {
	act := action{Signal: signal, Payload: payload, Origin: origin}

	bus.buffer.PushBack(msg{channel: channel, action: act})
}

func (bus *bus) ReceiveFrom(channel string) chan action {
	return bus.channels[channel]
}

func (bus *bus) send() {
	go func() {
		front := bus.buffer.Front()
		if front != nil {
			el := front.Value.(msg)

			bus.buffer.Remove(front)
			bus.channels[el.channel] <- el.action
		}
	}()
}

func (bus *bus) Start() {

}
