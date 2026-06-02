package obacht.npm

import rego.v1

findings contains f if {
	not input.restricted
	f := {
		"rule_id": "NPM003",
		"evidence": sprintf("npm allow-git is %q (expected \"none\" or \"root\")", [input.raw]),
	}
}
