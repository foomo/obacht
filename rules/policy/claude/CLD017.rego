package obacht.claude

import rego.v1

findings contains f if {
	input.env_disable_install_github_app_command != "1"
	f := {
		"rule_id": "CLD017",
		"evidence": sprintf("env.DISABLE_INSTALL_GITHUB_APP_COMMAND is %s, expected \"1\"", [input.env_disable_install_github_app_command]),
	}
}
