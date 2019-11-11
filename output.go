package kizcool

import (
	"encoding/json"
	"fmt"
	"io"

	yaml "gopkg.in/yaml.v2"
)

// Output prints the given object to the writer in the desired format
// format can be 'text', 'json' or 'yaml'
func Output(w io.Writer, format string, obj interface{}) error {
	switch format {
	case "yaml":
		return printYAML(w, obj)
	case "json":
		return printJSON(w, obj)
	case "text":
		return printText(w, obj)
	default:
		return fmt.Errorf("Unknown output format: %s", format)
	}
}

func printJSON(w io.Writer, obj interface{}) error {
	j, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err = io.WriteString(w, fmt.Sprintln(string(j))); err != nil {
		return err
	}

	return nil
}

func printYAML(w io.Writer, obj interface{}) error {
	y, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, fmt.Sprintf("---\n%s\n", string(y)))
	return err
}

// printText prints the object most useful fields in readable text format
func printText(w io.Writer, obj interface{}) (err error) {
	switch t := obj.(type) {
	default:
		return fmt.Errorf("printText does not support type %T", t)
	case []string:
		ls := obj.([]string)
		for _, s := range ls {
			if _, err = io.WriteString(w, s+"\n"); err != nil {
				return err
			}
		}
	case string:
		s := obj.(string)
		if _, err = io.WriteString(w, s); err != nil {
			return err
		}
	case Device:
		err := printTextDevice(w, obj.(Device))
		if err != nil {
			return err
		}
	case []Device:
		for _, d := range obj.([]Device) {
			err := printTextDevice(w, d)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// printTextDevice prints useful values of a single device
func printTextDevice(w io.Writer, d Device) (err error) {
	wantedStates := map[StateName]bool{
		"core:ClosureState":        true,
		"core:OpenClosedState":     true,
		"core:LightIntensityState": true,
		"core:OnOffState":          true,
		// "core:RSSILevelState":      true,
	}
	var states []State
	for _, state := range d.States {
		if wantedStates[state.Name] {
			states = append(states, state)
		}
	}
	if _, err = io.WriteString(w, fmt.Sprintf("| %-22s | %-33s | %v |\n",
		d.Label, d.DeviceURL, states)); err != nil {
		return err
	}
	return nil
}
