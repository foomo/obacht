package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.show_all_extensions
	f := {"rule_id": "OS044", "evidence": "Finder hides file extensions — 'Evil.jpg.app' can masquerade as 'Evil.jpg'"}
}
