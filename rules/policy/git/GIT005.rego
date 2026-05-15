package obacht.git

import rego.v1

findings contains f if {
	input.installed
	count(input.gitignore_missing_keys) > 0
	f := {
		"rule_id": "GIT005",
		"evidence": sprintf("Global gitignore is missing %d key/cert entries: %s", [count(input.gitignore_missing_keys), concat(", ", input.gitignore_missing_keys)]),
	}
}
