package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// loginCmd represents the on command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Force a login to the server",
	Long:  `Authenticate to the server to get a sessionID. Normally this is not needed.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := kiz.Login(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
