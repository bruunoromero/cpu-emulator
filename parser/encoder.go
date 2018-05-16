package parser

import (
	"strconv"
	"strings"

	"github.com/bruunoromero/cpu-emulator/utils"
)

type Encoder struct {
	wordLength int
	registers  map[string]int
}

type Msg struct {
	Index  int
	Lenght int
	Value  byte
}

// This constants represents all possible actions in the cpu
const (
	Add = iota
	Mov
	Inc
	Imul
)

var actions = map[string]byte{
	"mov":  Mov,
	"add":  Add,
	"inc":  Inc,
	"imul": Imul,
}

func NewEncoder(registers []string, word int) Encoder {
	rgs := make(map[string]int)

	for i, register := range registers {
		rgs[register] = -(i + 1)
	}

	return Encoder{registers: rgs, wordLength: word}
}

func (encoder *Encoder) encode(action string, params []string) []byte {
	payload := encoder.expandValue(int(getAction(action)))
	return append(payload, encoder.mapParams(params)...)
}

func (encoder *Encoder) Parse(code string) [][]byte {
	var exprs [][]byte
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

func (encoder *Encoder) mapParams(params []string) []byte {
	var prs []byte

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

func (encoder *Encoder) expandValue(value int) []byte {
	return utils.ToBytes(encoder.wordLength, value)
}

func getAction(action string) byte {
	val, ok := actions[action]

	if !ok {
		utils.Abort("Unexpected action")
	}

	return val
}
