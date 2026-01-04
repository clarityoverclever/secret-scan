package models

import "regexp"

type PatternDefinition struct {
	Name      string
	Regex     string
	Severity  string
	Validator string
}

type CompiledPattern struct {
	Name      string
	Severity  string
	Regex     *regexp.Regexp
	Validator string
}
