package obacht.claude

import rego.v1

findings contains f if {
	input.settings_respect_gitignore != "true"
	f := {
		"rule_id": "CLD027",
		"evidence": sprintf("respectGitignore is %s, expected true", [input.settings_respect_gitignore]),
	}
}
