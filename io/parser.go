package io

import (
	"strings"
)

// Parse will transform a string into a list of expressions
func parse(code string) []expr {
	var exprs []expr
	lines := strings.Split(code, ";")

	for _, line := range lines {
		commaReplaced := strings.Replace(line, ",", " ", -1)
		values := strings.Fields(commaReplaced)

		if len(values) > 1 {
			action := values[0]
			params := values[1:len(values)]

			expr := encode(action, params)
			exprs = append(exprs, expr)
		}

	}

	return exprs
}
