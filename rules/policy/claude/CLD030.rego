package obacht.claude

import rego.v1

findings contains f if {
	input.settings_plans_directory != ".claude/plans"
	f := {
		"rule_id": "CLD030",
		"evidence": sprintf("plansDirectory is %s, expected \".claude/plans\"", [input.settings_plans_directory]),
	}
}
