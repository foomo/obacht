package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.content_caching_disabled
	f := {"rule_id": "OS029", "evidence": "Content Caching is enabled"}
}
