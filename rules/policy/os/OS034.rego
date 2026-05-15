package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.airplay_receiver_enabled
	f := {"rule_id": "OS034", "evidence": "AirPlay Receiver is enabled"}
}
