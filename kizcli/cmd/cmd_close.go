package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// closeCmd represents the close command
var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "Close device",
	Long: `Close the device.
	The first argument is the device url or label.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("You must specify a device.")
		}
		dev, err := kiz.GetDeviceByText(args[0])
		if err != nil {
			log.Fatal(err)
		}
		_, err = kiz.Close(dev)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(closeCmd)
}
