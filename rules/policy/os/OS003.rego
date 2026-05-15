package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.firewall_enabled
	f := {"rule_id": "OS003", "evidence": "Application Firewall is disabled"}
}
