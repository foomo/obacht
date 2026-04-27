package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/foomo/obacht/internal/reporter"
	"github.com/foomo/obacht/pkg/engine"
	"github.com/foomo/obacht/pkg/schema"
	"github.com/foomo/obacht/rules"
)

var (
	category    string
	severity    string
	rule        string
	excludeRule string
	showPassing bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan the local development environment for security issues",
	RunE:  runScan,
}

func init() {
	scanCmd.Flags().StringVar(&category, "category", "", "comma-separated list of categories to scan (e.g. ssh,git,env)")
	scanCmd.Flags().StringVar(&severity, "severity", "", "comma-separated list of severities to include (critical,high,warn,info)")
	scanCmd.Flags().StringVar(&rule, "rule", "", "comma-separated list of rule IDs to run (e.g. SSH001,GIT003)")
	scanCmd.Flags().StringVar(&excludeRule, "exclude-rule", "", "comma-separated list of rule IDs to exclude (e.g. SSH001,GIT003)")
	scanCmd.Flags().BoolVar(&showPassing, "show-passing", false, "include passing checks in pretty output (no effect on --format json)")
	rootCmd.AddCommand(scanCmd)
}

// parseCategories splits the comma-separated category flag into a set.
func parseCategories() map[string]bool {
	if category == "" {
		return nil
	}

	cats := make(map[string]bool)

	for c := range strings.SplitSeq(category, ",") {
		c = strings.TrimSpace(c)
		if c != "" {
			cats[c] = true
		}
	}

	if len(cats) == 0 {
		return nil
	}

	return cats
}

// parseSeverities splits the comma-separated severity flag into a set.
func parseSeverities() map[schema.Severity]bool {
	if severity == "" {
		return nil
	}

	sevs := make(map[schema.Severity]bool)

	for s := range strings.SplitSeq(severity, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			sevs[schema.Severity(s)] = true
		}
	}

	if len(sevs) == 0 {
		return nil
	}

	return sevs
}

// parseRuleIDs splits a comma-separated rule ID string into a set.
func parseRuleIDs(s string) map[string]bool {
	if s == "" {
		return nil
	}

	ids := make(map[string]bool)

	for id := range strings.SplitSeq(s, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			ids[id] = true
		}
	}

	if len(ids) == 0 {
		return nil
	}

	return ids
}

// collectRuleIDs returns the set of all rule IDs across the given rule files.
func collectRuleIDs(ruleFiles []schema.RulesFile) map[string]bool {
	ids := make(map[string]bool)

	for _, rf := range ruleFiles {
		for _, r := range rf.Rules {
			ids[r.ID] = true
		}
	}

	return ids
}

// validateRuleIDs returns an error if any key in requested is not present in known.
func validateRuleIDs(requested, known map[string]bool) error {
	var unknown []string

	for id := range requested {
		if !known[id] {
			unknown = append(unknown, id)
		}
	}

	if len(unknown) == 0 {
		return nil
	}

	sort.Strings(unknown)

	return fmt.Errorf("unknown rule IDs: %s", strings.Join(unknown, ", "))
}

// filterRuleFilesByID returns rule files with only the rules whose ID is in the set.
func filterRuleFilesByID(ruleFiles []schema.RulesFile, ids map[string]bool) []schema.RulesFile {
	if ids == nil {
		return ruleFiles
	}

	var filtered []schema.RulesFile

	for _, rf := range ruleFiles {
		var rules []schema.Rule

		for _, r := range rf.Rules {
			if ids[r.ID] {
				rules = append(rules, r)
			}
		}

		if len(rules) > 0 {
			filtered = append(filtered, schema.RulesFile{
				Input:  rf.Input,
				Policy: rf.Policy,
				Rules:  rules,
			})
		}
	}

	return filtered
}

// excludeRuleFilesByID returns rule files with rules whose ID is not in the set.
func excludeRuleFilesByID(ruleFiles []schema.RulesFile, ids map[string]bool) []schema.RulesFile {
	if ids == nil {
		return ruleFiles
	}

	var filtered []schema.RulesFile

	for _, rf := range ruleFiles {
		var rules []schema.Rule

		for _, r := range rf.Rules {
			if !ids[r.ID] {
				rules = append(rules, r)
			}
		}

		if len(rules) > 0 {
			filtered = append(filtered, schema.RulesFile{
				Input:  rf.Input,
				Policy: rf.Policy,
				Rules:  rules,
			})
		}
	}

	return filtered
}

