package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.lockdown_enabled
	f := {"rule_id": "OS039", "evidence": "Lockdown Mode is not enabled"}
}
