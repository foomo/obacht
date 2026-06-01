package obacht.claude

import rego.v1

findings contains f if {
	input.permissions_deny_git_missing != ""
	f := {
		"rule_id": "CLD038",
		"evidence": sprintf("permissions.deny missing destructive git entries:\n        • %s", [concat("\n        • ", split(input.permissions_deny_git_missing, " "))]),
	}
}
