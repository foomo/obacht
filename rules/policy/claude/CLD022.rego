package obacht.claude

import rego.v1

findings contains f if {
	input.env_force_autoupdate_plugins != "1"
	f := {
		"rule_id": "CLD022",
		"evidence": sprintf("env.FORCE_AUTOUPDATE_PLUGINS is %s, expected \"1\"", [input.env_force_autoupdate_plugins]),
	}
}
