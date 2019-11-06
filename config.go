package kizcool

import (
	"github.com/spf13/viper"
)

// Config required to connect to a server
type Config struct {
	Username  string
	Password  string
	BaseURL   string
	Debug     bool
	SessionID string
}

// GetConfig reads in config file and returns a Config struct
func GetConfig() (Config, error) {
	viper.SetConfigName(".kizcool") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.SetDefault("base_url", "https://tahomalink.com/enduser-mobile-web")
	viper.SetDefault("debug", false)

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	return Config{
		Username:  viper.GetString("username"),
		Password:  viper.GetString("password"),
		BaseURL:   viper.GetString("base_url"),
		Debug:     viper.GetBool("debug"),
		SessionID: viper.GetString("session_id"),
	}, nil
}

// setConfigValue updates an entry in the config struct, in the
// viper singleton and saves it back to the configuration file.
func setConfigValue(key string, value interface{}) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}

// SaveSessionID stores the session ID in the configuration file
func SaveSessionID(ID string) error {
	return setConfigValue("session_id", ID)
}
