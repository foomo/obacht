package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.in_app_review_disabled
	f := {"rule_id": "OS037", "evidence": "App Store in-app review prompts are not disabled"}
}
