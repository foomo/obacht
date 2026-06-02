package obacht.brew

import rego.v1

findings contains f if {
	input.homebrew_installed
	input.homebrew_auto_update_disabled
	f := {
		"rule_id": "BRW001",
		"evidence": "HOMEBREW_NO_AUTO_UPDATE is set, preventing automatic formula updates",
	}
}
