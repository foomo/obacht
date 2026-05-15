package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.stealth_mode_enabled
	f := {"rule_id": "OS004", "evidence": "Stealth Mode is disabled"}
}
