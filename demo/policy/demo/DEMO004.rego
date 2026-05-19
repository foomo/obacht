package obacht.demo

import rego.v1

findings contains f if {
	input.world_readable
	f := {
		"rule_id": "DEMO004",
		"evidence": sprintf("%s has mode %s (world-readable)", [input.path, input.mode]),
	}
}
