package cmd

import (
	"log"
	"regexp"

	"github.com/sgrimee/kizcool"
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
		var deviceURL kizcool.DeviceURLT
		validURL := regexp.MustCompile(`^[a-z]+://\d{4}-\d{4}-\d{4}/\d+`)
		if validURL.MatchString(args[0]) {
			// a DeviceURL was given
			deviceURL = kizcool.DeviceURLT(args[0])
		} else {
			// try to match a Label from all devices
			devices, err := kiz.GetDevices()
			if err != nil {
				log.Fatal(err)
			}
			deviceURL, err = kizcool.DeviceURLByLabel(args[0], devices)
			if err != nil {
				log.Fatalf("Unable to identify device with url or label: %s", args[0])
			}
		}
		dev, err := kiz.GetDevice(deviceURL)
		if err != nil {
			log.Fatal(err)
		}
		output(outputFormat, dev)
	},
}

func init() {
	getCmd.AddCommand(deviceCmd)
}
