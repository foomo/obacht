package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.filevault_enabled
	f := {"rule_id": "OS002", "evidence": "FileVault disk encryption is disabled"}
}
