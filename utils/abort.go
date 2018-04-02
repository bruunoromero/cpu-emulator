package utils

import (
	"fmt"
	"os"
)

// Abort will log a error and stop the applicaiton
func Abort(err string) {
	fmt.Println(err)
	os.Exit(1)
}
