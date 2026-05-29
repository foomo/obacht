package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_extra_usage_command != "1"
	f := {
		"rule_id": "CLD015",
		"evidence": sprintf("env.DISABLE_EXTRA_USAGE_COMMAND is %s, expected \"1\"", [input.env_disable_extra_usage_command]),
	}
}
