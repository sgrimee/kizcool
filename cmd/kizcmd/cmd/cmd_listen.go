package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/sgrimee/kizcool"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for events",
	Long:  "Continuously poll for events from the server and display them on the console.",
	Run: func(cmd *cobra.Command, arge []string) {
		events := make(chan kizcool.Event)
		finish := make(chan struct{})
		e := make(chan error)

		go kiz.PollEventsContinuous(events, e, finish)

		for {
			select {
			case err := <-e:
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Polling error, will resume after a pause.")
			case event := <-events:
				log.WithFields(log.Fields{
					"event": event,
					"type":  fmt.Sprintf("%T", event),
				}).Info("Received event")
			default:
				time.Sleep(time.Millisecond * 100) // avoid burning the CPU
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listenCmd)
}
