// +build integration

package kizcool

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var kiz *Kiz

var config Config

func TestMain(m *testing.M) {
	cfg, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	config = cfg
	os.Exit(m.Run())
}

func TestIntegrationBadLogin(t *testing.T) {
	badConfig := config
	badConfig.Username = "baduser"
	badConfig.Password = "badpass"
	kiz, _ := New(badConfig)
	err := kiz.Login()
	assert.EqualError(t, err, "401: Authentication error")
}

func TestIntegrationGoodLogin(t *testing.T) {
	kiz, err := New(config)
	assert.Nil(t, err)
	err = kiz.Login()
	assert.Nil(t, err)
}
