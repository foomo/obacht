package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.app_auto_update_enabled
	f := {"rule_id": "OS010", "evidence": "Automatic App Store updates are disabled"}
}
