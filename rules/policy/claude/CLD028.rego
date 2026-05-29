package obacht.claude

import rego.v1

findings contains f if {
	input.settings_skip_web_fetch_preflight != "true"
	f := {
		"rule_id": "CLD028",
		"evidence": sprintf("skipWebFetchPreflight is %s, expected true", [input.settings_skip_web_fetch_preflight]),
	}
}
