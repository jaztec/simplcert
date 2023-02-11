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
			checkVerboseFlag(c)

			p := c.String("root_cert_path")
			if p == "" {
				return fmt.Errorf("empty path parameter")
			}

			m, err := security.NewManager(p)
			if err != nil {
				return err
			}

			cfg, err := certConfigFromFlags(c)
			if err != nil {
				return err
			}

			printConfig(cfg)

			crt, priv, pub, err := m.CreateNamedCert(cfg)
			if err != nil {
				return err
			}

			fmt.Printf("%s\n\n", crt)
			fmt.Printf("%s\n\n", priv)
			fmt.Printf("%s\n\n", pub)

			return err
		},
		Flags: flags(
			nameFlag(),
			hostFlag(),
			certsPathFlag(),
			countryFlag(),
			organizationFlag(),
			isServerFlag(),
			verboseFlag(),
		),
	}
}

func verifyCertsCmd() *cli.Command {
	return &cli.Command{
		Name:     "verify",
		HelpName: "verify if all required certificates are present",
		Action: func(c *cli.Context) error {
			checkVerboseFlag(c)

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
			certsPathFlag(),
			verboseFlag(),
		),
	}
}

func showRootCertCmd() *cli.Command {
	return &cli.Command{
		Name:     "root-crt",
		HelpName: "display the root certificate which should be used as chain verifier",
		Action: func(c *cli.Context) error {
			checkVerboseFlag(c)

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
			certsPathFlag(),
			verboseFlag(),
		),
	}
}

func printConfig(cfg security.CertConfig) {
	log.WithFields(log.Fields{
		"name":         cfg.Name,
		"host":         cfg.Host,
		"is_ca":        cfg.IsCA,
		"is_server":    cfg.IsServer,
		"country":      cfg.Country,
		"organization": cfg.Organization,
	}).Debug("Using config")
}
