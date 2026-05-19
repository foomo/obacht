package obacht.tools

import rego.v1

findings contains f if {
	not input.telemetry_disabled
	f := {
		"rule_id": "TOL005",
		"evidence": sprintf("Go telemetry mode is %q (expected \"off\")", [input.mode]),
	}
}
