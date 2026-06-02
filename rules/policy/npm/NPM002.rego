package obacht.npm

import rego.v1

findings contains f if {
	not input.ignore_scripts
	f := {
		"rule_id": "NPM002",
		"evidence": sprintf("npm ignore-scripts is %q (expected \"true\")", [input.raw]),
	}
}
