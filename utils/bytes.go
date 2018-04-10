package utils

import (
	"strconv"
	"strings"
)

// IntToByteString converts an integer to its binary representation in string
func IntToByteString(v int) string {
	return strconv.FormatInt(int64(v), 2)
}

// ByteStringToInt converts a binary representation in string to its value in int8
func ByteStringToInt(v string) int {
	vl, err := strconv.ParseInt(v, 2, 64)

	if err != nil {
		Abort("Could not parse string to byte")
	}

	return int(vl)
}

// SumIntByteArray sums a array of int8 into a int
func SumIntByteArray(args []int8) int {
	isNegative := false
	str := ""

	for _, v := range args {
		str += IntToByteString(int(v))

		if strings.HasPrefix(str, "-") {
			isNegative = true
		}
	}

	if isNegative {
		str = "-" + str
	}

	res := ByteStringToInt(str)

	return res
}
