package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.guest_account_disabled
	f := {"rule_id": "OS007", "evidence": "Guest account is enabled"}
}
