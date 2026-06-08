package obacht.mise

import rego.v1

findings contains f if {
	input.value
	f := {
		"rule_id": "MIS002",
		"evidence": "mise not_found_auto_install is true (default); should be false to block silent supply-chain installs",
	}
}
