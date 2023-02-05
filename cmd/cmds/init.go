package cmds

import (
	"github.com/urfave/cli/v2"
	"os"
)

func Run(version string) error {
	a := &cli.App{
		Name:                 "simplcert",
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Usage:                "Control certificates for the infrastructure",
		Version:              version,
		Commands:             commands(),
	}

	return a.Run(os.Args)
}

func commands() []*cli.Command {
	cs := []*cli.Command{
		verifyCertsCmd(),
		createCertCmd(),
		showRootCertCmd(),
	}

	for _, c := range cs {
		for _, s := range c.Subcommands {
			s.Flags = append(s.Flags, verboseFlag())
		}
	}

	return cs
}
