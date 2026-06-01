package obacht.claude

import rego.v1

findings contains f if {
	input.permissions_deny_destructive_fs_missing != ""
	f := {
		"rule_id": "CLD037",
		"evidence": sprintf("permissions.deny missing destructive filesystem entries:\n        • %s", [concat("\n        • ", split(input.permissions_deny_destructive_fs_missing, " "))]),
	}
}
