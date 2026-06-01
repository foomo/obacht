package obacht.claude

import rego.v1

findings contains f if {
	input.sandbox_allow_unsandboxed_commands != "false"
	f := {
		"rule_id": "CLD033",
		"evidence": sprintf("sandbox.allowUnsandboxedCommands is %s, expected false", [input.sandbox_allow_unsandboxed_commands]),
	}
}
