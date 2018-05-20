package parser

import (
	"github.com/bradfitz/slice"
	"github.com/bruunoromero/cpu-emulator/utils"
)

// Decoder is an parser type
type Decoder struct {
	wordLength int
}

// Action represents an expression
type Action struct {
	Action     int
	Signal     int
	Origin     int
	Location   Parameter
	Parameters []Parameter
}

type Parameter struct {
	Value int
	Type  int
}

// NewDecoder instanciate an returns a Decoder
func NewDecoder(word int) Decoder {
	return Decoder{
		wordLength: word,
	}
}

func makeChunks(size int, arr []Msg) [][]Msg {
	var divided [][]Msg

	for i := 0; i < len(arr); i += size {
		end := i + size

		if end > len(arr) {
			end = len(arr)
		}

		divided = append(divided, arr[i:end])
	}

	return divided
}

func (decoder *Decoder) isMsgComplete(msg []Msg) bool {
	if len(msg) == 0 {
		return false
	}

	firstMsg := msg[0]
	return firstMsg.Lenght == len(msg)-1
}

func (decoder *Decoder) groupMessages(msgs []Msg) [][]Msg {
	groups := make([][]Msg, 0)
	tmp := make([]Msg, len(msgs))
	copy(tmp, msgs)

	slice.Sort(tmp, func(left int, right int) bool {
		return tmp[left].Key < tmp[right].Key
	})

	group := make([]Msg, 0)
	for index, msg := range tmp {
		if index == 0 {
			group = append(group, msg)

			if index == len(tmp)-1 {
				groups = append(groups, group)
			}
		} else {
			if tmp[index-1].Key == msg.Key {
				group = append(group, msg)

				if index == len(tmp)-1 {
					groups = append(groups, group)
				}

			} else {
				groups = append(groups, group)

				group = make([]Msg, 0)
				group = append(group, msg)
			}
		}
	}

	return groups
}

func mapSlice(vs []Msg, f func(Msg) byte) []byte {
	vsm := make([]byte, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func getValue(msg Msg) byte {
	return msg.Value
}

// Decode decodes an array of messages into an action
func (decoder *Decoder) Decode(payload []Msg) Action {

	tmp := make([]Msg, len(payload))
	copy(tmp, payload)

	slice.Sort(tmp, func(left int, right int) bool {
		return tmp[left].Index < tmp[right].Index
	})

	numBytes := decoder.wordLength / 8

	actionValue := mapSlice(tmp[:numBytes], getValue)
	locationValue := mapSlice(tmp[numBytes:numBytes*2], getValue)

	action := Action{
		Action: utils.FromBytes(decoder.wordLength, actionValue),
	}

	msg := tmp[numBytes]
	value := utils.FromBytes(decoder.wordLength, locationValue)

	if msg.Type == REGISTER {
		value = -(value + 1)
	}

	action.Location = Parameter{
		Type:  msg.Type,
		Value: value,
	}

	chunks := makeChunks(numBytes, tmp[numBytes*2:])

	for _, chunk := range chunks {
		chunkValue := mapSlice(chunk, getValue)

		v := utils.FromBytes(decoder.wordLength, chunkValue)
		vl := Parameter{}

		vl.Type = chunk[0].Type

		if vl.Type == REGISTER {
			vl.Value = (v + 1) * -1
		} else {
			vl.Value = v
		}

		action.Parameters = append(action.Parameters, vl)
	}

	return action
}

func (decoder *Decoder) GetMessagesWithQueue(address []Msg, data []Msg, instructions []Msg, queue *[]Msg) [][]Msg {
	messages := make([][]Msg, 0)

	msg := make([]Msg, 0)

	if data != nil {
		msg = append(msg, data...)
	}

	if address != nil {
		msg = append(msg, address...)
	}

	if instructions != nil {
		msg = append(msg, instructions...)
	}

	*queue = append(*queue, msg...)

	groups := decoder.groupMessages(*queue)

	for _, msgs := range groups {
		if !decoder.isMsgComplete(msgs) {
			break
		}

		messages = append(messages, msgs)
		*queue = decoder.removeMessagesFromQueue(msgs, *queue)
	}

	return messages
}

func (decoder *Decoder) removeMessagesFromQueue(msgs []Msg, queue []Msg) []Msg {
	tmp := make([]Msg, len(queue))

	copy(tmp, queue)

	for _, msg := range msgs {
		for index, queueMsg := range tmp {
			if queueMsg == msg {
				tmp = append(tmp[:index], tmp[index+1:]...)
			}
		}
	}

	return tmp
}
