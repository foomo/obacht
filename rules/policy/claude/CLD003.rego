package obacht.claude

import rego.v1

findings contains f if {
	input.pr_status_footer_enabled != "false"
	f := {
		"rule_id": "CLD003",
		"evidence": sprintf("prStatusFooterEnabled is %s, expected false", [input.pr_status_footer_enabled]),
	}
}
