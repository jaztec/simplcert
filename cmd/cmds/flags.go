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
		Name:    "root-cert-path",
		Usage:   "pathname",
		EnvVars: []string{"SCM_ROOT_CERT_PATH"},
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
		Usage: "provide host, comma separated for multiple inputs. IP addresses allowed",
	}
}

func isServerFlag() appFlag {
	return &cli.BoolFlag{
		Name:  "is-server",
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

func outputPathFlag() appFlag {
	return &cli.StringFlag{
		Name:    "output-path",
		Aliases: []string{"o"},
	}
}

func daysValidFlag() appFlag {
	return &cli.IntFlag{
		Name: "days-valid",
	}
}

func outputName() appFlag {
	return &cli.StringFlag{
		Name: "output-name",
	}
}

func ecdsaFlag() appFlag {
	return &cli.BoolFlag{
		Name: "ecdsa",
	}
}

func rsaFlag() appFlag {
	return &cli.BoolFlag{
		Name: "rsa",
	}
}

func ed25519Flag() appFlag {
	return &cli.BoolFlag{
		Name: "ed25519",
	}
}

func checkVerboseFlag(c *cli.Context) {
	l := log.InfoLevel
	if v := c.Bool("verbose"); v {
		l = log.DebugLevel
	}
	log.SetLevel(l)
}
