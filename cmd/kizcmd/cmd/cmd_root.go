package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/sgrimee/kizcool"
	"github.com/sgrimee/kizcool/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var kiz *kizcool.Kiz

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kizcmd",
	Short: "Overkiz command-line client",
	Long:  `kizcmd implements a partial client for the Overkiz home automation api.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if (len(os.Args) > 1) && (os.Args[1] != "configure") {
		initKizFromConfig()
	}

	RootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debugging")
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	if config.UsingConfigFile() {
		id := kiz.SessionID()
		config.SetSessionID(id)
		if err := config.Write(); err != nil {
			log.Fatalf("Unable to save session ID to config file: %s\n", err)
		}
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
