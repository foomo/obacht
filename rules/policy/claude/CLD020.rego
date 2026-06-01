package obacht.claude

import rego.v1

findings contains f if {
	input.env_claude_code_disable_file_checkpointing != "1"
	f := {
		"rule_id": "CLD020",
		"evidence": sprintf("env.CLAUDE_CODE_DISABLE_FILE_CHECKPOINTING is %s, expected \"1\"", [input.env_claude_code_disable_file_checkpointing]),
	}
}
