package obacht.pnpm

import rego.v1

findings contains f if {
	not input.no_downgrade
	f := {
		"rule_id": "PNP003",
		"evidence": sprintf("pnpm trustPolicy is %q (expected \"no-downgrade\")", [input.raw]),
	}
}
