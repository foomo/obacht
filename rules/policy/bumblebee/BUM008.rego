package obacht.bumblebee

import rego.v1

findings contains f if {
	some entry in input.findings
	f := {
		"rule_id": "BUM008",
		"evidence": sprintf(
			"%s@%s matches %s (%s) at %s",
			[entry.package_name, entry.version, entry.catalog_id, entry.catalog_name, entry.source_file],
		),
	}
}
