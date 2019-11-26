package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get all devices",
	Long:  "Get list of all devices in the installation.",
	Run: func(cmd *cobra.Command, arge []string) {
		devices, err := kiz.GetDevices()
		if err != nil {
			log.Fatal(err)
		}
		output(outputFormat, devices)
	},
}

func init() {
	getCmd.AddCommand(devicesCmd)
}
