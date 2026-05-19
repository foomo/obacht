package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.media_sharing_disabled
	f := {"rule_id": "OS027", "evidence": "Media Sharing is enabled"}
}
