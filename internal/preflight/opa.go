package preflight

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var versionRegexp = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

// CheckOPA verifies that the opa binary is available in PATH and that its
// version is >= 1.0.0. It returns a descriptive error with installation hints
// when the check fails.
func CheckOPA(ctx context.Context) error {
	path, err := exec.LookPath("opa")
	if err != nil {
		return fmt.Errorf(
			"opa not found in PATH: install it with \"brew install opa\" or \"mise install opa\"",
		)
	}

	out, err := exec.CommandContext(ctx, path, "version").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run opa version: %w", err)
	}

	major, minor, patch, err := parseOPAVersion(string(out))
	if err != nil {
		return err
	}

	if major < 1 {
		return fmt.Errorf(
			"opa version %d.%d.%d is too old (>= 1.0.0 required): upgrade with \"brew install opa\" or \"mise install opa\"",
			major, minor, patch,
		)
	}

	return nil
}

func parseOPAVersion(output string) (int, int, int, error) {
	for line := range strings.SplitSeq(output, "\n") {
		matches := versionRegexp.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		major, _ := strconv.Atoi(matches[1])
		minor, _ := strconv.Atoi(matches[2])
		patch, _ := strconv.Atoi(matches[3])

		return major, minor, patch, nil
	}

	return 0, 0, 0, fmt.Errorf("could not parse opa version from output: %s", output)
}
