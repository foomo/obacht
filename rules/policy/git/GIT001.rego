package obacht.git

import rego.v1

findings contains f if {
	input.credential_helper == "store"
	f := {
		"rule_id": "GIT001",
		"evidence": "Git credential helper is set to 'store' which saves passwords in plaintext",
	}
}
