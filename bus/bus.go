package bus

import (
	"container/list"
	"time"

	"github.com/bruunoromero/cpu-emulator/parser"
)

type bus struct {
	length    int
	frequency int
	buffer    *list.List
	channels  map[string]chan action
}

type action struct {
	Signal  int
	Origin  string
	Payload []parser.Msg
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

// DATA constant
const DATA = "Data"

// ADDRESS constant
const ADDRESS = "Address"

// INSTUCTION constant
const INSTUCTION = "Instruction"

var lanes = []string{ADDRESS, DATA, INSTUCTION}

// Instance is the interface of the bus type
type Instance interface {
	Start()
	MakeChannel(string)
	ReceiveFrom(string) chan action
	SendTo(string, string, int, []parser.Msg)
}

// New returns a new instance of bus
func New(frequency int, length int) Instance {
	return &bus{
		length:    length,
		frequency: frequency,
		buffer:    list.New(),
		channels:  make(map[string]chan action),
	}
}

func (bus *bus) MakeChannel(channel string) {
	for _, lane := range lanes {
		bus.channels[channel+lane] = make(chan action)
	}
}

func getLane(msg parser.Msg) string {
	if msg.Type == parser.CALL {
		return INSTUCTION
	} else if msg.Type == parser.LITERAL {
		return DATA
	} else if msg.Type == parser.MEMORY || msg.Type == parser.REGISTER {
		return ADDRESS
	}

	return ""
}

func (bus *bus) SendTo(channel string, origin string, signal int, payload []parser.Msg) {
	lanes := categorizeMsgs(payload)
	expandedLanes := bus.expandCategories(lanes)

	for _, category := range expandedLanes {
		for lane, msgs := range category {
			act := action{Signal: signal, Payload: msgs, Origin: origin}
			bus.buffer.PushBack(msg{channel: channel + lane, action: act})
		}
	}
}

func (bus *bus) ReceiveFrom(channel string) chan action {
	return bus.channels[channel]
}

func (bus *bus) expandCategories(categories map[string][]parser.Msg) []map[string][]parser.Msg {
	expanded := make([]map[string][]parser.Msg, 0)

	for lane, msgs := range categories {
		chunks := bus.chuckifyMsgs(msgs)

		for _, chunk := range chunks {
			var category map[string][]parser.Msg

			category[lane] = chunk
			expanded = append(expanded, category)
		}
	}

	return expanded
}

func (bus *bus) chuckifyMsgs(msgs []parser.Msg) [][]parser.Msg {
	var divided [][]parser.Msg

	chunkSize := bus.length / 8

	for i := 0; i < len(msgs); i += chunkSize {
		end := i + chunkSize

		if end > len(msgs) {
			end = len(msgs)
		}

		divided = append(divided, msgs[i:end])
	}

	return divided
}

func categorizeMsgs(msgs []parser.Msg) map[string][]parser.Msg {
	var lanes map[string][]parser.Msg

	for _, msg := range msgs {
		lane := getLane(msg)

		if lanes[lane] == nil {
			lanes[lane] = make([]parser.Msg, 0)
		}

		lanes[lane] = append(lanes[lane], msg)
	}

	return lanes
}

func (bus *bus) send(front *list.Element, channelsLength map[string]int) bool {
	if front != nil {
		el := front.Value.(msg)
		if channelsLength[el.channel]+8 < bus.length {
			channelsLength[el.channel] += 8
			go func() {
				bus.channels[el.channel] <- el.action
			}()
			return true
		}
	}

	return false
}

func (bus *bus) Start() {
	go func() {
		for {
			sent := make([]*list.Element, 0)
			channelsLength := make(map[string]int)
			for front := bus.buffer.Front(); front != nil; front = front.Next() {
				if bus.send(front, channelsLength) {
					sent = append(sent, front)
				}
			}

			for _, msg := range sent {
				bus.buffer.Remove(msg)
			}

			time.Sleep(time.Second)
		}
	}()
}
