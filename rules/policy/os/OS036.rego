package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.macos_major_version > 0
	input.macos_major_version < 13
	f := {
		"rule_id": "OS036",
		"evidence": sprintf("macOS major version %d is no longer receiving security updates (minimum supported: 13)", [input.macos_major_version]),
	}
}
