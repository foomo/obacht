package obacht.git

import rego.v1

findings contains f if {
	input.installed
	input.safe_directory_wildcard
	f := {
		"rule_id": "GIT003",
		"evidence": "Git safe.directory is set to '*', disabling ownership checks for all repositories",
	}
}
