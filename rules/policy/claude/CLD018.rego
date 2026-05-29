package obacht.claude

import rego.v1

findings contains f if {
	input.env_claude_code_disable_cron != "1"
	f := {
		"rule_id": "CLD018",
		"evidence": sprintf("env.CLAUDE_CODE_DISABLE_CRON is %s, expected \"1\"", [input.env_claude_code_disable_cron]),
	}
}
