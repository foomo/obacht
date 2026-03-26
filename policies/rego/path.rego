package bouncer.path

import rego.v1

findings contains f if {
	dir := input.path.dirs[_]
	dir.exists
	dir.writable
	f := {
		"rule_id": "PTH001",
		"evidence": sprintf("Writable directory in PATH: %s", [dir.path]),
	}
}

findings contains f if {
	dir := input.path.dirs[_]
	dir.is_relative
	f := {
		"rule_id": "PTH002",
		"evidence": sprintf("Relative path in PATH: %s", [dir.path]),
	}
}
