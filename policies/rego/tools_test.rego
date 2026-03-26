package bouncer.tools_test

import data.bouncer.tools
import rego.v1

test_tol001_missing_tool if {
	findings := tools.findings with input as {"tools": {"tools": [
		{"name": "git", "installed": false, "version": "", "path": ""},
		{"name": "opa", "installed": true, "version": "0.62.0", "path": "/usr/local/bin/opa"},
	]}}
	count(findings) == 1
	some f in findings
	f.rule_id == "TOL001"
}

test_tol001_all_installed if {
	findings := tools.findings with input as {"tools": {"tools": [
		{"name": "git", "installed": true, "version": "2.45.0", "path": "/usr/bin/git"},
		{"name": "opa", "installed": true, "version": "0.62.0", "path": "/usr/local/bin/opa"},
	]}}
	count(findings) == 0
}

test_tol001_multiple_missing if {
	findings := tools.findings with input as {"tools": {"tools": [
		{"name": "git", "installed": false, "version": "", "path": ""},
		{"name": "opa", "installed": false, "version": "", "path": ""},
		{"name": "gpg", "installed": true, "version": "2.4.5", "path": "/usr/bin/gpg"},
	]}}
	count(findings) == 2
}
