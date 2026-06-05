package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.mac_randomization_disabled_count > 0
	f := {
		"rule_id": "OS043",
		"evidence": sprintf("%d known Wi-Fi network(s) have private MAC address disabled", [input.mac_randomization_disabled_count]),
	}
}
