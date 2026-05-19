package obacht.credentials

import rego.v1

findings contains f if {
	input.gcp_exists
	input.gcp_mode != "0600"
	f := {
		"rule_id": "CRD003",
		"evidence": sprintf("GCP credentials file has mode %s (expected 0600)", [input.gcp_mode]),
	}
}
