package knxbridge

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sgrimee/kizcool"
	"github.com/sgrimee/kizcool/config/knxcfg"
	"github.com/sgrimee/knx-go/knx"
)

// A Bridge forwards events between a knx network and a kiz system
type Bridge struct {
	kiz     *kizcool.Kiz
	devices []knxcfg.Device
}

// New returns a new Bridge from an initiated Kiz and a group address device config
func New(k *kizcool.Kiz, d []knxcfg.Device) *Bridge {
	return &Bridge{
		kiz:     k,
		devices: d,
	}
}

// Start causes the bridge to listen on both systems and forward packets. It is blocking.
func (br *Bridge) Start(finish chan struct{}) error {
	// Listen for overkiz events
	ovkEvent := make(chan kizcool.Event)
	ovkErr := make(chan error)
	ovkFinish := make(chan struct{})
	go br.kiz.PollEventsContinuous(ovkEvent, ovkErr, ovkFinish)

	// Listen for knx events
	knxClt, err := knx.NewGroupRouter("224.0.23.12:3671", knx.DefaultRouterConfig)
	if err != nil {
		return err
	}
	defer knxClt.Close()
	knxInbound := knxClt.Inbound()

	// Process events
	for {
		select {
		case kizEvent, open := <-ovkEvent:
			if !open {
				return errors.New("Kiz channel closed")
			}
			br.processKizEvent(kizEvent)
		case knxEvent, open := <-knxInbound:
			if !open {
				return errors.New("Knx channel closed")
			}
			br.processKnxEvent(knxEvent)
		case err := <-ovkErr:
			log.WithFields(log.Fields{"err": err}).Error("Kiz polling error, will resume after a pause.")
		case <-finish:
			log.Debug("Bridge termination was requested")
			return nil
		default:
			time.Sleep(time.Millisecond * 100) // avoid burning the CPU
		}
	}
}

func (br *Bridge) processKizEvent(kizEvent kizcool.Event) error {
	log.WithFields(log.Fields{
		"kind":  kizEvent.Kind(),
		"event": kizEvent,
	}).Debug("Kiz event")
	// TODO: Bridge device state change events to knx
	return nil
}

func (br *Bridge) processKnxEvent(knxEvent knx.GroupEvent) error {
	log.WithFields(log.Fields{
		"Command":     knxEvent.Command,
		"Destination": knxEvent.Destination,
		"Source":      knxEvent.Source,
		"Value":       knxEvent.Data,
	}).Debug("KNX event")

	// TODO: support GroupRead commands
	if knxEvent.Command != knx.GroupWrite {
		return nil
	}

	for _, d := range br.devices {
		for _, gcmd := range d.Commands {
			if gcmd.GroupAddr == knxEvent.Destination {
				log.WithFields(log.Fields{
					"Label":       d.Label,
					"CommandName": gcmd.Name,
				}).Debug("KNX command")
				// TODO: avoid retrieving all device details each time
				kizDevice, err := br.kiz.GetDevice(kizcool.DeviceURL(d.URL))
				if err != nil {
					return fmt.Errorf("Could not GetDevice with url %s: %w", d.URL, err)
				}
				br.processCommand(gcmd, kizDevice, knxEvent.Data)
				return nil
			}
		}
	}
	return nil
}

func (br *Bridge) processCommand(gcmd knxcfg.Command, d kizcool.Device, data []byte) error {
	switch gcmd.Name {
	case "setOnOff":
		if len(data) < 1 {
			err := errors.New("Invalid data field for setOnOff")
			log.WithFields(log.Fields{
				"gcmd":   gcmd,
				"device": d,
				"data":   data,
			}).Error(err)
			return err
		}
		if data[0] == 1 {
			log.Infof("Turning device %s on\n", d.Label)
			br.kiz.On(d)
		} else {
			log.Infof("Turning device %s off\n", d.Label)
			br.kiz.Off(d)
		}
	case "setIntensity":
		if len(data) < 2 {
			err := errors.New("Invalid data field for setIntensity")
			log.WithFields(log.Fields{
				"gcmd":   gcmd,
				"device": d,
				"data":   data,
			}).Error(err)
			return err
		}
		intensity := int(data[1]) * 100 / 255
		log.Infof("Setting device %s intensity to %d\n", d.Label, intensity)
		br.kiz.SetIntensity(d, intensity)
	default:
		return errors.New("Unhandled command")
	}
	return nil
}
