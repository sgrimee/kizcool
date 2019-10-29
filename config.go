package kizcool

import (
	"github.com/spf13/viper"
)

// Config required to connect to a server
type Config struct {
	Username string
	Password string
	BaseURL  string
	Debug    bool
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
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
		BaseURL:  viper.GetString("base_url"),
		Debug:    viper.GetBool("debug"),
	}, nil
}
