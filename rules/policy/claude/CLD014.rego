package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_feedback_command != "1"
	f := {
		"rule_id": "CLD014",
		"evidence": sprintf("env.DISABLE_FEEDBACK_COMMAND is %s, expected \"1\"", [input.env_disable_feedback_command]),
	}
}
