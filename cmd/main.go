package main

import (
	"github.com/jaztec/simplcert/cmd/cmds"
	log "github.com/sirupsen/logrus"
)

var (
	Version string
	Build   string
)

func main() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	err := cmds.Run(getVersion())
	if err != nil {
		log.Fatal(err)
	}
}

func getVersion() string {
	v := Version
	if Version != "" && Build != "" {
		v += " - "
	}
	if Build != "" {
		v += Build
	}
	return v
}
