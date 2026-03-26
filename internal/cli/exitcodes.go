package cli

// Exit codes for the bouncer CLI.
const (
	OK       = 0 // No findings
	Findings = 1 // Security findings detected
	Error    = 2 // Runtime error
)
