package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.screen_sharing_disabled
	f := {"rule_id": "OS013", "evidence": "Screen Sharing is enabled"}
}
