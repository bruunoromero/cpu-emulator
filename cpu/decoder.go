package cpu

import (
	"github.com/bruunoromero/cpu-emulator/utils"
)

type decoder struct {
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

func newDecoder(word int) decoder {
	return decoder{
		wordLength: word,
	}

}

func makeChunks(size int, arr []int8) [][]int8 {
	var divided [][]int8

	for i := 0; i < len(arr); i += size {
		end := i + size

		if end > len(arr) {
			end = len(arr)
		}

		divided = append(divided, arr[i:end])
	}

	return divided
}

func (decoder *decoder) decode(payload []int8) action {
	action := action{
		params: make([]value, 0),
	}

	numBytes := decoder.wordLength / 8

	action.isRegister = true
	action.location = -(utils.SumIntByteArray(payload[numBytes:numBytes*2]) + 1)

	action.action = utils.SumIntByteArray(payload[:numBytes])
	chunks := makeChunks(numBytes, payload[numBytes*2:])

	for _, chunk := range chunks {
		v := utils.SumIntByteArray(chunk)
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
