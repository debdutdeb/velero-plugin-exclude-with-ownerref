package config

import (
	"fmt"
	"regexp"
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

func (rw *RegexWrapper) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var regexStr string

	if err := unmarshal(&regexStr); err != nil {
		return err
	}

	compiled, err := regexp.Compile(regexStr)
	if err != nil {
		return fmt.Errorf("failed to compile regex %q: %w", regexStr, err)
	}

	rw.Regexp = compiled
	return nil
}
