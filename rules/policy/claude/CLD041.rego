package obacht.claude

import rego.v1

findings contains f if {
	count(input.claude_desktop_native_messaging_manifests) > 0
	f := {
		"rule_id": "CLD041",
		"evidence": sprintf("Unmitigated Claude Desktop Native Messaging manifest(s):\n        • %s", [concat("\n        • ", input.claude_desktop_native_messaging_manifests)]),
	}
}
