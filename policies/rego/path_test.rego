package bouncer.path_test

import data.bouncer.path
import rego.v1

test_pth001_writable_dir if {
	findings := path.findings with input as {"path": {"dirs": [
		{"path": "/usr/local/bin", "exists": true, "writable": true, "is_relative": false},
		{"path": "/usr/bin", "exists": true, "writable": false, "is_relative": false},
	]}}
	count(findings) == 1
	some f in findings
	f.rule_id == "PTH001"
}

test_pth002_relative_path if {
	findings := path.findings with input as {"path": {"dirs": [
		{"path": "relative/bin", "exists": false, "writable": false, "is_relative": true},
		{"path": "/usr/bin", "exists": true, "writable": false, "is_relative": false},
	]}}
	count(findings) == 1
	some f in findings
	f.rule_id == "PTH002"
}

test_pth001_and_pth002_combined if {
	findings := path.findings with input as {"path": {"dirs": [
		{"path": "/tmp", "exists": true, "writable": true, "is_relative": false},
		{"path": ".", "exists": true, "writable": true, "is_relative": true},
	]}}

	# /tmp triggers PTH001, "." triggers PTH002 and PTH001
	count(findings) == 3
}

test_no_findings_clean_path if {
	findings := path.findings with input as {"path": {"dirs": [
		{"path": "/usr/bin", "exists": true, "writable": false, "is_relative": false},
		{"path": "/usr/local/bin", "exists": true, "writable": false, "is_relative": false},
	]}}
	count(findings) == 0
}
