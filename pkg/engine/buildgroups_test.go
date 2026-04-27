package engine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/foomo/obacht/pkg/engine"
	"github.com/foomo/obacht/pkg/schema"
)

// TestBuildRuleGroups_MergesByInput verifies that rules in a RulesFile sharing
// the same effective input are merged into a single group, even when each
// rule defines its own inline policy. Without this, every rule would trigger
// a separate input run + rego compile/eval, which is the bug that made `obacht
// scan` take ~15s instead of <2s.
func TestBuildRuleGroups_MergesByInput(t *testing.T) {
	rf := schema.RulesFile{
		Input: `printf '{"x":1}'`,
		Rules: []schema.Rule{
			{ID: "A1", Category: "a", Policy: `findings contains f if { f := {"rule_id":"A1","evidence":"a1"} }`},
			{ID: "A2", Category: "a", Policy: `findings contains f if { f := {"rule_id":"A2","evidence":"a2"} }`},
			{ID: "A3", Category: "a", Policy: `findings contains f if { f := {"rule_id":"A3","evidence":"a3"} }`},
		},
	}

	groups := engine.BuildRuleGroups([]schema.RulesFile{rf})
	require.Len(t, groups, 1, "rules sharing the same input must merge into one group")

	g := groups[0]
	assert.Len(t, g.Rules, 3)
	assert.Contains(t, g.Policy, `"rule_id":"A1"`)
	assert.Contains(t, g.Policy, `"rule_id":"A2"`)
	assert.Contains(t, g.Policy, `"rule_id":"A3"`)
	assert.Contains(t, g.Policy, "package obacht.a")
}

// TestBuildRuleGroups_IsolatesPackageDecl verifies that a rule whose policy
// declares its own `package` is given its own group (concatenating multiple
// package declarations would produce an invalid rego module).
func TestBuildRuleGroups_IsolatesPackageDecl(t *testing.T) {
	rf := schema.RulesFile{
		Input: `printf '{"x":1}'`,
		Rules: []schema.Rule{
			{ID: "A1", Category: "a", Policy: `findings contains f if { f := {"rule_id":"A1","evidence":"a1"} }`},
			{ID: "A2", Category: "a", Policy: "package obacht.custom\nimport rego.v1\n\nfindings contains f if { f := {\"rule_id\":\"A2\",\"evidence\":\"a2\"} }"},
		},
	}

	groups := engine.BuildRuleGroups([]schema.RulesFile{rf})
	require.Len(t, groups, 2)

	ids := []string{groups[0].Rules[0].ID, groups[1].Rules[0].ID}
	assert.ElementsMatch(t, []string{"A1", "A2"}, ids)
	assert.Len(t, groups[0].Rules, 1)
	assert.Len(t, groups[1].Rules, 1)
}

// TestBuildRuleGroups_RuleLevelInputSeparates verifies that a rule with a
// distinct rule-level input is placed in its own group.
func TestBuildRuleGroups_RuleLevelInputSeparates(t *testing.T) {
	rf := schema.RulesFile{
		Input:  `printf '{"x":1}'`,
		Policy: `findings contains f if { f := {"rule_id":"A1","evidence":"a1"} }`,
		Rules: []schema.Rule{
			{ID: "A1", Category: "a"},
			{ID: "A2", Category: "a", Input: `printf '{"x":2}'`},
		},
	}

	groups := engine.BuildRuleGroups([]schema.RulesFile{rf})
	require.Len(t, groups, 2)
}
