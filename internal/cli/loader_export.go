package cli

import "github.com/foomo/obacht/pkg/schema"

// LoadEmbeddedRuleFiles is the exported wrapper for loading the embedded
// rule set. Exposed for use by integration tests and tooling outside this
// package.
func LoadEmbeddedRuleFiles() ([]schema.RulesFile, error) {
	return loadEmbeddedRuleFiles()
}
