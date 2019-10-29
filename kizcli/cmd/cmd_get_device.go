package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Get a single device",
	Long: `Get infos on a single device identified by its device url or label (case insensitive)
	kizclient get device "io://1111-0000-4444/15332221"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("You must specify a device.")
		}
		dev, err := kiz.GetDeviceByText(args[0])
		if err != nil {
			log.Fatal(err)
		}
		output(outputFormat, dev)
	},
}

func init() {
	getCmd.AddCommand(deviceCmd)
}
