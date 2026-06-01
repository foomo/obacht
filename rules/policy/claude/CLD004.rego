package obacht.claude

import rego.v1

findings contains f if {
	input.claude_in_chrome_default_enabled != "false"
	f := {
		"rule_id": "CLD004",
		"evidence": sprintf("claudeInChromeDefaultEnabled is %s, expected false", [input.claude_in_chrome_default_enabled]),
	}
}
