// +build integration

package kizcool

import (
	"log"
	"os"
	"testing"

	"github.com/sgrimee/kizcool/api"
	"github.com/sgrimee/kizcool/config"
	"github.com/stretchr/testify/assert"
)

var kiz *Kiz

func TestMain(m *testing.M) {
	err := config.Read(false)
	if err != nil {
		log.Fatal(err)
	}
	kiz, err = New(config.Username(), config.Password(), config.BaseURL(), "")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestIntBadLogin(t *testing.T) {
	k, err := New("baduser", "badpass", config.BaseURL(), "")
	assert.NoError(t, err)
	err = k.Login()
	assert.Error(t, err)
	_, ok := err.(*api.AuthenticationError)
	assert.True(t, ok)
}

func TestIntGoodLogin(t *testing.T) {
	err := kiz.Login()
	assert.NoError(t, err)
}
