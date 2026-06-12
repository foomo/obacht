package obacht.yarn

import rego.v1

findings contains f if {
	some t in input.tokens
	f := {
		"rule_id": "YRN003",
		"evidence": sprintf("plaintext %s in %s:%d (value: %s)", [t.key, input.config_path, t.line, t.masked]),
	}
}
