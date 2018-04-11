package utils

import (
	"encoding/binary"
	"math"
)

// ToBytes convert a value to an array of bytes
func ToBytes(size int, value int) []byte {
	b := make([]byte, size/8)
	v := int(math.Pow(2, float64(size)) / 2)

	if size == 16 {
		binary.LittleEndian.PutUint16(b, uint16(value+v))
	} else if size == 32 {
		v := int(math.Pow(2, 16) / 2)
		binary.LittleEndian.PutUint32(b, uint32(value+v))
	} else if size == 64 {
		binary.LittleEndian.PutUint64(b, uint64(value+v))
	}

	return b
}

// FromBytes convert an array of bytes to a value
func FromBytes(size int, bytes []byte) int {
	v := int(math.Pow(2, float64(size)) / 2)

	if size == 16 {
		return int(binary.LittleEndian.Uint16(bytes)) - v
	} else if size == 32 {
		return int(binary.LittleEndian.Uint32(bytes)) - v
	} else if size == 64 {
		return int(binary.LittleEndian.Uint64(bytes)) - v
	}

	Abort("Unexpected size")
	return 0
}
