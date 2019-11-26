package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// openCmd represents the on command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open device",
	Long: `Open the device.
	The first argument is the device url or label.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("You must specify a device.")
		}
		dev, err := kiz.GetDeviceByText(args[0])
		if err != nil {
			log.Fatal(err)
		}
		_, err = kiz.Open(dev)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(openCmd)
}
