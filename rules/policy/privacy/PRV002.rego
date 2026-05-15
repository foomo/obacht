package obacht.privacy

import rego.v1

findings contains f if {
	not input.vpn_configured
	f := {
		"rule_id": "PRV002",
		"evidence": "No VPN configuration or VPN process detected",
	}
}
