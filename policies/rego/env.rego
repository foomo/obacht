package bouncer.env

import rego.v1

findings contains f if {
	var := input.env.suspicious_vars[_]
	f := {
		"rule_id": "ENV001",
		"evidence": sprintf("Suspicious env var: %s (matched pattern: %s)", [var.name, var.pattern]),
	}
}
