package obacht.ssh

import rego.v1

findings contains f if {
	key := input.keys[_]
	key.mode != "0600"
	f := {
		"rule_id": "SSH001",
		"evidence": sprintf("Key %s has mode %s (expected 0600)", [key.path, key.mode]),
	}
}
