package cli

import (
	"log"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Create configuration file it`s doesn`t exist",
	Run: func(cmd *cobra.Command, args []string) {
		settingsFile, err := cmd.Flags().GetString(flagCfgFile)
		if err != nil {
			log.Fatalf(err.Error())
		}
		configFromYaml.CreateSettingFile(settingsFile)

	},
}

func init() {
	setupCmd.Flags().String(flagCfgFile, cfgFile, "Only config filename.yaml")
	rootCmd.AddCommand(setupCmd)
}
