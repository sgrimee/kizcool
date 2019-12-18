package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/sgrimee/kizcool/config"
	"github.com/sgrimee/kizcool/knxbridge"
	"github.com/sgrimee/knx-go/knx/util"
	"github.com/spf13/cobra"
)

// knxBridgeCmd represents the knx bridge command
var knxBridgeCmd = &cobra.Command{
	Use:   "bridge",
	Short: "Bridge kiz events to a knx network",
	Long:  "Bridge kiz events to a knx network",
	Run: func(cmd *cobra.Command, arge []string) {
		if config.Debug() {
			util.Logger = log.New() // see logs from knx package
		}

		devices, err := config.Devices()
		if err != nil {
			log.Fatal(err)
		}

		bridge := knxbridge.New(kiz, devices)
		if err := bridge.Start(nil); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	knxCmd.AddCommand(knxBridgeCmd)
}
