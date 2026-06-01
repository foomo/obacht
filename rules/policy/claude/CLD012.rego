package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_error_reporting != "1"
	f := {
		"rule_id": "CLD012",
		"evidence": sprintf("env.DISABLE_ERROR_REPORTING is %s, expected \"1\"", [input.env_disable_error_reporting]),
	}
}
