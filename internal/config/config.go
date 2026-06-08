package config

import (
	"fmt"
	"regexp"

	"go.yaml.in/yaml/v2"
)

type RegexWrapper struct {
	*regexp.Regexp
}

type ExcludeRegexForKind struct {
	Kind    string          `yaml:"kind"`
	Regexes []*RegexWrapper `yaml:"regexes"`
}

type Config struct {
	ExcludeResourcesWithOwnerReferences bool                  `yaml:"excludeResourcesWithOwnerReferences,omitempty"`
	ExcludeNamesRegexes                 []*RegexWrapper       `yaml:"excludeNamesRegexes,omitempty"`
	ExcludeRegexesForKinds              []ExcludeRegexForKind `yaml:"excludeRegexesForKinds,omitempty"`
}

func (rw *RegexWrapper) UnmarshalJSON(b []byte) error {
	var regexStr string

	if err := yaml.Unmarshal(b, &regexStr); err != nil {
		return err
	}

	compiled, err := regexp.Compile(regexStr)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %w", err)
	}

	rw.Regexp = compiled
	return nil
}
