package bouncer.docker

import rego.v1

findings contains f if {
	input.docker.socket_exists
	mode := input.docker.socket_mode
	mode != "0660"
	mode != "0600"
	mode != "0700"
	f := {
		"rule_id": "DOC001",
		"evidence": sprintf("Docker socket has mode %s (expected 0660 or stricter)", [mode]),
	}
}

findings contains f if {
	input.docker.user_in_group
	f := {
		"rule_id": "DOC002",
		"evidence": "Current user is a member of the docker group",
	}
}
