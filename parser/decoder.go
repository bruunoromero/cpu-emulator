package parser

import (
	"github.com/bradfitz/slice"
)

// Decoder is an parser type
type Decoder struct {
	wordLength int
}

// Action represents an expression
type Action struct {
	Action     int
	Location   int
	IsRegister bool
	Params     []Value
}

type Message struct {
	Signal      int
	Origin      string
	Data        []Msg
	Address     []Msg
	Inctruction []Msg
}

// Value represents a parameter
type Value struct {
	Value      int
	IsRegister bool
}

// NewDecoder instanciate an returns a Decoder
func NewDecoder(word int) Decoder {
	return Decoder{
		wordLength: word,
	}
}

func makeChunks(size int, arr []Message) [][]Message {
	var divided [][]Message

	for i := 0; i < len(arr); i += size {
		end := i + size

		if end > len(arr) {
			end = len(arr)
		}

		divided = append(divided, arr[i:end])
	}

	return divided
}

func (decoder *Decoder) IsMsgComplete(msg []Msg) bool {
	slice.Sort(msg, func(left int, right int) bool {
		return msg[left].Index < msg[right].Index
	})

	lastMsg := msg[len(msg)-1]
	return lastMsg.Lenght == lastMsg.Index
}

func (decoder *Decoder) GroupMessages(msgs []Msg) [][]Msg {
	groups := make([][]Msg, 0)

	slice.Sort(msgs, func(left int, right int) bool {
		return msgs[left].Key < msgs[right].Key
	})

	group := make([]Msg, 0)
	for index, msg := range msgs {
		if index == 0 {
			group = append(group, msg)
		} else {
			if msgs[index-1].Key == msg.Key {
				group = append(group, msg)
			} else {
				groups = append(groups, group)

				group := make([]Msg, 0)
				group = append(group, msg)
			}
		}
	}

	return groups
}

// Decode decodes an array of messages into an action
func (decoder *Decoder) Decode(payload []Message) Action {
	action := Action{
		Params: make([]Value, 0),
	}

	// numBytes := decoder.wordLength / 8

	// action.IsRegister = true
	// action.Location = -(utils.FromBytes(decoder.wordLength, payload[numBytes:numBytes*2]) + 1)

	// action.Action = utils.FromBytes(decoder.wordLength, payload[:numBytes])
	// chunks := makeChunks(numBytes, payload[numBytes*2:])

	// for _, chunk := range chunks {
	// 	v := utils.FromBytes(decoder.wordLength, chunk)
	// 	vl := Value{}
	// 	vl.IsRegister = v < 0

	// 	if vl.IsRegister {
	// 		vl.Value = -(v + 1)
	// 	} else {
	// 		vl.Value = v
	// 	}

	// 	action.Params = append(action.Params, vl)
	// }

	return action
}
