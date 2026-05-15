package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	not input.edr_deployed
	f := {"rule_id": "OS018", "evidence": "No Endpoint Detection & Response agent deployed"}
}
