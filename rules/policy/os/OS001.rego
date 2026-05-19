package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.sip_enabled
	f := {"rule_id": "OS001", "evidence": "System Integrity Protection is disabled"}
}
