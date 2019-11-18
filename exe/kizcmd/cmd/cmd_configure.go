package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sgrimee/kizcool"
	"github.com/sgrimee/kizcool/config"
	"github.com/spf13/cobra"
)

func promptCredentials() (username, password string, err error) {
	fmt.Println("Please enter your credentials for the overkiz api.")

	fmt.Print("username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err = reader.ReadString('\n')
	username = strings.TrimSuffix(username, "\n")
	if err != nil {
		return
	}
	fmt.Print("password: ")
	reader = bufio.NewReader(os.Stdin)
	password, err = reader.ReadString('\n')
	password = strings.TrimSuffix(password, "\n")
	return
}

// configureCmd represents the on command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Generates the config file",
	Long:  `Prompts the user for username and password and saves to the config file. If no config file exists it will be created.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Read(true); err != nil {
			log.Fatal(err)
		}
		username, password, err := promptCredentials()
		if err != nil {
			log.Fatal(err)
		}
		config.SetUsername(username)
		config.SetPassword(password)
		kiz, err = kizcool.New(username, password, config.BaseURL(), "")
		if err != nil {
			log.Fatal(err)
		}
		if err := kiz.Login(); err != nil {
			log.Fatal(err)
		}
		if err := config.Write(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Changes written to config file: %s\n", config.File())
	},
}

func init() {
	RootCmd.AddCommand(configureCmd)
}
