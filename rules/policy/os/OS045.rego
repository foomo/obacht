package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.save_to_icloud
	f := {"rule_id": "OS045", "evidence": "New documents default to saving in iCloud Drive"}
}
