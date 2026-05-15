package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.internet_sharing_disabled
	f := {"rule_id": "OS014", "evidence": "Internet Sharing is enabled"}
}
