package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.password_required_on_lock
	input.password_lock_delay_seconds > 0
	f := {
		"rule_id": "OS031",
		"evidence": sprintf("Password not required immediately after screen lock (delay: %ds)", [input.password_lock_delay_seconds]),
	}
}

findings contains f if {
	input.os == "darwin"
	input.password_lock_delay_seconds == -2
	f := {
		"rule_id": "OS031",
		"evidence": "Screen lock is disabled — no password required after screensaver",
	}
}

skips contains s if {
	input.os == "darwin"
	input.password_lock_delay_seconds == -1
	s := {
		"rule_id": "OS031",
		"evidence": "Could not determine screen-lock password delay (sysadminctl unavailable and legacy askForPasswordDelay unreadable)",
	}
}
