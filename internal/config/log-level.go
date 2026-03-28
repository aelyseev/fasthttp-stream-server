package config

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
)

var allowedLevels = []string{LevelDebug, LevelInfo, LevelWarn}

type LogLevel string

func (e LogLevel) SlogLevel() slog.Level {
	switch e {
	case LevelDebug:
		return slog.LevelDebug
	case LevelWarn:
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func (e *LogLevel) UnmarshalYAML(value *yaml.Node) error {
	var env string
	if err := value.Decode(&env); err != nil {
		return err
	}

	if slices.Contains(allowedLevels, env) {
		*e = LogLevel(env)
		return nil
	}

	return fmt.Errorf(
		"unsupported value '%s' for 'level' field, only [%s] are allowed ",
		env,
		strings.Join(allowedLevels, ", "),
	)
}
