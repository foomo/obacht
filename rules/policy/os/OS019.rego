package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.legacy_kexts_blocked
	f := {"rule_id": "OS019", "evidence": "Legacy kernel extensions (kexts) are not blocked"}
}
