package parser

import (
	"github.com/bruunoromero/cpu-emulator/utils"
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

func makeChunks(size int, arr []byte) [][]byte {
	var divided [][]byte

	for i := 0; i < len(arr); i += size {
		end := i + size

		if end > len(arr) {
			end = len(arr)
		}

		divided = append(divided, arr[i:end])
	}

	return divided
}

// Decode decodes an array of messages into an action
func (decoder *Decoder) Decode(payload []byte) Action {
	action := Action{
		Params: make([]Value, 0),
	}

	numBytes := decoder.wordLength / 8

	action.IsRegister = true
	action.Location = -(utils.FromBytes(decoder.wordLength, payload[numBytes:numBytes*2]) + 1)

	action.Action = utils.FromBytes(decoder.wordLength, payload[:numBytes])
	chunks := makeChunks(numBytes, payload[numBytes*2:])

	for _, chunk := range chunks {
		v := utils.FromBytes(decoder.wordLength, chunk)
		vl := Value{}
		vl.IsRegister = v < 0

		if vl.IsRegister {
			vl.Value = -(v + 1)
		} else {
			vl.Value = v
		}

		action.Params = append(action.Params, vl)
	}

	return action
}
