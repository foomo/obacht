package bouncer.env_test

import data.bouncer.env
import rego.v1

test_env001_suspicious_var_found if {
	findings := env.findings with input as {"env": {"suspicious_vars": [{"name": "GITHUB_TOKEN", "pattern": "exact:GITHUB_TOKEN"}]}}
	count(findings) == 1
	some f in findings
	f.rule_id == "ENV001"
}

test_env001_multiple_vars if {
	findings := env.findings with input as {"env": {"suspicious_vars": [
		{"name": "GITHUB_TOKEN", "pattern": "exact:GITHUB_TOKEN"},
		{"name": "MY_API_KEY", "pattern": "*_API_KEY"},
	]}}
	count(findings) == 2
}

test_env001_no_suspicious_vars if {
	findings := env.findings with input as {"env": {"suspicious_vars": []}}
	count(findings) == 0
}
