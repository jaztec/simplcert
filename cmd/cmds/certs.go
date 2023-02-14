package cmds

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jaztec/simplcert"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func createCertCmd() *cli.Command {
	return &cli.Command{
		Name:     "create",
		HelpName: "create a named certificate",
		Action: func(c *cli.Context) error {
			checkVerboseFlag(c)
			certPath, err := promptRootCertPath(c)
			if err != nil {
				return err
			}

			m, err := simplcert.NewManager(certPath)
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
			crtBytes := simplcert.EncodeCertificate(crt.Raw)
			privBytes, err := m.MarshalPrivateKey(priv)
			if err != nil {
				return err
			}

			if cfg.OutputPath != "" {
				name := cfg.OutputName
				if name == "" {
					name = strcase.ToKebab(cfg.Name)
				}
				err = outputToFile(cfg.OutputPath, name, m.RootPEM(), crtBytes, privBytes, pub)
			} else {
				outputToScreen(crtBytes, privBytes, pub)
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
			ecdsaFlag(),
			rsaFlag(),
			ed25519Flag(),
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

			certPath, err := promptRootCertPath(c)
			if err != nil {
				return err
			}

			log.WithField("path", certPath).Info("Verify certificates")
			if fi, err := os.Stat(certPath); err != nil || !fi.IsDir() {
				return fmt.Errorf("invalid path parameter (%s)", certPath)
			}

			if _, err = simplcert.NewManager(certPath); errors.Is(err, simplcert.NoCertsError) {
				return simplcert.CreateRootCAFiles(getCertType(c), certPath)
			}

			return nil
		},
		Flags: flags(
			certsPathFlag(),
			ecdsaFlag(),
			rsaFlag(),
			ed25519Flag(),
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

			certPath, err := promptRootCertPath(c)
			if err != nil {
				return err
			}

			if fi, err := os.Stat(certPath); err != nil || !fi.IsDir() {
				return fmt.Errorf("invalid path parameter (%s)", certPath)
			}

			m, err := simplcert.NewManager(certPath)
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

func printConfig(cfg simplcert.CertConfig) {
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

func getCertType(c *cli.Context) simplcert.CertType {
	if c.Bool("ecdsa") {
		return simplcert.TypeECDSA
	}

	if c.Bool("rsa") {
		return simplcert.TypeRSA
	}

	if c.Bool("ed25519") {
		return simplcert.TypeED25519
	}

	return simplcert.TypeECDSA
}
