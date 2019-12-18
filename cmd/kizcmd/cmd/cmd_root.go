package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/sgrimee/kizcool"
	"github.com/sgrimee/kizcool/config"
	"github.com/spf13/cobra"
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

	RootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debugging")
	if config.Debug() {
		log.SetLevel(log.DebugLevel)
	}
	log.Debug("Debugging mode")

	if len(os.Args) > 1 {
		switch cmd := os.Args[1]; cmd {
		case "knx":
			// Do not use cached sessionID for the bridge to allow parallel runs with another cli
			kiz = kizFromConfig(false)
		case "config":
			// do not try to initialise a kiz from config as we are generating the config file
		default:
			// use the cached sessionID from the config if present
			kiz = kizFromConfig(true)
		}
	}

	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	// save sessionID to config file
	if config.UsingConfigFile() {
		id := kiz.SessionID()
		config.SetSessionID(id)
		if err := config.Write(); err != nil {
			log.Fatalf("Unable to save session ID to config file: %s\n", err)
		}
	}
}

// Initialise the global kiz from config file
func kizFromConfig(useCachedSessionID bool) *kizcool.Kiz {
	if err := config.Read(false); err != nil {
		log.Fatal(err)
	}
	var sessionID string
	if useCachedSessionID {
		sessionID = config.SessionID()
	}
	k, err := kizcool.New(config.Username(), config.Password(), config.BaseURL(), sessionID)
	if err != nil {
		log.Fatal(err)
	}
	return k
}
