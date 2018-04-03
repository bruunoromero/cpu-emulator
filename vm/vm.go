package vm

import "sync"

var once sync.Once

// Start initiates the Von Neumann loop
func Start(busLength int, ramLength int, wordLength int) {
	once.Do(func() {
		
	})
}
