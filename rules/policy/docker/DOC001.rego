package obacht.docker

import rego.v1

findings contains f if {
	input.socket_exists
	mode := input.socket_mode
	mode != "0660"
	mode != "0600"
	mode != "0700"
	mode != "0755"
	f := {
		"rule_id": "DOC001",
		"evidence": sprintf("Docker socket has mode %s (expected 0660 or stricter)", [mode]),
	}
}
