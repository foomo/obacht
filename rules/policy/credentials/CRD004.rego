package obacht.credentials

import rego.v1

findings contains f if {
	input.npmrc_has_token
	input.npmrc_mode != "0600"
	f := {
		"rule_id": "CRD004",
		"evidence": sprintf("~/.npmrc contains auth token and has mode %s (expected 0600)", [input.npmrc_mode]),
	}
}
