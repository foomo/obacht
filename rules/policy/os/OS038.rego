package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.allow_signed_builtin
	f := {"rule_id": "OS038", "evidence": "Application firewall auto-allows built-in signed software"}
}

findings contains f if {
	input.os == "darwin"
	input.allow_signed_downloaded
	f := {"rule_id": "OS038", "evidence": "Application firewall auto-allows downloaded signed software"}
}
