package main

import (
	"github.com/bruunoromero/cpu-emulator/io"
)

func main() {
	ioInstance := io.GetInstance()
	ioInstance.Run()
}
