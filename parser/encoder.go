package parser

import (
	"strconv"
	"strings"

	"github.com/bruunoromero/cpu-emulator/utils"
)

// This constants represents all possible actions in the cpu
const (
	Add = iota
	Mov
	Inc
	Imul
)

// This constants represents all possible types of messages
const (
	CALL = iota
	MEMORY
	LITERAL
	REGISTER
)

// Encoder is an parser type
type Encoder struct {
	wordLength int
	registers  map[string]int
}

// Msg is the payload type of the bus
type Msg struct {
	Key    int
	Type   int
	Index  int
	Lenght int
	Value  byte
}

var actions = map[string]byte{
	"mov":  Mov,
	"add":  Add,
	"inc":  Inc,
	"imul": Imul,
}

// NewEncoder instanciate an returns an encoder
func NewEncoder(registers []string, word int) Encoder {
	rgs := make(map[string]int)

	for i, register := range registers {
		rgs[register] = -(i + 1)
	}

	return Encoder{registers: rgs, wordLength: word}
}

func (encoder *Encoder) encode(key int, action string, params []string) []Msg {
	payload := encoder.expandValue(int(getAction(action)), CALL)
	payload = append(payload, encoder.mapParams(params)...)

	for index := range payload {
		payload[index].Key = key
		payload[index].Index = index
		payload[index].Lenght = len(payload) - 1
	}

	return payload
}

// Parse parses an string into an matrix of bytes
func (encoder *Encoder) Parse(codeIndex int, code string) [][]Msg {
	var exprs [][]Msg
	lines := strings.Split(code, ";")

	for _, line := range lines {
		commaReplaced := strings.Replace(line, ",", " ", -1)
		values := strings.Fields(commaReplaced)

		if len(values) > 1 {
			action := values[0]
			params := values[1:len(values)]

			expr := encoder.encode(codeIndex, action, params)
			exprs = append(exprs, expr)
		}

	}

	return exprs
}

func (encoder *Encoder) mapParams(params []string) []Msg {
	var prs []Msg

	for _, param := range params {
		if strings.HasPrefix(param, "0x") {
			// If the case matches, the parameter is a memory
			value, err := strconv.ParseInt(strings.TrimPrefix(param, "0x"), 16, 64)

			if err != nil {
				utils.Abort("Cannot convert value")
			}

			prs = append(prs, encoder.expandValue(int(value), MEMORY)...)
		} else {
			value, err := strconv.Atoi(param)

			// If theres a error, than the value is a register
			if err != nil {
				register := encoder.registers[param]

				// If register is 0, then there is no register defined with that name
				if register != 0 {
					prs = append(prs, encoder.expandValue(int(register), REGISTER)...)
				} else {
					utils.Abort("Invalid register")
				}
			} else {
				prs = append(prs, encoder.expandValue(value, LITERAL)...)
			}

		}
	}

	return prs
}

func (encoder *Encoder) expandValue(value int, msgType int) []Msg {
	var msgs []Msg
	bytes := utils.ToBytes(encoder.wordLength, value)

	for _, value := range bytes {
		msg := Msg{
			Value: value,
			Type:  msgType,
		}

		msgs = append(msgs, msg)
	}

	return msgs
}

func getAction(action string) byte {
	val, ok := actions[action]

	if !ok {
		utils.Abort("Unexpected action")
	}

	return val
}
