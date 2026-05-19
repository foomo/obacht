package runner

import (
	"encoding/json"
	"fmt"
)

// Envelope is the standard output shape for input scripts.
//
//	{"rule_id":"SSH001","status":"ok","data":{...}}
//	{"rule_id":"SSH001","status":"skip","skip_reason":"..."}
//
// During migration, scripts that emit bare JSON (no envelope) are accepted
// and treated as status=ok with the bare payload as Data. The bare-JSON
// fallback can be removed once every rules/inputs/*.sh emits the envelope
// shape — verify with `grep -L 'emit_ok\|emit_skip' rules/inputs/**/*.sh`.
type Envelope struct {
	RuleID     string `json:"rule_id"`
	Status     string `json:"status"`
	Data       any    `json:"data,omitempty"`
	SkipReason string `json:"skip_reason,omitempty"`
}

// ParseEnvelope parses raw stdout bytes into an Envelope.
// expectedID is the rule ID the caller expects; if non-empty, it must match
// the envelope's rule_id. When the payload is bare JSON (no envelope), the
// expectedID is used as Envelope.RuleID and Status defaults to "ok".
func ParseEnvelope(raw []byte, expectedID string) (Envelope, error) {
	var probe map[string]any
	if err := json.Unmarshal(raw, &probe); err != nil {
		return Envelope{}, fmt.Errorf("envelope not valid JSON: %w", err)
	}

	// Heuristic for envelope detection: presence of "status" with a string value
	// equal to "ok" or "skip". Anything else falls back to the legacy bare-JSON
	// path so legitimate legacy scripts (e.g. {"status":"online"}) are not
	// misclassified as malformed envelopes.
	statusStr, _ := probe["status"].(string)

	if statusStr != "ok" && statusStr != "skip" {
		return Envelope{
			RuleID: expectedID,
			Status: "ok",
			Data:   probe,
		}, nil
	}

	env := Envelope{}

	if err := json.Unmarshal(raw, &env); err != nil {
		return Envelope{}, fmt.Errorf("parsing envelope: %w", err)
	}

	if expectedID != "" && env.RuleID != "" && env.RuleID != expectedID {
		return Envelope{}, fmt.Errorf("rule_id mismatch: envelope=%q expected=%q", env.RuleID, expectedID)
	}

	if env.RuleID == "" {
		env.RuleID = expectedID
	}

	switch statusStr {
	case "ok":
		if env.Data == nil {
			return Envelope{}, fmt.Errorf("envelope status=ok requires data field")
		}
	case "skip":
		if env.SkipReason == "" {
			return Envelope{}, fmt.Errorf("envelope status=skip requires skip_reason field")
		}
	}

	return env, nil
}
