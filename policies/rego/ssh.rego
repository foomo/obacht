package bouncer.ssh

import rego.v1

findings contains f if {
	key := input.ssh.keys[_]
	key.mode != "0600"
	f := {
		"rule_id": "SSH001",
		"evidence": sprintf("Key %s has mode %s (expected 0600)", [key.path, key.mode]),
	}
}

findings contains f if {
	input.ssh.directory_exists
	input.ssh.directory_mode != "0700"
	f := {
		"rule_id": "SSH002",
		"evidence": sprintf("~/.ssh has mode %s (expected 0700)", [input.ssh.directory_mode]),
	}
}
