package cmd

import (
	"log"
	"os"

	"github.com/sgrimee/kizcool"
	"github.com/sgrimee/kizcool/config"
	"github.com/spf13/cobra"
)

var kiz *kizcool.Kiz

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kizcli",
	Short: "Overkiz command-line client",
	Long:  `kizclient implements a partial client for the Overkiz home automation api.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if (len(os.Args) > 1) && (os.Args[1] != "configure") {
		initKizFromConfig()
	}

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Initialise the global kiz from config file
func initKizFromConfig() {
	if err := config.Read(false); err != nil {
		log.Fatal(err)
	}
	k, err := kizcool.New(config.Username(), config.Password(), config.BaseURL(), config.SessionID())
	if err != nil {
		log.Fatal(err)
	}
	kiz = k
}
