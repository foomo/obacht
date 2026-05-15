package obacht.docker

import rego.v1

findings contains f if {
	input.user_in_group
	f := {
		"rule_id": "DOC002",
		"evidence": "Current user is a member of the docker group",
	}
}
