package obacht.kube

import rego.v1

findings contains f if {
	some ctx in input.contexts
	contains(lower(ctx.name), "prod")
	f := {
		"rule_id": "KUB002",
		"evidence": sprintf("Production context found: %s (cluster: %s)", [ctx.name, ctx.cluster]),
	}
}
