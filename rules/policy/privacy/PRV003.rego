package obacht.privacy

import rego.v1

findings contains f if {
	not input.encrypted_dns
	f := {
		"rule_id": "PRV003",
		"evidence": "No encrypted DNS (DoH/DoT) configuration detected",
	}
}
