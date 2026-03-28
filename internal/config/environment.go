package config

import (
	"fmt"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

var allowedEnvs = []string{EnvLocal, EnvDev, EnvProd}

type Environment string

func (e *Environment) UnmarshalYAML(value *yaml.Node) error {
	var env string
	if err := value.Decode(&env); err != nil {
		return err
	}

	if slices.Contains(allowedEnvs, env) {
		*e = Environment(env)
		return nil
	}

	return fmt.Errorf(
		"unsupported value '%s' for 'level' field, only [%s] are allowed ",
		env,
		strings.Join(allowedEnvs, ", "),
	)
}
