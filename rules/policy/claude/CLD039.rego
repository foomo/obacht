package obacht.claude

import rego.v1

findings contains f if {
	input.permissions_deny_home_secrets_missing != ""
	f := {
		"rule_id": "CLD039",
		"evidence": sprintf("permissions.deny missing home credential entries:\n        • %s", [concat("\n        • ", split(input.permissions_deny_home_secrets_missing, " "))]),
	}
}
