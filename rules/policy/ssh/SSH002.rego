package obacht.ssh

import rego.v1

findings contains f if {
	input.directory_exists
	input.directory_mode != "0700"
	f := {
		"rule_id": "SSH002",
		"evidence": sprintf("~/.ssh has mode %s (expected 0700)", [input.directory_mode]),
	}
}
