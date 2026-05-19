package obacht.credentials

import rego.v1

findings contains f if {
	input.netrc_exists
	input.netrc_mode != "0600"
	f := {
		"rule_id": "CRD002",
		"evidence": sprintf("~/.netrc has mode %s (expected 0600)", [input.netrc_mode]),
	}
}
