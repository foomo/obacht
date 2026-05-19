package obacht.path

import rego.v1

findings contains f if {
	some dir in input.dirs
	dir.exists
	dir.world_writable
	f := {
		"rule_id": "PTH001",
		"evidence": sprintf("World-writable directory in PATH: %s (mode %s)", [dir.path, dir.mode]),
	}
}
