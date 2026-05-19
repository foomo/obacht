package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.mdm_enrolled
	f := {"rule_id": "OS020", "evidence": "Device is not enrolled in MDM"}
}
