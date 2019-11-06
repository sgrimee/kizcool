package cmd

import (
	"log"

	"github.com/sgrimee/kizcool"
)

var kiz *kizcool.Kiz

func init() {
	config, err := kizcool.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	k, err := kizcool.New(config)
	if err != nil {
		log.Fatal(err)
	}
	kiz = k
}
