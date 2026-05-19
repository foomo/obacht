package obacht.credentials

import rego.v1

findings contains f if {
	input.aws_exists
	input.aws_mode != "0600"
	f := {
		"rule_id": "CRD001",
		"evidence": sprintf("~/.aws/credentials has mode %s (expected 0600)", [input.aws_mode]),
	}
}