// filterRuleFiles returns rule files with only the rules matching the category set.
func filterRuleFiles(ruleFiles []schema.RulesFile, cats map[string]bool) []schema.RulesFile {
	if cats == nil {
		return ruleFiles
	}

	var filtered []schema.RulesFile

	for _, rf := range ruleFiles {
		var rules []schema.Rule

		for _, r := range rf.Rules {
			if cats[r.Category] {
				rules = append(rules, r)
			}
		}

		if len(rules) > 0 {
			filtered = append(filtered, schema.RulesFile{
				Input:  rf.Input,
				Policy: rf.Policy,
				Rules:  rules,
			})
		}
	}

	return filtered
}

// filterRuleFilesBySeverity returns rule files with only the rules matching the severity set.
func filterRuleFilesBySeverity(ruleFiles []schema.RulesFile, sevs map[schema.Severity]bool) []schema.RulesFile {
	if sevs == nil {
		return ruleFiles
	}

	var filtered []schema.RulesFile

	for _, rf := range ruleFiles {
		var rules []schema.Rule

		for _, r := range rf.Rules {
			if sevs[r.Severity] {
				rules = append(rules, r)
			}
		}

		if len(rules) > 0 {
			filtered = append(filtered, schema.RulesFile{
				Input:  rf.Input,
				Policy: rf.Policy,
				Rules:  rules,
			})
		}
	}

	return filtered
}

func runScan(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Load built-in rule files from embedded policies.
	ruleFiles, err := loadEmbeddedRuleFiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading embedded rules: %v\n", err)
		os.Exit(Error)
	}

	// Optionally load external rule files from --rules-dir.
	if rulesDir != "" {
		extRuleFiles, err := loadExternalRuleFiles(rulesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "loading external rules: %v\n", err)
			os.Exit(Error)
		}

		ruleFiles = mergeRuleFiles(ruleFiles, extRuleFiles)
	}

	// Parse rule ID filters.
	ruleIDs := parseRuleIDs(rule)
	excludeIDs := parseRuleIDs(excludeRule)

	// --rule is mutually exclusive with --category and --severity.
	if ruleIDs != nil && (category != "" || severity != "") {
		return fmt.Errorf("--rule cannot be combined with --category or --severity")
	}

	// Validate rule IDs against the full set of loaded rules.
	allIDs := collectRuleIDs(ruleFiles)

	if err := validateRuleIDs(ruleIDs, allIDs); err != nil {
		return err
	}

	if err := validateRuleIDs(excludeIDs, allIDs); err != nil {
		return err
	}

	// Filter by rule ID (allowlist).
	ruleFiles = filterRuleFilesByID(ruleFiles, ruleIDs)

	// Filter by category.
	cats := parseCategories()
	ruleFiles = filterRuleFiles(ruleFiles, cats)

	// Filter by severity.
	sevs := parseSeverities()
	ruleFiles = filterRuleFilesBySeverity(ruleFiles, sevs)

	// Exclude specific rules (blocklist, applied last).
	ruleFiles = excludeRuleFilesByID(ruleFiles, excludeIDs)

	// Evaluate all rule files.
	var scanResult *schema.ScanResult

	if format == "pretty" {
		// Run scan with animated progress display.
		model := newScanModel(ctx, ruleFiles)
		p := tea.NewProgram(model, tea.WithOutput(os.Stderr), tea.WithInput(os.Stdin))
		model.SetProgram(p)

		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "running progress: %v\n", err)
			os.Exit(Error)
		}

		m := finalModel.(*scanModel) //nolint:forcetypeassert
		if m.err != nil {
			fmt.Fprintf(os.Stderr, "evaluating policies: %v\n", m.err)
			os.Exit(Error)
		}

		scanResult = m.result
	} else {
		var err error

		scanResult, err = engine.Evaluate(ctx, ruleFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "evaluating policies: %v\n", err)
			os.Exit(Error)
		}
	}

	// Report results.
	rep := reporter.ForFormat(format, showPassing)
	if err := rep.Report(os.Stdout, scanResult); err != nil {
		fmt.Fprintf(os.Stderr, "reporting results: %v\n", err)
		os.Exit(Error)
	}

	// Exit with appropriate code.
	if scanResult.Summary.Failed > 0 {
		os.Exit(Findings)
	}

	return nil
}

