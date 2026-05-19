package obacht.os

import rego.v1

tm_dest_evaluable if {
	input.os == "darwin"
	input.timemachine_enabled
	input.timemachine_destination_connected
}

skips contains s if {
	input.os == "darwin"
	input.timemachine_enabled
	not input.timemachine_destination_connected
	s := {"rule_id": "OS033", "evidence": "Time Machine destination not connected — cannot evaluate backup recency"}
}

findings contains f if {
	tm_dest_evaluable
	not input.timemachine_recent_backup
	f := {"rule_id": "OS033", "evidence": "Time Machine has no backup within the last 14 days"}
}
