package io

import (
	"container/list"
	"strings"
)

// Parse will transform a string into a list of expressions
func parse(code string) *list.List {
	exprs := list.New()
	lines := strings.Split(code, ";")

	for _, line := range lines {
		commaReplaced := strings.Replace(line, ",", " ", -1)
		values := strings.Fields(commaReplaced)

		if len(values) > 1 {
			action := values[0]
			params := values[1:len(values)]

			expr := encode(action, params)
			exprs.PushBack(expr)
		}

	}

	return exprs
}
