package obacht.gh

import rego.v1

telemetry_disabled if input.telemetry == "false"

telemetry_disabled if input.gh_telemetry_env in {"0", "false", "off", "no"}

telemetry_disabled if input.do_not_track_env in {"1", "true"}

findings contains f if {
	not telemetry_disabled
	f := {
		"rule_id": "GH001",
		"evidence": sprintf(
			"gh telemetry is not disabled (config=%q, GH_TELEMETRY=%q, DO_NOT_TRACK=%q)",
			[input.telemetry, input.gh_telemetry_env, input.do_not_track_env],
		),
	}
}
