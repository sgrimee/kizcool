// Package config provides handling of configuration through a config file
package config

import (
	"fmt"
	"reflect"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/sgrimee/kizcool/config/knxcfg"
	"github.com/spf13/viper"
	"github.com/vapourismo/knx-go/knx/cemi"
)

// Defaults for config file
const (
	DefaultConfigFileBaseName  = ".kizcool"
	DefaultConfigFileExtension = ".yaml"
	TahomaBaseURL              = "https://tahomalink.com/enduser-mobile-web"
)

var usingConfigFile bool

var v *viper.Viper

func init() {
	v = viper.New()
	v.SetDefault("base_url", TahomaBaseURL)

	v.SetEnvPrefix("KIZ")
	v.AutomaticEnv()
}

// Username returns the username
func Username() string {
	return v.GetString("username")
}

// SetUsername sets the username
func SetUsername(username string) {
	v.Set("username", username)
}

// Password returns the password
func Password() string {
	return v.GetString("password")
}

// SetPassword sets the password
func SetPassword(password string) {
	v.Set("password", password)
}

// BaseURL returns the BaseURL
func BaseURL() string {
	return v.GetString("base_url")
}

// SetBaseURL sets the BaseURL
func SetBaseURL(url string) {
	v.Set("base_url", url)
}

// SessionID returns the SessionID
func SessionID() string {
	return v.GetString("session_id")
}

// SetSessionID stores the session ID in the configuration file
func SetSessionID(ID string) {
	v.Set("session_id", ID)
}

// toGaddrHookFunc decodes strings into cemi.GroupAddr integers
func toGaddrHookFunc() mapstructure.DecodeHookFunc {
	return func(
		from reflect.Type,
		to reflect.Type,
		data interface{}) (interface{}, error) {
		// Only handle conversion from string to Uint16 (cemi.GroupAddr)
		if to.Kind() != reflect.Uint16 {
			return data, nil
		}
		if from.Kind() != reflect.String {
			return data, nil
		}
		str := fmt.Sprintf("%v", data)
		return cemi.NewGroupAddrString(str)
	}
}

// Devices returns the list of group devices defined in the config
func Devices() ([]knxcfg.Device, error) {
	var result []knxcfg.Device
	opt := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		toGaddrHookFunc(),
	))

	err := v.UnmarshalKey("devices", &result, opt)
	return result, err
}

// SetDevices saves the list of devices to the config
func SetDevices(devices []knxcfg.Device) {
	v.Set("devices", devices)
}

// Read reads in config file. It should be called before using other functions in this package.
// The local directory is searched first, then the user's home directory
// If no file is found and create is true, a config file with defaults is created.
func Read(create bool) error {
	v.SetConfigName(DefaultConfigFileBaseName) // name of config file (without extension)
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")

	if err := v.ReadInConfig(); err != nil {
		// create empty config file if needed
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && create {
			configFile, _ := homedir.Expand("~/" + DefaultConfigFileBaseName + DefaultConfigFileExtension)
			v.SetConfigPermissions(0600)
			if err := v.SafeWriteConfigAs(configFile); err != nil {
				usingConfigFile = false
				return err
			}
			usingConfigFile = true
			return v.ReadInConfig()
		}
		usingConfigFile = false
		return nil
	}
	usingConfigFile = true
	return nil
}

// Write saves the current config to file
func Write() error {
	return v.WriteConfig()
}

// UsingConfigFile is true if a config file is used, false otherwise (e.g. only env variables)
func UsingConfigFile() bool {
	return usingConfigFile
}

// File returns the name and path of the configuration file in use.
func File() string {
	return v.ConfigFileUsed()
}

// Debug returns true if debugging was requested
func Debug() bool {
	return v.GetBool("debug")
}
