package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.timemachine_enabled
	f := {"rule_id": "OS023", "evidence": "Time Machine backup is not enabled or no backup destination is configured"}
}
