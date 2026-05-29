package obacht.claude

import rego.v1

findings contains f if {
	not input.gitignore_excludes_settings
	f := {
		"rule_id": "CLD001",
		"evidence": "Global gitignore does not exclude **/.claude/settings.local.json",
	}
}
