package obacht.git

import rego.v1

findings contains f if {
	input.installed
	count(input.gitignore_missing_env) > 0
	f := {
		"rule_id": "GIT006",
		"evidence": sprintf("Global gitignore is missing %d env entries: %s", [count(input.gitignore_missing_env), concat(", ", input.gitignore_missing_env)]),
	}
}
