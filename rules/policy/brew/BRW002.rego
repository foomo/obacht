package obacht.brew

import rego.v1

findings contains f if {
	input.age_days >= 0
	input.age_days > 7
	f := {
		"rule_id": "BRW002",
		"evidence": sprintf("Homebrew metadata last refreshed %d days ago (>7)", [input.age_days]),
	}
}

skips contains s if {
	input.age_days < 0
	s := {
		"rule_id": "BRW002",
		"evidence": "Homebrew last-update timestamp unavailable",
	}
}
