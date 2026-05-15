package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.remote_apple_events_disabled
	f := {"rule_id": "OS016", "evidence": "Remote Apple Events are enabled"}
}
