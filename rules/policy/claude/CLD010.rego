package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_login_command != "1"
	f := {
		"rule_id": "CLD010",
		"evidence": sprintf("env.DISABLE_LOGIN_COMMAND is %s, expected \"1\"", [input.env_disable_login_command]),
	}
}
