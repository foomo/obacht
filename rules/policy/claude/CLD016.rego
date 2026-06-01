package obacht.claude

import rego.v1

findings contains f if {
	input.env_claude_code_disable_fast_mode != "1"
	f := {
		"rule_id": "CLD016",
		"evidence": sprintf("env.CLAUDE_CODE_DISABLE_FAST_MODE is %s, expected \"1\"", [input.env_claude_code_disable_fast_mode]),
	}
}
