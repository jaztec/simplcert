package cmds

import (
	"fmt"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gitlab.jaztec.info/certs/manager"
	"os"
)

func createCertCmd() *cli.Command {
	return &cli.Command{
		Name:     "create",
		HelpName: "create a named certificate",
		Action: func(c *cli.Context) error {
			checkVerboseFlag(c)
			p, err := promptRootCertPath(c)
			if err != nil {
				return err
			}

			m, err := manager.NewManager(p)
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

			if cfg.OutputPath != "" {
				name := cfg.OutputName
				if name == "" {
					name = strcase.ToKebab(cfg.Name)
				}
				err = outputToFile(cfg.OutputPath, name, m.RootPEM(), crt, priv, pub)
			} else {
				outputToScreen(crt, priv, pub)
			}

			return err
		},
		Flags: flags(
			nameFlag(),
			hostFlag(),
			certsPathFlag(),
			countryFlag(),
			organizationFlag(),
			isServerFlag(),
			outputPath(),
			outputName(),
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

			p, err := promptRootCertPath(c)
			if err != nil {
				return err
			}

			log.WithField("path", p).Info("Verify certificates")
			if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
				return fmt.Errorf("invalid path parameter (%s)", p)
			}

			_, err = manager.NewManager(p)
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

			p, err := promptRootCertPath(c)
			if err != nil {
				return err
			}

			if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
				return fmt.Errorf("invalid path parameter (%s)", p)
			}

			m, err := manager.NewManager(p)
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

func printConfig(cfg manager.CertConfig) {
	log.WithFields(log.Fields{
		"name":         cfg.Name,
		"host":         cfg.Host,
		"is-server":    cfg.IsServer,
		"country":      cfg.Country,
		"organization": cfg.Organization,
	}).Debug("Using config")
}

func outputToScreen(crt, key, pub []byte) {
	fmt.Printf("%s\n\n", crt)
	fmt.Printf("%s\n\n", key)
	fmt.Printf("%s\n\n", pub)
}

func outputToFile(path, name string, root, crt, key, pub []byte) error {
	sep := string(os.PathSeparator)
	if err := os.WriteFile(fmt.Sprintf("%s%s%s.crt", path, sep, "root-ca"), root, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(fmt.Sprintf("%s%s%s.crt", path, sep, name), crt, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(fmt.Sprintf("%s%s%s.key", path, sep, name), key, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(fmt.Sprintf("%s%s%s.pub", path, sep, name), pub, 0644); err != nil {
		return err
	}
	return nil
}