// loadEmbeddedRuleFiles reads all YAML rule files from the embedded filesystem.
func loadEmbeddedRuleFiles() ([]schema.RulesFile, error) {
	var ruleFiles []schema.RulesFile

	err := fs.WalkDir(rules.Embedded, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		data, err := rules.Embedded.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		var rf schema.RulesFile
		if err := yaml.Unmarshal(data, &rf); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		// Resolve input from inputs/<name>.sh if not set inline.
		if rf.Input == "" {
			baseName := strings.TrimSuffix(filepath.Base(path), ".yaml")

			input, err := resolveInputFromFS(rules.Embedded, "inputs", baseName)
			if err != nil {
				return fmt.Errorf("resolving input for %s: %w", path, err)
			}

			rf.Input = input
		}

		ruleFiles = append(ruleFiles, rf)

		return nil
	})

	return ruleFiles, err
}

// loadExternalRuleFiles loads rule YAML files from an external directory.
// The directory is expected to contain a policies/ subdirectory with YAML files
// and an optional inputs/ subdirectory with shell scripts.
func loadExternalRuleFiles(dir string) ([]schema.RulesFile, error) {
	var ruleFiles []schema.RulesFile

	policiesDir := filepath.Join(dir, "policies")

	entries, err := os.ReadDir(policiesDir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w (expected policies/ subdirectory in %s)", policiesDir, err, dir)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		path := filepath.Join(policiesDir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}

		var rf schema.RulesFile
		if err := yaml.Unmarshal(data, &rf); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}

		// Resolve input from inputs/<name>.sh if not set inline.
		if rf.Input == "" {
			baseName := strings.TrimSuffix(entry.Name(), ".yaml")

			input, err := resolveInputFromDir(dir, baseName)
			if err != nil {
				return nil, fmt.Errorf("resolving input for %s: %w", path, err)
			}

			rf.Input = input
		}

		// Resolve policy file references.
		rf.Policy, err = resolvePolicy(policiesDir, rf.Policy)
		if err != nil {
			return nil, fmt.Errorf("resolving policy in %s: %w", path, err)
		}

		for i := range rf.Rules {
			rf.Rules[i].Policy, err = resolvePolicy(policiesDir, rf.Rules[i].Policy)
			if err != nil {
				return nil, fmt.Errorf("resolving policy for rule %s in %s: %w", rf.Rules[i].ID, path, err)
			}
		}

		ruleFiles = append(ruleFiles, rf)
	}

	return ruleFiles, nil
}

// resolvePolicy checks if the policy value is a file reference (ends with .rego)
// and reads the file content. Otherwise returns the value as-is (inline rego).
func resolvePolicy(dir, policy string) (string, error) {
	if policy == "" || !strings.HasSuffix(policy, ".rego") {
		return policy, nil
	}

	data, err := os.ReadFile(filepath.Join(dir, policy))
	if err != nil {
		return "", fmt.Errorf("reading rego file %s: %w", policy, err)
	}

	return string(data), nil
}

// resolveInputFromFS checks for inputs/<name>.sh in the given FS and returns
// its content if found. Returns empty string if not found.
func resolveInputFromFS(fsys fs.FS, inputsDir, baseName string) (string, error) {
	scriptPath := inputsDir + "/" + baseName + ".sh"

	data, err := fs.ReadFile(fsys, scriptPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil
		}

		return "", err
	}

	return string(data), nil
}

// resolveInputFromDir checks for <dir>/inputs/<name>.sh on the real filesystem.
func resolveInputFromDir(dir, baseName string) (string, error) {
	scriptPath := filepath.Join(dir, "inputs", baseName+".sh")

	data, err := os.ReadFile(scriptPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", err
	}

	return string(data), nil
}

// mergeRuleFiles merges external rule files into built-in ones.
// External rules with the same ID override built-in ones.
func mergeRuleFiles(builtIn, external []schema.RulesFile) []schema.RulesFile {
	// Build a map of all external rules by ID.
	extMap := make(map[string]schema.Rule)

	for _, rf := range external {
		for _, r := range rf.Rules {
			extMap[r.ID] = r
		}
	}

	// Override built-in rules where external ones have the same ID.
	for i, rf := range builtIn {
		for j, r := range rf.Rules {
			if ext, ok := extMap[r.ID]; ok {
				builtIn[i].Rules[j] = ext

				delete(extMap, r.ID)
			}
		}
	}

	// Collect remaining external rules (entirely new) into rule files.
	if len(extMap) > 0 {
		// Add them as part of the external rule files that define their input/policy.
		for _, rf := range external {
			var newRules []schema.Rule

			for _, r := range rf.Rules {
				if _, ok := extMap[r.ID]; ok {
					newRules = append(newRules, r)
				}
			}

			if len(newRules) > 0 {
				builtIn = append(builtIn, schema.RulesFile{
					Input:  rf.Input,
					Policy: rf.Policy,
					Rules:  newRules,
				})
			}
		}
	}

	return builtIn
}
