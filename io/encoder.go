package io

import (
	"strconv"
	"strings"

	"github.com/bruunoromero/cpu-emulator/utils"
)

type encoder struct {
	registers map[string]int
}

// This constants represents all possible actions in the cpu
const (
	Add = iota
	Mov
	Inc
	Imul
)

var actions = map[string]int{
	"mov":  Mov,
	"add":  Add,
	"inc":  Inc,
	"imul": Imul,
}

func newEncoder(registers []string) encoder {
	rgs := make(map[string]int)

	for i, register := range registers {
		rgs[register] = -(i + 1)
	}

	return encoder{registers: rgs}
}

func (encoder *encoder) encode(action string, params []string) []int {
	payload := []int{getAction(action)}
	return append(payload, encoder.mapParams(params)...)
}

func (encoder *encoder) parse(code string) [][]int {
	var exprs [][]int
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

func (encoder *encoder) mapParams(params []string) []int {
	var prs []int

	for _, param := range params {
		if strings.HasPrefix(param, "0x") {
			// If the case matches, the parameter is a memory
			value, err := strconv.ParseInt(strings.TrimPrefix(param, "0x"), 16, 64)

			if err != nil {
				utils.Abort("Cannot convert value")
			}

			prs = append(prs, int(value))
		} else {
			value, err := strconv.Atoi(param)

			// If theres a error, than the value is a register
			if err != nil {
				register := encoder.registers[param]

				// If register is 0, then there is no register defined with that name
				if register != 0 {
					prs = append(prs, register)
				} else {
					utils.Abort("Invalid register")
				}
			} else {
				prs = append(prs, int(value))
			}

		}
	}

	return prs
}

func getAction(action string) int {
	val, ok := actions[action]

	if !ok {
		utils.Abort("Unexpected action")
	}

	return val
}
