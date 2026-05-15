package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.os_auto_update_enabled
	f := {"rule_id": "OS009", "evidence": "Automatic OS updates are disabled"}
}
