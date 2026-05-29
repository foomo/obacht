package obacht.claude

import rego.v1

findings contains f if {
	input.settings_auto_memory_directory != ".claude/memory"
	f := {
		"rule_id": "CLD029",
		"evidence": sprintf("autoMemoryDirectory is %s, expected \".claude/memory\"", [input.settings_auto_memory_directory]),
	}
}
