//go:build safe

package cli

import "context"

var (
	ParseRuleIDs         = parseRuleIDs
	CollectRuleIDs       = collectRuleIDs
	ValidateRuleIDs      = validateRuleIDs
	FilterRuleFilesByID  = filterRuleFilesByID
	ExcludeRuleFilesByID = excludeRuleFilesByID
)

// SetVersionForTest swaps the package-level version variable and returns the
// previous value so callers can restore it.
func SetVersionForTest(v string) string {
	prev := version
	version = v

	return prev
}

// RenderHeroForTest renders the scan model view with no rules so tests can
// assert hero/tagline/version output without standing up a full bubbletea run.
func RenderHeroForTest() string {
	m := newScanModel(context.Background(), nil)

	return m.View().Content
}
