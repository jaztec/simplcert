package cmds

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type appFlag cli.Flag

func flags(f ...appFlag) []cli.Flag {
	r := make([]cli.Flag, 0, len(f))
	for _, fl := range f {
		r = append(r, fl)
	}
	return r
}

func certsPathFlag() appFlag {
	return &cli.StringFlag{
		Name:     "root_cert_path",
		Usage:    "pathname",
		Required: true,
		EnvVars:  []string{"SCM_ROOT_CERT_PATH"},
	}
}

func nameFlag() appFlag {
	return &cli.StringFlag{
		Name:    "name",
		Aliases: []string{"n"},
		Usage:   "provide name",
	}
}

func hostFlag() appFlag {
	return &cli.StringFlag{
		Name:  "host",
		Usage: "provide host",
	}
}

func isServerFlag() appFlag {
	return &cli.BoolFlag{
		Name:  "is_server",
		Usage: "indicate if target is server",
	}
}

func verboseFlag() appFlag {
	return &cli.BoolFlag{
		Name:    "verbose",
		Aliases: []string{"v"},
		EnvVars: []string{"SCM_VERBOSE"},
	}
}

func countryFlag() appFlag {
	return &cli.StringFlag{
		Name: "country",
	}
}

func organizationFlag() appFlag {
	return &cli.StringFlag{
		Name: "organization",
	}
}

func checkVerboseFlag(c *cli.Context) {
	l := log.InfoLevel
	if v := c.Bool("verbose"); v {
		l = log.DebugLevel
	}
	log.SetLevel(l)
}
