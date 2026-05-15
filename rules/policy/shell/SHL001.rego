package obacht.shell

import rego.v1

findings contains f if {
	input.history_file_mode != ""
	input.history_file_mode != "0600"
	f := {
		"rule_id": "SHL001",
		"evidence": sprintf("History file %s has mode %s (expected 0600)", [input.history_file, input.history_file_mode]),
	}
}
