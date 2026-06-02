package obacht.bun

import rego.v1

findings contains f if {
	not input.ignore_scripts
	f := {
		"rule_id": "BUN002",
		"evidence": sprintf("bun [install].ignoreScripts is %q (expected \"true\")", [input.raw]),
	}
}
