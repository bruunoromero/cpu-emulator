package bus

import "github.com/bruunoromero/cpu-emulator/io"

type bus struct {
	io <-chan io.Expr
}

// New returns a new instance of bus
func New() {

}
