package obacht.os

import rego.v1

findings contains f if {
	input.os == "darwin"
	input.env_keep_home
	f := {"rule_id": "OS041", "evidence": "sudoers env_keep retains HOME — allows non-root dotfiles to execute as root via 'sudo zsh'"}
}
