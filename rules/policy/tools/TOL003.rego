package obacht.tools

import rego.v1

findings contains f if {
	some pm in input.package_managers
	pm.age_days >= 0
	pm.age_days > 7
	f := {
		"rule_id": "TOL003",
		"evidence": sprintf("%s metadata last refreshed %d days ago (>7)", [pm.name, pm.age_days]),
	}
}

skips contains s if {
	some pm in input.package_managers
	pm.age_days < 0
	s := {
		"rule_id": "TOL003",
		"evidence": sprintf("%s last-update timestamp unavailable", [pm.name]),
	}
}
