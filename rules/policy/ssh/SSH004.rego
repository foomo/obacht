package obacht.ssh

import rego.v1

findings contains f if {
	input.config_exists
	input.forward_agent_global
	f := {
		"rule_id": "SSH004",
		"evidence": "SSH agent forwarding is enabled globally (Host *)",
	}
}
