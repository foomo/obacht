package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.remote_login_disabled
	f := {"rule_id": "OS024", "evidence": "Remote Login (SSH server) is enabled"}
}
