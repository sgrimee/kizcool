package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// onCmd represents the on command
var offCmd = &cobra.Command{
	Use:   "off",
	Short: "Turn device off",
	Long: `Turn device off.
	The first argument is the device url or label.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("You must specify a device.")
		}
		dev, err := kiz.GetDeviceByText(args[0])
		if err != nil {
			log.Fatal(err)
		}
		_, err = kiz.Off(dev)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(offCmd)
}
