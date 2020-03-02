package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var developer string

var rootCmd = &cobra.Command{
	Use:   "waf-tuning",
	Short: "waf-tuning start analysing waf log",
	Long:  `waf-tuning start analysing waf log`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&developer, "verbosity", "v", "info", "Log level to use")
}

func initConfig() {
	developer, _ := rootCmd.Flags().GetString("developer")
	if developer != "" {
		fmt.Println("Developer:", developer)
	}
}
