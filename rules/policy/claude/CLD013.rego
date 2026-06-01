package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_upgrade_command != "1"
	f := {
		"rule_id": "CLD013",
		"evidence": sprintf("env.DISABLE_UPGRADE_COMMAND is %s, expected \"1\"", [input.env_disable_upgrade_command]),
	}
}
