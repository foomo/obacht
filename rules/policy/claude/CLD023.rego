package obacht.claude

import rego.v1

findings contains f if {
	input.env_is_demo != "1"
	f := {
		"rule_id": "CLD023",
		"evidence": sprintf("env.IS_DEMO is %s, expected \"1\"", [input.env_is_demo]),
	}
}
