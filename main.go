package main

import (
	"fmt"

	"github.com/bruunoromero/cpu-emulator/vm"
)

func getWordLengh() int {
	var vl int

	fmt.Println("Qual o tamano da palavra em bits (16, 32, 64): ")

	for {
		fmt.Scanf("%d", &vl)
		if vl == 16 || vl == 32 || vl == 64 {
			return int(vl)
		}
	}
}

func getMemoryLengh() int {
	var vl int

	fmt.Println("Qual o tamano da ram em bytes (8, 16, 32): ")

	for {
		fmt.Scanf("%d", &vl)
		if vl == 8 || vl == 16 || vl == 32 {
			return int(vl)
		}
	}
}

func main() {

	word := getWordLengh()
	memory := getMemoryLengh()

	fmt.Println("")
	fmt.Println("-----------------------")
	fmt.Println("      Starting VM      ")
	fmt.Println("-----------------------")
	fmt.Println("")
	fmt.Println("Log: VM Started")
	fmt.Println("")

	vm.Start([]string{"A", "B", "C", "D", "E"}, word, memory)

}
