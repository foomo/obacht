package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.bluetooth_sharing_disabled
	f := {"rule_id": "OS026", "evidence": "Bluetooth Sharing is enabled"}
}
