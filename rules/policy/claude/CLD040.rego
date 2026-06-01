package obacht.claude

import rego.v1

findings contains f if {
	input.permissions_deny_project_secrets_missing != ""
	f := {
		"rule_id": "CLD040",
		"evidence": sprintf("permissions.deny missing project secret entries:\n        • %s", [concat("\n        • ", split(input.permissions_deny_project_secrets_missing, " "))]),
	}
}
