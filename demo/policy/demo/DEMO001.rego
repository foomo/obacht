package obacht.demo

import rego.v1

findings contains f if {
	input.world_writable
	f := {
		"rule_id": "DEMO001",
		"evidence": sprintf("%s has mode %s (world-writable)", [input.path, input.mode]),
	}
}
