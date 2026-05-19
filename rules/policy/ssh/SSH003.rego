package obacht.ssh

import rego.v1

findings contains f if {
	input.config_exists
	input.strict_host_key_checking_disabled
	f := {
		"rule_id": "SSH003",
		"evidence": "SSH config has StrictHostKeyChecking set to 'no'",
	}
}
