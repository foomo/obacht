package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.touchid_configured
	f := {"rule_id": "OS042", "evidence": "Touch ID for sudo is not configured (pam_tid.so absent from /etc/pam.d/sudo and sudo_local)"}
}
