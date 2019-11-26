package cmd

import (
	log "github.com/sirupsen/logrus"
	"strconv"

	"github.com/spf13/cobra"
)

// onCmd represents the on command
var intensityCmd = &cobra.Command{
	Use:   "intensity",
	Short: "Set device intensity",
	Long: `Set device to given intensity.
	The first argument is the device url or label.
	The second argument is the intensity in range 0-100`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal("You must specify a device and intensity")
		}
		dev, err := kiz.GetDeviceByText(args[0])
		if err != nil {
			log.Fatal(err)
		}
		intensity, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(err)

		}
		_, err = kiz.SetIntensity(dev, intensity)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(intensityCmd)
}
