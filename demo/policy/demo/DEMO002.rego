package obacht.demo

import rego.v1

findings contains f if {
	some key in input.keys
	key.algorithm == "RSA"
	key.bits < 3072
	f := {
		"rule_id": "DEMO002",
		"evidence": sprintf("SSH key %s uses RSA/%d bits (minimum 3072)", [key.path, key.bits]),
	}
}
