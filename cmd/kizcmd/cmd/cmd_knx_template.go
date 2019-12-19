package cmd

import (
	"github.com/sgrimee/kizcool/config"
	"github.com/sgrimee/kizcool/config/knxcfg"
	"github.com/sgrimee/kizcool/knxbridge"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var knxTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Generate a template for a knx group address configuration file",
	Long:  `Generate a template for a knx group address configuration file from the list of kiz devices"`,
	Run: func(cmd *cobra.Command, args []string) {

		// Use a map to merge devices from kiz with devices already in config file
		// Devices already in the config are kept there and ignored (not updated)
		dm := make(map[string]knxcfg.Device)

		cfgDevices, err := config.Devices()
		if err != nil {
			log.Fatal(err)
		}
		for _, cd := range cfgDevices {
			url := cd.URL
			if url == "" {
				log.Fatalf("Empty URL for device: %+v\n", cd)
			}
			dm[url] = cd
		}
		log.Tracef("config before: %+v", dm)

		kizDevices, err := kiz.GetDevices()
		if err != nil {
			log.Fatal(err)
		}
		for _, kd := range kizDevices {
			url := string(kd.DeviceURL)
			if _, ok := dm[url]; ok {
				continue // ignore devices already in the config
			}
			d := knxcfg.Device{
				URL:   url,
				Label: kd.Label,
			}

			for _, cdef := range kd.Definition.Commands {
				name := knxbridge.ConfigNameForKizCommand(cdef.CommandName)
				if name == "" {
					log.Tracef("Ignoring unsupported command %s\n", cdef.CommandName)
					continue
				}
				if cmdAlreadyThere(d.Commands, name) {
					log.Tracef("Ignoring command %s, already present\n", name)
					continue
				}
				c := knxcfg.Command{
					Name:      name,
					GroupAddr: 0,
				}
				d.Commands = append(d.Commands, c)
			}

			for _, sdef := range kd.Definition.States {
				name := knxbridge.ConfigNameForKizState(sdef.QualifiedName)
				if name == "" {
					log.Tracef("Ignoring unsupported state %s\n", sdef.QualifiedName)
					continue
				}
				if stateAlreadyThere(d.States, name) {
					log.Tracef("Ignoring state %s, already present\n", name)
					continue
				}
				s := knxcfg.State{
					Name:      name,
					GroupAddr: 0,
				}
				d.States = append(d.States, s)
			}

			if (len(d.Commands) > 0) || (len(d.States) > 0) {
				dm[url] = d
			}
		}
		log.Tracef("config after: %+v", dm)

		// In config, devices are saved as a list
		// TODO: consider saving as a map
		var devList []knxcfg.Device
		for _, v := range dm {
			devList = append(devList, v)
		}
		config.SetDevices(devList)
	},
}

// cmdAlreadyThere returns true if command name s is in list of commands l, false otherwise
func cmdAlreadyThere(l []knxcfg.Command, s string) bool {
	for _, v := range l {
		if v.Name == s {
			return true
		}
	}
	return false
}

// stateAlreadyThere returns true if command name s is in list of commands l, false otherwise
func stateAlreadyThere(l []knxcfg.State, s string) bool {
	for _, v := range l {
		if v.Name == s {
			return true
		}
	}
	return false
}

func init() {
	knxCmd.AddCommand(knxTemplateCmd)
}
