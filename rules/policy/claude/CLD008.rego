package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_bug_command != "1"
	f := {
		"rule_id": "CLD008",
		"evidence": sprintf("env.DISABLE_BUG_COMMAND is %s, expected \"1\"", [input.env_disable_bug_command]),
	}
}
