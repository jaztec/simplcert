package cmds

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gitlab.jaztec.info/certs/manager/pkg/security"
	"os"
)

func createCertCmd() *cli.Command {
	return &cli.Command{
		Name:     "create",
		HelpName: "create a named certificate",
		Action: func(c *cli.Context) error {
			p := c.String("certs_path")
			if p == "" {
				return fmt.Errorf("empty path parameter")
			}
			n := c.String("name")
			if n == "" {
				return fmt.Errorf("invalid name parameter (%s)", n)
			}
			h := c.String("host")
			if h == "" {
				return fmt.Errorf("invalid host parameter (%s)", h)
			}

			m, err := security.NewManager(p)
			if err != nil {
				return err
			}

			crt, priv, pub, err := m.CreateNamedCert(n, h, c.Bool("is_server"))
			if err != nil {
				return err
			}

			fmt.Printf("%s\n\n", crt)
			fmt.Printf("%s\n\n", priv)
			fmt.Printf("%s\n\n", pub)

			return err
		},
		Flags: flags(
			nameFlag(true),
			hostFlag(true),
			certsPathFlag(true),
			isServerFlag(),
		),
	}
}

func verifyCertsCmd() *cli.Command {
	return &cli.Command{
		Name:     "verify",
		HelpName: "verify if all required certificates are present",
		Action: func(c *cli.Context) error {
			p := c.String("certs_path")
			if p == "" {
				return fmt.Errorf("empty path parameter")
			}
			log.WithField("path", p).Info("Verify certificates")
			if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
				return fmt.Errorf("invalid path parameter (%s)", p)
			}

			_, err := security.NewManager(p)
			return err
		},
		Flags: flags(
			certsPathFlag(true),
		),
	}
}

func showRootCertCmd() *cli.Command {
	return &cli.Command{
		Name:     "root-crt",
		HelpName: "display the root certificate which should be used as chain verifier",
		Action: func(c *cli.Context) error {
			p := c.String("certs_path")
			if p == "" {
				return fmt.Errorf("empty path parameter")
			}
			if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
				return fmt.Errorf("invalid path parameter (%s)", p)
			}

			m, err := security.NewManager(p)
			if err != nil {
				return err
			}

			fmt.Printf("%s\n\n", string(m.RootPEM()))

			return nil
		},
		Flags: flags(
			certsPathFlag(true),
		),
	}
}
