package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kunalkushwaha/octool/plugins"
	_ "github.com/kunalkushwaha/octool/plugins/linux"
	"github.com/spf13/cobra"
)

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:     "lint",
	Short:   "validate container config file(s)",
	Run:     validateContainerConfig,
	Example: "octool lint <path-of-image>",
}

func init() {
	RootCmd.AddCommand(lintCmd)

	// lintCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	lintCmd.Flags().StringP("image", "", ".", "path of image to validate")
	lintCmd.Flags().StringP("os", "", "linux", "Target OS")

}

func validateContainerConfig(cmd *cobra.Command, args []string) {
	imagePath, _ := cmd.Flags().GetString("image")
	targetOS, _ := cmd.Flags().GetString("os")

	_, err := os.Stat(imagePath)
	if os.IsNotExist(err) {
		cmd.HelpFunc()
		return
	}
	//FIXME: Instead of default as linux, detect os

	plugin, err := plugin.NewPlugin(targetOS)
	if err != nil {
		//fmt.Println(err)
		log.Error(err)
		return
	}
	errors, valid := plugin.ValidatePluginSpecs(imagePath)
	if !valid {
		fmt.Println("")
		for _, err := range errors {
			log.Warn(err)
			//fmt.Println(err)
		}
		fmt.Printf("\nInvalid OCI config format\n")
	} else {
		fmt.Printf("\nConfig is Valid OCI\n")
	}
	return
}
