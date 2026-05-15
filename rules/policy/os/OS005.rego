package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.gatekeeper_enabled
	f := {"rule_id": "OS005", "evidence": "Gatekeeper is disabled"}
}
