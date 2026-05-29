package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_telemetry != "1"
	f := {
		"rule_id": "CLD007",
		"evidence": sprintf("env.DISABLE_TELEMETRY is %s, expected \"1\"", [input.env_disable_telemetry]),
	}
}
