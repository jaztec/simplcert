package cmds

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gitlab.jaztec.info/certs/manager/pkg/security"
	"os"
	"strings"
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

type promptDatatype int

const (
	stringPrompt promptDatatype = iota
	boolPrompt
)

type promptQuestion struct {
	question     string
	datatype     promptDatatype
	defaultValue string
}

func newStringPromptQuestion(name string) promptQuestion {
	return promptQuestion{
		question:     name,
		datatype:     stringPrompt,
		defaultValue: "",
	}
}

func newBoolPromptQuestion(name string) promptQuestion {
	return promptQuestion{
		question:     name,
		datatype:     boolPrompt,
		defaultValue: "",
	}
}

type promptResult[T string | bool] struct {
	value T
}

func certConfigFromFlags(c *cli.Context) (security.CertConfig, error) {
	cfg := security.CertConfig{}

	appFlags := []struct {
		name           string
		promptDatatype promptDatatype
		promptQuestion promptQuestion
		target         any
	}{
		{
			name:           "name",
			promptQuestion: newStringPromptQuestion("Name"),
			target:         &cfg.Name,
		},
		{
			name:           "host",
			promptQuestion: newStringPromptQuestion("Host"),
			target:         &cfg.Host,
		},
		{
			name:           "country",
			promptQuestion: newStringPromptQuestion("Country"),
			target:         &cfg.Country,
		},
		{
			name:           "organization",
			promptQuestion: newStringPromptQuestion("Organization"),
			target:         &cfg.Organization,
		},
		{
			name:           "is_server",
			promptQuestion: newBoolPromptQuestion("Target is a server"),
			target:         &cfg.IsServer,
		},
	}

	for _, flag := range appFlags {
		switch flag.promptDatatype {
		case boolPrompt:
			v := c.Bool(flag.name)
			if v == false {
				r, err := prompt[bool](flag.promptQuestion, func(val string) bool {
					switch val {
					case "Y":
						fallthrough
					case "y":
						return true
					default:
						return false
					}
				})
				if err != nil {
					return cfg, err
				}
				v = r.value
			}
			flag.target = &v
			log.WithFields(log.Fields{
				"v":           fmt.Sprintf("%+v", v),
				"flag.target": flag.target,
			})
		case stringPrompt:
			v := c.String(flag.name)
			if v == "" {
				r, err := prompt[string](flag.promptQuestion, func(val string) string { return val })
				if err != nil {
					return cfg, err
				}
				v = r.value
			}
			flag.target = &v
		}
	}

	return cfg, nil
}

func prompt[T string | bool](q promptQuestion, setter func(val string) T) (promptResult[T], error) {
	var base T

	switch q.datatype {
	case boolPrompt:
		fmt.Printf("%s [y/N]: ", q.question)
	case stringPrompt:
		fmt.Printf("%s: ", q.question)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return promptResult[T]{base}, err
	}
	intermediate := strings.TrimSpace(input)
	set := setter(intermediate)

	log.WithFields(log.Fields{
		"question":   q.question,
		"registered": intermediate,
		"setter":     set,
	}).Debug("Setting flag value")

	return promptResult[T]{set}, nil
}
