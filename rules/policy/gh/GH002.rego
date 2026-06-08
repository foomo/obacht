package obacht.gh

import rego.v1

expected_helper := "gh auth git-credential"

findings contains f if {
	not contains(input.github_helpers, expected_helper)
	f := {
		"rule_id": "GH002",
		"evidence": "git credential helper for https://github.com does not delegate to `gh auth git-credential`",
	}
}

findings contains f if {
	not contains(input.gist_helpers, expected_helper)
	f := {
		"rule_id": "GH002",
		"evidence": "git credential helper for https://gist.github.com does not delegate to `gh auth git-credential`",
	}
}
