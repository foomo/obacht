package obacht.claude

import rego.v1

findings contains f if {
	input.env_claude_code_disable_feedback_survey != "1"
	f := {
		"rule_id": "CLD019",
		"evidence": sprintf("env.CLAUDE_CODE_DISABLE_FEEDBACK_SURVEY is %s, expected \"1\"", [input.env_claude_code_disable_feedback_survey]),
	}
}
