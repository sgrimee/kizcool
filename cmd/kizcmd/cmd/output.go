package cmd

import (
	"log"
	"os"

	"github.com/sgrimee/kizcool"
)

// output prints the object in the given format to stdout
func output(format string, obj interface{}) {
	if err := kizcool.Output(os.Stdout, format, obj); err != nil {
		log.Fatal(err)
	}
}
