package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/vshn/waf-tool/cmd"
)

var (
	// These will be populated/overridden by Goreleaser
	version = "latest"
	commit  = "dirty"
	date    = "today"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		ForceColors:            true,
	})

	cmd.SetVersion(fmt.Sprintf("%s, commit %s, date %s", version, commit, date))
	if err := cmd.Execute(); err != nil {
		log.WithError(err).Fatal()
	}
}
