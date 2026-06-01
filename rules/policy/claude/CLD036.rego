package obacht.claude

import rego.v1

findings contains f if {
	input.permissions_deny_network_missing != ""
	f := {
		"rule_id": "CLD036",
		"evidence": sprintf("permissions.deny missing network/exfiltration entries:\n        • %s", [concat("\n        • ", split(input.permissions_deny_network_missing, " "))]),
	}
}
