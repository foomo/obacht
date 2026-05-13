package obacht.path

import rego.v1

findings contains f if {
	dir := input.dirs[_]
	dir.is_relative
	f := {
		"rule_id": "PTH002",
		"evidence": sprintf("Relative path in PATH: %s", [dir.path]),
	}
}
