package obacht.claude

import rego.v1

findings contains f if {
	input.env_claude_code_disable_experimental_betas != "1"
	f := {
		"rule_id": "CLD021",
		"evidence": sprintf("env.CLAUDE_CODE_DISABLE_EXPERIMENTAL_BETAS is %s, expected \"1\"", [input.env_claude_code_disable_experimental_betas]),
	}
}
