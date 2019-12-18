// kizknx is a two-way bridge between a knx network and an overkiz setup, for pre-configured events
package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sgrimee/kizcool"
	"github.com/sgrimee/kizcool/config"
	"github.com/sgrimee/kizcool/group"

	"github.com/sgrimee/knx-go/knx"

	"github.com/sgrimee/knx-go/knx/util"
)

var kiz *kizcool.Kiz

func init() {
	if err := config.Read(false); err != nil {
		log.Fatal(err)
	}
	if config.Debug() {
		log.SetLevel(log.DebugLevel)
	}
	util.Logger = log.New()
	log.Debug("Debugging mode")
}

func main() {
	// We do not use the sessionId stored in the config to avoid sharing the session with a kizcmd
	// that may be running in parallel. Since this is a long-running server it is not needed.
	_kiz, err := kizcool.New(config.Username(), config.Password(), config.BaseURL(), "")
	if err != nil {
		log.Fatal(err)
	}
	kiz = _kiz

	// Listen for overkiz events
	ovkEvent := make(chan kizcool.Event)
	ovkErr := make(chan error)
	ovkFinish := make(chan struct{})
	go kiz.PollEventsContinuous(ovkEvent, ovkErr, ovkFinish)

	// Listen for knx events
	knxClt, err := knx.NewGroupRouter("224.0.23.12:3671", knx.DefaultRouterConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer knxClt.Close()
	knxInbound := knxClt.Inbound()

	// Process events
	for {
		select {
		case kizEvent, open := <-ovkEvent:
			if !open {
				log.Fatal("Kiz channel closed")
			}
			processKizEvent(kizEvent)
		case knxEvent, open := <-knxInbound:
			if !open {
				log.Fatal("Knx channel closed")
			}
			processKnxMsg(knxEvent)
		case err := <-ovkErr:
			log.WithFields(log.Fields{"err": err}).Error("Kiz polling error, will resume after a pause.")
		default:
			time.Sleep(time.Millisecond * 100) // avoid burning the CPU
		}
	}
}

func processKizEvent(kizEvent kizcool.Event) error {
	log.WithFields(log.Fields{
		"kind":  kizEvent.Kind(),
		"event": kizEvent,
	}).Debug("Kiz event")
	return nil
}

func processKnxMsg(knxEvent knx.GroupEvent) error {
	log.WithFields(log.Fields{
		"Command":     knxEvent.Command,
		"Destination": knxEvent.Destination,
		"Source":      knxEvent.Source,
		"Value":       knxEvent.Data,
	}).Debug("KNX event")

	if knxEvent.Command != knx.GroupWrite {
		return nil
	}

	for _, d := range managedDevices() {
		for _, cmd := range d.Commands {
			if cmd.GroupAddr == knxEvent.Destination {
				log.WithFields(log.Fields{
					"Label":       d.Label,
					"CommandName": cmd.Name,
				}).Debug("KNX command")
				processCommand(cmd, d, knxEvent.Data)
				return nil
			}
		}
	}
	return nil
}

func processCommand(gcmd group.Command, gd group.Device, data []byte) {
	device, err := kiz.GetDevice(kizcool.DeviceURL(gd.URL))
	if err != nil {
		log.Fatal(err)
	}
	switch gcmd.Name {
	case "setOnOff":
		if len(data) < 1 {
			log.WithFields(log.Fields{
				"gcmd":   gcmd,
				"device": gd,
				"data":   data,
			}).Error("Invalid data field")
			return
		}
		if data[0] == 1 {
			log.Infof("Turning device %s on\n", device.Label)
			kiz.On(device)
		} else {
			log.Infof("Turning device %s off\n", device.Label)
			kiz.Off(device)
		}
	case "setIntensity":
		if len(data) < 2 {
			log.WithFields(log.Fields{
				"gcmd":   gcmd,
				"device": gd,
				"data":   data,
			}).Error("Invalid data field")
			return
		}
		intensity := int(data[1]) * 100 / 255
		log.Infof("Setting device %s intensity to %d\n", device.Label, intensity)
		kiz.SetIntensity(device, intensity)
	default:
		log.Warn("Unhandled command")
		return
	}
}

// managedDevices simulates reading a device list from config
func managedDevices() []group.Device {
	devices, err := config.Devices()
	if err != nil {
		log.Fatal(err)
	}
	return devices
}
