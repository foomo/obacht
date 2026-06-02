package obacht.apt

import rego.v1

findings contains f if {
	input.age_days >= 0
	input.age_days > 7
	f := {
		"rule_id": "APT001",
		"evidence": sprintf("APT metadata last refreshed %d days ago (>7)", [input.age_days]),
	}
}

skips contains s if {
	input.age_days < 0
	s := {
		"rule_id": "APT001",
		"evidence": "APT last-update timestamp unavailable",
	}
}
