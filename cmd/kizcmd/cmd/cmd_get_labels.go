package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Get labels from all devices",
	Long:  "Get list of all device labels in the installation.",
	Run: func(cmd *cobra.Command, arge []string) {
		var labels []string

		devices, err := kiz.GetDevices()
		if err != nil {
			log.Fatal(err)
		}
		for _, d := range devices {
			labels = append(labels, d.Label)
		}
		output(outputFormat, labels)
	},
}

func init() {
	getCmd.AddCommand(labelsCmd)
}
