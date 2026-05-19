package obacht.privacy

import rego.v1

findings contains f if {
	not input.trusted_dns
	count(input.dns_servers) > 0
	f := {
		"rule_id": "PRV004",
		"evidence": sprintf("Untrusted DNS resolver(s): %s", [input.untrusted_dns_servers]),
	}
}

findings contains f if {
	not input.trusted_dns
	count(input.dns_servers) == 0
	f := {
		"rule_id": "PRV004",
		"evidence": "No DNS resolvers detected via scutil --dns",
	}
}
