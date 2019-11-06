package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// loginCmd represents the on command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the server",
	Long:  `Authenticate to the server to get a sessionID`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := kiz.Login(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
}
