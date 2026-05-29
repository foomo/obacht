package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_auto_compact != "1"
	f := {
		"rule_id": "CLD009",
		"evidence": sprintf("env.DISABLE_AUTO_COMPACT is %s, expected \"1\"", [input.env_disable_auto_compact]),
	}
}
