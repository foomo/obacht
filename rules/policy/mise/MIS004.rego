package obacht.mise

import rego.v1

findings contains f if {
	input.value != "gh auth token"
	f := {
		"rule_id": "MIS004",
		"evidence": sprintf("mise github.credential_command is %q (recommend \"gh auth token\" to reuse gh auth)", [input.value]),
	}
}
