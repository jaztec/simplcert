package cmds

import (
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

func certsPathFlag(required bool) appFlag {
	return &cli.StringFlag{
		Name:     "certs_path",
		Usage:    "pathname",
		Required: required,
		EnvVars:  []string{"SCM_CERTS_PATH"},
	}
}

func nameFlag(required bool) appFlag {
	return &cli.StringFlag{
		Name:     "name",
		Aliases:  []string{"n"},
		Usage:    "provide name",
		Required: required,
	}
}

func hostFlag(required bool) appFlag {
	return &cli.StringFlag{
		Name:     "host",
		Usage:    "provide host",
		Required: required,
	}
}

func isServerFlag() appFlag {
	return &cli.BoolFlag{
		Name:  "is_server",
		Usage: "indicate is server",
	}
}

func verboseFlag() appFlag {
	return &cli.BoolFlag{
		Name:    "verbose",
		Aliases: []string{"v"},
		EnvVars: []string{"SCM_VERBOSE"},
	}
}
