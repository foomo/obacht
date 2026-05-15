package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.rosetta_installed
	f := {"rule_id": "OS021", "evidence": "Rosetta 2 is installed; remove if no longer needed to reduce attack surface"}
}
