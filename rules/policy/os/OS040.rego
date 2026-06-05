package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.auto_submit
	f := {"rule_id": "OS040", "evidence": "Diagnostic data is auto-submitted to Apple (AutoSubmit=1)"}
}

findings contains f if {
	input.os == "darwin"
	input.third_party_submit
	f := {"rule_id": "OS040", "evidence": "Third-party diagnostic data submission is enabled"}
}
