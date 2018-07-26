package cmd

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/tok"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"

	"github.com/spf13/cobra"
)

// InitCmd - `tok init`
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Tokaido project",
	Long:  "Initialize a Tokaido project. Installs dependencies, starts unison synchronization, builds a configuration file, starts your container. Docker is required to be installed manually for this to work",
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckCmdHard("docker-compose")
		conf.LoadConfig(cmd)

		fmt.Printf("The command 'tok init' has been deprecated. Please use 'tok up' instead.\n\n")

		tok.Init()
		tok.InitMessage()
	},
}
