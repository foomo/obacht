package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.airdrop_setting != "off"
	input.airdrop_setting != "everyone"
	f := {"rule_id": "OS022", "evidence": sprintf("AirDrop is set to %s; consider disabling entirely", [input.airdrop_setting])}
}
