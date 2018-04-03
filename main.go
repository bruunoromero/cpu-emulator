package main

import (
	"github.com/bruunoromero/cpu-emulator/io"
)

func main() {
	ioInstance := io.GetInstance(10)
	ioInstance.Run()
}
