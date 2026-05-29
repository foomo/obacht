package obacht.claude

import rego.v1

findings contains f if {
	input.settings_attribution_commit != ""
	f := {
		"rule_id": "CLD026",
		"evidence": sprintf("attribution.commit is %s, expected \"\"", [input.settings_attribution_commit]),
	}
}

findings contains f if {
	input.settings_attribution_pr != ""
	f := {
		"rule_id": "CLD026",
		"evidence": sprintf("attribution.pr is %s, expected \"\"", [input.settings_attribution_pr]),
	}
}
