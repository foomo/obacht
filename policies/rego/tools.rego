package bouncer.tools

import rego.v1

findings contains f if {
	tool := input.tools.tools[_]
	not tool.installed
	f := {
		"rule_id": "TOL001",
		"evidence": sprintf("Tool '%s' is not installed", [tool.name]),
	}
}
