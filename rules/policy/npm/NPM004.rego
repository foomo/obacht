package obacht.npm

import rego.v1

findings contains f if {
	not input.secure
	f := {
		"rule_id": "NPM004",
		"evidence": sprintf("npm registry is %q (expected https://...)", [input.raw]),
	}
}
