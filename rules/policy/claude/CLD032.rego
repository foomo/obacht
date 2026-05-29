package obacht.claude

import rego.v1

findings contains f if {
	input.sandbox_auto_allow_bash_if_sandboxed != "false"
	f := {
		"rule_id": "CLD032",
		"evidence": sprintf("sandbox.autoAllowBashIfSandboxed is %s, expected false", [input.sandbox_auto_allow_bash_if_sandboxed]),
	}
}
