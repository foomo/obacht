package engine

import "github.com/foomo/obacht/pkg/schema"

// BuildRuleGroups exposes buildRuleGroups for white-box testing.
type RuleGroup = ruleGroup

func BuildRuleGroups(ruleFiles []schema.RulesFile) []RuleGroup {
	return buildRuleGroups(ruleFiles)
}
