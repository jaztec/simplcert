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

type promptFlag[T string | bool] struct {
	name           string
	promptDatatype promptDatatype
	promptQuestion promptQuestion
	setter         func(config *security.CertConfig, val T)
	required       bool
}

func certConfigFromFlags(c *cli.Context) (security.CertConfig, error) {
	cfg := security.CertConfig{}

	stringFlags := []promptFlag[string]{
		{
			name:           "name",
			promptQuestion: newStringPromptQuestion("Name"),
			required:       true,
			setter: func(cfg *security.CertConfig, val string) {
				cfg.Name = val
			},
		},
		{
			name:           "host",
			promptQuestion: newStringPromptQuestion("Host"),
			required:       true,
			setter: func(cfg *security.CertConfig, val string) {
				cfg.Host = val
			},
		},
		{
			name:           "country",
			promptQuestion: newStringPromptQuestion("Country"),
			setter: func(cfg *security.CertConfig, val string) {
				cfg.Country = val
			},
		},
		{
			name:           "organization",
			promptQuestion: newStringPromptQuestion("Organization"),
			setter: func(cfg *security.CertConfig, val string) {
				cfg.Organization = val
			},
		},
	}
	boolFlags := []promptFlag[bool]{
		{
			name:           "is_server",
			promptQuestion: newBoolPromptQuestion("Target is a server"),
			setter: func(cfg *security.CertConfig, val bool) {
				cfg.IsServer = val
			},
		},
	}

	for _, flag := range stringFlags {
		v := c.String(flag.name)
		if v == "" {
			r, err := prompt[string](flag.promptQuestion, stringParser)
			if err != nil {
				return cfg, err
			}
			v = r.value
		}
		if flag.required && v == "" {
			return cfg, fmt.Errorf("flag \"%s\" is required nu no value provided", flag.name)
		}
		flag.setter(&cfg, v)
	}

	for _, flag := range boolFlags {
		v := c.Bool(flag.name)
		if v == false {
			r, err := prompt[bool](flag.promptQuestion, boolParser)
			if err != nil {
				return cfg, err
			}
			v = r.value
		}
		flag.setter(&cfg, v)
	}

	return cfg, nil
}

func prompt[T string | bool](q promptQuestion, parser func(val string) T) (promptResult[T], error) {
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
	set := parser(intermediate)

	log.WithFields(log.Fields{
		"question":   q.question,
		"registered": intermediate,
		"parser":     set,
	}).Debug("Setting flag value")

	return promptResult[T]{set}, nil
}

func stringParser(val string) string { return val }
func boolParser(val string) bool {
	switch val {
	case "Y":
		fallthrough
	case "y":
		return true
	default:
		return false
	}
}
