package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop device activity",
	Long: `Stop the current activity of the device.
	The first argument is the device url or label.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("You must specify a device.")
		}
		dev, err := kiz.GetDeviceByText(args[0])
		if err != nil {
			log.Fatal(err)
		}
		_, err = kiz.Stop(dev)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(stopCmd)
}
