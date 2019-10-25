// +build integration

package kizcool

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var kiz *Kiz

var (
	validUsername string
	validPassword string
)

func TestMain(m *testing.M) {
	initTestIntegrationConfig()
	os.Exit(m.Run())
}

// initConfig reads in config file
func initTestIntegrationConfig() {
	viper.SetConfigName(".kizcool") // name of config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	validUsername = viper.GetString("username")
	validPassword = viper.GetString("password")
}
func TestIntegrationBadLogin(t *testing.T) {
	kiz := NewKiz("baduser", "badpass")
	err := kiz.Login()
	assert.EqualError(t, err, "401: Invalid credentials")
}

func TestIntegrationGoodLogin(t *testing.T) {
	kiz := NewKiz(validUsername, validPassword)
	err := kiz.Login()
	assert.Nil(t, err)
}
