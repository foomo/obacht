package obacht.yarn

import rego.v1

findings contains f if {
	not input.scripts_disabled
	f := {
		"rule_id": "YRN002",
		"evidence": sprintf("yarn enableScripts is %q (expected \"false\")", [input.raw]),
	}
}
