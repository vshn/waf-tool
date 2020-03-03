package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vshn/waf-tool/cfg"
	"github.com/vshn/waf-tool/cmd"
	"os"
	"strings"
)

var (
	// These will be populated/overridden by Goreleaser
	version = "latest"
	commit  = "dirty"
	date    = "today"

	rootCmd = &cobra.Command{
		Use:     "waf-tool",
		Version: fmt.Sprintf("%s, commit %s, date %s", version, commit, date),
	}
)

func initialize() {

	defaults := cfg.CreateDefaultConfig()

	rootCmd.PersistentFlags().String("log.level", defaults.Log.Level, "Log level")
	rootCmd.AddCommand(cmd.CreateDemoCommand())
	flags := rootCmd.PersistentFlags()
	if err := viper.BindPFlags(flags); err != nil {
		log.Fatal(err)
	}
	if err := flags.Parse(flags.Args()); err != nil {
		log.Fatal(err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Overwrite defaults with new variables from flags as well as ENV variables:
	err := viper.Unmarshal(&defaults)
	if err != nil {
		log.Fatal(err)
	}
	cfg.SetupLogging(defaults.Log)
}

func main() {
	cobra.OnInitialize(initialize)

	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Error("Command error.")
		os.Exit(1)
	}
}
