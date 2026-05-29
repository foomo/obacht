package obacht.claude

import rego.v1

findings contains f if {
	input.settings_disable_deep_link_registration != "disable"
	f := {
		"rule_id": "CLD025",
		"evidence": sprintf("disableDeepLinkRegistration is %s, expected \"disable\"", [input.settings_disable_deep_link_registration]),
	}
}
