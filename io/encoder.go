package io

import (
	"strconv"
	"strings"

	"github.com/bruunoromero/cpu-emulator/utils"
)

type encoder struct {
	wordLength int
	registers  map[string]int8
}

// This constants represents all possible actions in the cpu
const (
	Add = iota
	Mov
	Inc
	Imul
)

var actions = map[string]int8{
	"mov":  Mov,
	"add":  Add,
	"inc":  Inc,
	"imul": Imul,
}

func newEncoder(registers []string, word int) encoder {
	rgs := make(map[string]int8)

	for i, register := range registers {
		rgs[register] = int8(-(i + 1))
	}

	return encoder{registers: rgs, wordLength: word}
}

func (encoder *encoder) encode(action string, params []string) []int8 {
	payload := encoder.expandValue(int(getAction(action)))
	return append(payload, encoder.mapParams(params)...)
}

func (encoder *encoder) parse(code string) [][]int8 {
	var exprs [][]int8
	lines := strings.Split(code, ";")

	for _, line := range lines {
		commaReplaced := strings.Replace(line, ",", " ", -1)
		values := strings.Fields(commaReplaced)

		if len(values) > 1 {
			action := values[0]
			params := values[1:len(values)]

			expr := encoder.encode(action, params)
			exprs = append(exprs, expr)
		}

	}

	return exprs
}

func (encoder *encoder) mapParams(params []string) []int8 {
	var prs []int8

	for _, param := range params {
		if strings.HasPrefix(param, "0x") {
			// If the case matches, the parameter is a memory
			value, err := strconv.ParseInt(strings.TrimPrefix(param, "0x"), 16, 64)

			if err != nil {
				utils.Abort("Cannot convert value")
			}

			prs = append(prs, encoder.expandValue(int(value))...)
		} else {
			value, err := strconv.Atoi(param)

			// If theres a error, than the value is a register
			if err != nil {
				register := encoder.registers[param]

				// If register is 0, then there is no register defined with that name
				if register != 0 {
					prs = append(prs, encoder.expandValue(int(register))...)
				} else {
					utils.Abort("Invalid register")
				}
			} else {
				prs = append(prs, encoder.expandValue(value)...)
			}

		}
	}

	return prs
}

func (encoder *encoder) expandValue(value int) []int8 {
	var res []int8

	v := value
	valueInserted := false

	for i := 0; i < encoder.wordLength/8; i++ {
		if v > 127 {
			res = append(res, 127)
			v -= 127
		} else if !valueInserted {
			valueInserted = true
			res = append([]int8{int8(v)}, res...)
		} else {
			res = append([]int8{0}, res...)
		}
	}

	return res
}

func getAction(action string) int8 {
	val, ok := actions[action]

	if !ok {
		utils.Abort("Unexpected action")
	}

	return val
}
