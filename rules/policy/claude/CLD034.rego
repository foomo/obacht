package obacht.claude

import rego.v1

findings contains f if {
	input.sandbox_network_allow_managed_domains_only != "true"
	f := {
		"rule_id": "CLD034",
		"evidence": sprintf("sandbox.network.allowManagedDomainsOnly is %s, expected true", [input.sandbox_network_allow_managed_domains_only]),
	}
}
