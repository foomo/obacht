package bouncer.kube

import rego.v1

findings contains f if {
	input.kube.config_exists
	input.kube.config_mode != "0600"
	f := {
		"rule_id": "KUB001",
		"evidence": sprintf("~/.kube/config has mode %s (expected 0600)", [input.kube.config_mode]),
	}
}

findings contains f if {
	ctx := input.kube.contexts[_]
	contains(lower(ctx.name), "prod")
	f := {
		"rule_id": "KUB002",
		"evidence": sprintf("Production context found: %s (cluster: %s)", [ctx.name, ctx.cluster]),
	}
}
