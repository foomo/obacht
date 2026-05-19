package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.auto_login_disabled
	f := {"rule_id": "OS006", "evidence": "Automatic login is enabled"}
}
