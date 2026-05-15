package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.user_is_admin
	f := {"rule_id": "OS030", "evidence": "Current user has local admin privileges"}
}
