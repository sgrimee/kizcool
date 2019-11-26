package main

import (
	"os"

	"github.com/sgrimee/kizcool/cmd/kizcmd/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	//log.SetReportCaller(true)
}

func main() {
	cmd.Execute()
}
