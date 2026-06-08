package obacht.mise

import rego.v1

findings contains f if {
	input.value != "always"
	f := {
		"rule_id": "MIS003",
		"evidence": sprintf("mise status.missing_tools is %q (recommend \"always\")", [input.value]),
	}
}
