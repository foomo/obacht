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
	s := {"rule_id": "OS032", "evidence": "Time Machine destination not connected — cannot evaluate encryption"}
}

findings contains f if {
	tm_dest_evaluable
	not input.timemachine_destination_encrypted
	f := {"rule_id": "OS032", "evidence": "Time Machine destination is not encrypted"}
}
