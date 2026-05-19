package obacht.ssh

import rego.v1

findings contains f if {
	some key in input.keys
	key.algorithm == "DSA"
	f := {
		"rule_id": "SSH005",
		"evidence": sprintf("SSH key %s uses weak algorithm: DSA/%d bits", [key.path, key.bits]),
	}
}

findings contains f if {
	some key in input.keys
	key.algorithm == "RSA"
	key.bits > 0
	key.bits < 3072
	f := {
		"rule_id": "SSH005",
		"evidence": sprintf("SSH key %s uses weak algorithm: RSA/%d bits (minimum 3072)", [key.path, key.bits]),
	}
}

skips contains s if {
	some key in input.keys
	key.bits == 0
	s := {
		"rule_id": "SSH005",
		"evidence": sprintf("Could not read SSH key bits: %s (ssh-keygen failed)", [key.path]),
	}
}
