package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.airdrop_setting == "everyone"
	f := {"rule_id": "OS017", "evidence": "AirDrop is set to Everyone"}
}
