package obacht.kube

import rego.v1

findings contains f if {
	input.config_exists
	input.config_mode != "0600"
	f := {
		"rule_id": "KUB001",
		"evidence": sprintf("~/.kube/config has mode %s (expected 0600)", [input.config_mode]),
	}
}
