package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_logout_command != "1"
	f := {
		"rule_id": "CLD011",
		"evidence": sprintf("env.DISABLE_LOGOUT_COMMAND is %s, expected \"1\"", [input.env_disable_logout_command]),
	}
}
