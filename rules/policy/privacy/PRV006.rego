package obacht.privacy

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.user_trust_overrides > 0
	f := {
		"rule_id": "PRV006",
		"evidence": sprintf("%d user-domain certificate trust override(s) found — verify each is legitimate", [input.user_trust_overrides]),
	}
}

findings contains f if {
	input.os == "darwin"
	input.admin_trust_overrides > 0
	f := {
		"rule_id": "PRV006",
		"evidence": sprintf("%d admin-domain certificate trust override(s) found — typically MDM-installed; verify each is legitimate", [input.admin_trust_overrides]),
	}
}
