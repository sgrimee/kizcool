package cmd

import (
	"github.com/spf13/cobra"
)

// knxCmd represents the knx command
var knxCmd = &cobra.Command{
	Use:   "knx",
	Short: "Operate on knx network",
	Long:  "Operate on knx network",
}

func init() {
	RootCmd.AddCommand(knxCmd)
}
