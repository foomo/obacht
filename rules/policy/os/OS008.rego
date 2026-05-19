package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.screen_lock_timeout_seconds > 300
	f := {
		"rule_id": "OS008",
		"evidence": sprintf("Screen lock timeout is %d seconds (maximum 300)", [input.screen_lock_timeout_seconds]),
	}
}
