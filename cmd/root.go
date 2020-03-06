package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vshn/waf-tool/cfg"
)

var (
	rootCmd = &cobra.Command{
		Use:           "waf-tool",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	config = cfg.NewDefaultConfig()
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// SetVersion sets the version
func SetVersion(version string) {
	rootCmd.Version = version
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("verbose", "v", config.Log.Verbose, "Verbose log output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.WithError(err).Fatal("Could not bind flags")
	}
	viper.SetEnvPrefix("WAF")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		log.WithError(err).Fatal("Could not read config")
	}

	if config.Log.Verbose {
		log.SetLevel(log.DebugLevel)
	}
}
