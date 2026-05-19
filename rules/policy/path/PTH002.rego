package obacht.path

import rego.v1

findings contains f if {
	some dir in input.dirs
	dir.is_relative
	f := {
		"rule_id": "PTH002",
		"evidence": sprintf("Relative path in PATH: %s", [dir.path]),
	}
}
