package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// GetConfigFromFlags loads config from flags
func GetConfigFromFlags(flags *pflag.FlagSet, cfg interface{}) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := v.BindPFlags(flags); err != nil {
		log.WithError(err).Error("Could not bind flags")
	}
	if err := v.Unmarshal(cfg); err != nil {
		log.WithError(err).Error("Could not unmarshal cfg")
	}
}
