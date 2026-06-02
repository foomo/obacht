package obacht.pnpm

import rego.v1

findings contains f if {
	not input.blocked
	f := {
		"rule_id": "PNP002",
		"evidence": sprintf("pnpm blockExoticSubdeps is %q (expected \"true\")", [input.raw]),
	}
}
