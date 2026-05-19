package obacht.tools

import rego.v1

findings contains f if {
	not input.analytics_disabled
	f := {
		"rule_id": "TOL004",
		"evidence": "HOMEBREW_NO_ANALYTICS is not set to \"1\"",
	}
}
