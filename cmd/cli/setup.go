package cli

import (
	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Create configuration file it`s doesn`t exist",
	Run: func(cmd *cobra.Command, args []string) {
		settingsFile, err := cmd.Flags().GetString(flagCfgFile)
		if err != nil {
			logger.LogError(errors.Errorf("ngonx: :%v", err).Error())

		}
		configFromYaml.CreateSettingFile(settingsFile)

	},
}

func init() {
	setupCmd.Flags().String(flagCfgFile, cfgFile, "Only config filename.yaml")
	rootCmd.AddCommand(setupCmd)
}
