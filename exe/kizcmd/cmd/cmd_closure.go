package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// onCmd represents the on command
var closureCmd = &cobra.Command{
	Use:   "closure",
	Short: "Set device closure",
	Long: `Set device to given closure.
	The first argument is the device url or label.
	The second argument is the closure in range 0-100`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("You must specify a device and closure")
		}
		dev, err := kiz.GetDeviceByText(args[0])
		if err != nil {
			log.Fatal(err)
		}
		closure, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(err)

		}
		_, err = kiz.SetClosure(dev, closure)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(closureCmd)
}
