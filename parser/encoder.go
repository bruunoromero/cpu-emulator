package parser

import (
	"strconv"
	"strings"

	"github.com/bruunoromero/cpu-emulator/utils"
)

// This constants represents all possible actions in the cpu
const (
	GT = iota
	LT
	EQ
	Add
	Mov
	Inc
	Imul
	Jump
	NULL
	GTEQ
	LTEQ
	Label
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
	Signal int
	Value  byte
	Origin string
}

var actions = map[string]byte{
	"EQ":    EQ,
	"GT":    GT,
	"LT":    LT,
	"mov":   Mov,
	"add":   Add,
	"inc":   Inc,
	"JMP":   Jump,
	"imul":  Imul,
	"NULL":  NULL,
	"GTEQ":  GTEQ,
	"LTEQ":  LTEQ,
	"label": Label,
}

var conditionals = map[string]string{
	"=":  "EQ",
	">":  "GT",
	"<":  "LT",
	">=": "GTEQ",
	"<=": "LTEQ",
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
	payload = append(payload, encoder.MapParams(params)...)

	for index := range payload {
		payload[index].Key = key
		payload[index].Index = index
		payload[index].Lenght = len(payload) - 1
	}

	return payload
}

func isConditional(line string) bool {
	for key := range conditionals {
		if strings.Contains(line, key) {
			return true
		}
	}

	return false
}

func (encoder *Encoder) ExpandInstruction(code string) []string {
	if isConditional(code) {
		parts := strings.Split(code, ":")

		if len(parts) != 3 {
			utils.Abort("Unexpected action")
		}

		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}

		for key, value := range conditionals {
			if strings.Contains(parts[0], key) {
				insts := strings.Split(parts[0], key)
				for i, inst := range insts {
					insts[i] = strings.TrimSpace(inst)
				}

				parts[0] = value + " " + insts[0] + ", " + insts[1]
			}
		}

		return parts
	}

	return []string{code}
}

// Parse parses an string into an matrix of bytes
func (encoder *Encoder) Parse(codeIndex int, code string) [][]Msg {
	var exprs [][]Msg
	lines := strings.Split(code, ";")

	for _, line := range lines {
		commaReplaced := strings.Replace(line, ",", " ", -1)
		values := strings.Fields(commaReplaced)

		params := make([]string, 0)
		action := values[0]

		if len(values) > 1 {
			params = values[1:len(values)]
		}

		expr := encoder.encode(codeIndex, action, params)
		exprs = append(exprs, expr)

	}

	return exprs
}

func (encoder *Encoder) MapParams(params []string) []Msg {
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
