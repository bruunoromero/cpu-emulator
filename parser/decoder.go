package parser

import (
	"github.com/bruunoromero/cpu-emulator/utils"
)

type Decoder struct {
	wordLength int
}

type action struct {
	action     int
	location   int
	isRegister bool
	params     []value
}

type value struct {
	value      int
	isRegister bool
}

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

func (decoder *Decoder) Decode(payload []byte) action {
	action := action{
		params: make([]value, 0),
	}

	numBytes := decoder.wordLength / 8

	action.isRegister = true
	action.location = -(utils.FromBytes(decoder.wordLength, payload[numBytes:numBytes*2]) + 1)

	action.action = utils.FromBytes(decoder.wordLength, payload[:numBytes])
	chunks := makeChunks(numBytes, payload[numBytes*2:])

	for _, chunk := range chunks {
		v := utils.FromBytes(decoder.wordLength, chunk)
		vl := value{}
		vl.isRegister = v < 0

		if vl.isRegister {
			vl.value = -(v + 1)
		} else {
			vl.value = v
		}

		action.params = append(action.params, vl)
	}

	return action
}
