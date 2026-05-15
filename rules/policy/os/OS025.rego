package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.remote_management_disabled
	f := {"rule_id": "OS025", "evidence": "Remote Management is enabled"}
}
