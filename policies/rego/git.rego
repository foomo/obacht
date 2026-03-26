package bouncer.git

import rego.v1

findings contains f if {
	input.git.credential_helper == "store"
	f := {
		"rule_id": "GIT001",
		"evidence": "Git credential helper is set to 'store' which saves passwords in plaintext",
	}
}

findings contains f if {
	input.git.installed
	not input.git.signing_enabled
	f := {
		"rule_id": "GIT002",
		"evidence": "Git commit signing is not enabled",
	}
}
