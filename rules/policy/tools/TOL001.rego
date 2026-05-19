package obacht.tools

import rego.v1

findings contains f if {
	some tool in input.tools
	not tool.installed
	f := {
		"rule_id": "TOL001",
		"evidence": sprintf("Tool '%s' is not installed", [tool.name]),
	}
}
