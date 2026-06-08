package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.multicast_disabled
	f := {"rule_id": "OS045", "evidence": "Bonjour multicast advertisements are enabled (machine advertises name and services on the local network)"}
}
