package obacht.demo

import rego.v1

findings contains f if {
	some key in input.keys
	key.mode != "0600"
	f := {
		"rule_id": "DEMO003",
		"evidence": sprintf("Key %s has mode %s (expected 0600)", [key.path, key.mode]),
	}
}
