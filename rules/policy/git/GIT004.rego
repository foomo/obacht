package obacht.git

import rego.v1

findings contains f if {
	input.installed
	count(input.gitignore_missing_macos) > 0
	f := {
		"rule_id": "GIT004",
		"evidence": sprintf("Global gitignore is missing %d macOS entries: %s", [count(input.gitignore_missing_macos), concat(", ", input.gitignore_missing_macos)]),
	}
}
