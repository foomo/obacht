package obacht.pnpm

import rego.v1

findings contains f if {
	not input.strict
	f := {
		"rule_id": "PNP004",
		"evidence": sprintf("pnpm strictDepBuilds is %q (expected \"true\")", [input.raw]),
	}
}
