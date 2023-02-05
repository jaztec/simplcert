package main

import (
	log "github.com/sirupsen/logrus"
	"gitlab.jaztec.info/certs/manager/cmd/cmds"
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
