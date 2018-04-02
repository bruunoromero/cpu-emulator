package io

import (
	"github.com/bruunoromero/cpu-emulator/utils"
)

// Expr is the type of a expression for the cpu
type Expr struct {
	Action int
	Params []string
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

func encode(action string, params []string) Expr {
	return Expr{Action: getAction(action), Params: params}
}

func getAction(action string) int {
	val, ok := actions[action]

	if !ok {
		utils.Abort("Unexpected action")
	}

	return val
}
