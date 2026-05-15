package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.rsr_enabled
	f := {"rule_id": "OS011", "evidence": "Rapid Security Responses are disabled"}
}
