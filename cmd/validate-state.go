package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/kunalkushwaha/octool/plugins"
	_ "github.com/kunalkushwaha/octool/plugins/linux"
	"github.com/spf13/cobra"
)

// validate-stateCmd represents the validate-state command
var validateStateCmd = &cobra.Command{
	Use:   "validate-state",
	Short: "Validates the Container state",
	Run:   validateContainerState,
}

func init() {
	RootCmd.AddCommand(validateStateCmd)

	validateStateCmd.Flags().StringP("id", "", "", "Container-ID")
}

func validateContainerState(cmd *cobra.Command, args []string) {
	containerID, _ := cmd.Flags().GetString("id")

	if len(containerID) == 0 {
		cmd.HelpFunc()
		return
	}

	//FIXME: Instead of default as linux, detect os
	targetOS := "linux"

	plugin, err := plugin.NewPlugin(targetOS)
	if err != nil {
		log.Error(err)
		//fmt.Println(err)
		return
	}
	errors, valid := plugin.ValidatePluginState(containerID)
	if !valid {
		for _, err := range errors {
			//fmt.Println(err)
			log.Warn(err)
		}
		fmt.Printf("\nInvalid OCI State format\n")
	} else {
		fmt.Println("Container State Valid OCI")
	}

	return
}
