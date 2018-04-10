package utils

import (
	"strconv"
)

// IntToByteString converts an integer to its binary representation in string
func IntToByteString(v int) string {
	return strconv.FormatInt(int64(v), 2)
}

// ByteStringToInt converts a binary representation in string to its value in int8
func ByteStringToInt(v string) int8 {
	vl, err := strconv.ParseInt(v, 2, 64)

	if err != nil {
		Abort("Could not parse string to byte")
	}

	return int8(vl)
}

func SumIntByte(args ...int) {
	
}
