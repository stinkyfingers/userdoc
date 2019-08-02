package foobar

import (
	"errors"
	"fmt"
	"os"

	"github.com/stinkyfingers/userdoc/test_files/funk"
)

// ErrNum is a test err
var ErrNum = errors.New("this is a num err")

// Foo is a test function
func Foo(name string) (int, error) {

	fmt.Printf("my name is %s\n", name)                 // displays user name
	return os.Stdout.Write([]byte("writing to stdout")) // writes a stdout message
}

func bar(number int) error {
	fmt.Println("HERE")
	if number > 0 {
		return ErrNum // indicates that the number is greater than zero
	}
	Foo("fooing here") // this foos
	return funk.Funk() // this is a FUNK comment
	// Foo("hey")  // we are fooing and heying
	// return nil
}
