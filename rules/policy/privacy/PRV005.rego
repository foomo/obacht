package obacht.privacy

import rego.v1

findings contains f if {
	not input.do_not_track_enabled
	f := {
		"rule_id": "PRV005",
		"evidence": "DO_NOT_TRACK env var not set to \"1\"",
	}
}
