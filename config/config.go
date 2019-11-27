// Package config provides handling of configuration through a config file
package config

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Defaults for config file
const (
	DefaultConfigFileBaseName  = ".kizcool"
	DefaultConfigFileExtension = ".yaml"
	TahomaBaseURL              = "https://tahomalink.com/enduser-mobile-web"
)

var usingConfigFile bool

// Username returns the username
func Username() string {
	return viper.GetString("username")
}

// SetUsername sets the username
func SetUsername(username string) {
	viper.Set("username", username)
}

// Password returns the password
func Password() string {
	return viper.GetString("password")
}

// SetPassword sets the password
func SetPassword(password string) {
	viper.Set("password", password)
}

// BaseURL returns the BaseURL
func BaseURL() string {
	return viper.GetString("base_url")
}

// SetBaseURL sets the BaseURL
func SetBaseURL(url string) {
	viper.Set("base_url", url)
}

// SessionID returns the SessionID
func SessionID() string {
	return viper.GetString("session_id")
}

// SetSessionID stores the session ID in the configuration file
func SetSessionID(ID string) {
	viper.Set("session_id", ID)
}

// Read reads in config file. It should be called before using other functions in this package.
// The local directory is searched first, then the user's home directory
// If no file is found and create is true, a config file with defaults is created.
func Read(create bool) error {
	viper.SetConfigName(DefaultConfigFileBaseName) // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.SetDefault("base_url", TahomaBaseURL)

	viper.SetEnvPrefix("KIZ")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// create empty config file if needed
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && create {
			configFile, _ := homedir.Expand("~/" + DefaultConfigFileBaseName + DefaultConfigFileExtension)
			viper.SetConfigPermissions(0600)
			if err := viper.SafeWriteConfigAs(configFile); err != nil {
				usingConfigFile = false
				return err
			}
			usingConfigFile = true
			return viper.ReadInConfig()
		}
		usingConfigFile = false
		return nil
	}
	usingConfigFile = true
	return nil
}

// UsingConfigFile is true if a config file is used, false otherwise (e.g. only env variables)
func UsingConfigFile() bool {
	return usingConfigFile
}

// Write saves the current config to file
func Write() error {
	return viper.WriteConfig()
}

// File returns the name and path of the configuration file in use.
func File() string {
	return viper.ConfigFileUsed()
}

// Debug returns true if debugging was requested
func Debug() bool {
	return viper.GetBool("debug")
}
