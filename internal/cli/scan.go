package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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

// ruleIDPattern matches the conventional rule-ID format: uppercase prefix +
// digits (e.g. SSH005, OS030, PTH001). Used to reject any rule.ID that could
// be used for path traversal when looking up per-rule input/policy files.
var ruleIDPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*[0-9]+$`)

// validateRuleID returns an error if id does not match the conventional
// PREFIX### pattern. Used to gate filesystem lookups whose path depends
// on the rule ID.
func validateRuleID(id string) error {
	if !ruleIDPattern.MatchString(id) {
		return fmt.Errorf("invalid rule ID %q: must match %s", id, ruleIDPattern.String())
	}

	return nil
}

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
	scanCmd.Flags().StringVar(&severity, "severity", "", "comma-separated list of severities to include (critical,high,warn,info,note)")
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

	// When --rules-dir is set, use only those rules. Otherwise use embedded rules.
	var (
		ruleFiles []schema.RulesFile
		err       error
	)

	var rulesFS fs.FS

	if rulesDir != "" {
		ruleFiles, err = loadExternalRuleFiles(rulesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "loading external rules: %v\n", err)
			os.Exit(Error)
		}

		rulesFS = os.DirFS(rulesDir)
	} else {
		ruleFiles, err = loadEmbeddedRuleFiles()
		if err != nil {
			fmt.Fprintf(os.Stderr, "loading embedded rules: %v\n", err)
			os.Exit(Error)
		}

		rulesFS = rules.Embedded
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

	// Materialize the embedded bumblebee exposure catalog under the user cache
	// directory and export its path so bumblebee-category input scripts can
	// find it. Persistent cache (rewritten each scan) — no cleanup needed,
	// avoids the defer-vs-os.Exit hazard.
	if catalogDir, err := materializeBumblebeeCatalog(); err == nil {
		_ = os.Setenv("OBACHT_BUMBLEBEE_CATALOG_DIR", catalogDir)
	} else {
		fmt.Fprintf(os.Stderr, "preparing bumblebee catalog: %v (bumblebee rules will skip)\n", err)
	}

	// Evaluate all rule files.
	var scanResult *schema.ScanResult

	if format == "pretty" {
		// Run scan with animated progress display.
		model := newScanModel(ctx, rulesFS, ruleFiles)
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

		scanResult, err = engine.EvaluateWithFS(ctx, rulesFS, ruleFiles)
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

// loadEmbeddedRuleFiles reads all YAML rule files from the embedded filesystem
// and resolves per-rule input/policy files where not set inline.
func loadEmbeddedRuleFiles() ([]schema.RulesFile, error) {
	var ruleFiles []schema.RulesFile

	err := fs.WalkDir(rules.Embedded, "policies", func(path string, d fs.DirEntry, err error) error {
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

		baseName := strings.TrimSuffix(filepath.Base(path), ".yaml")

		// Category-shared input: inputs/<cat>.sh
		if rf.Input == "" {
			shared, err := resolveInputFromFS(rules.Embedded, "inputs", baseName)
			if err != nil {
				return fmt.Errorf("resolving shared input for %s: %w", path, err)
			}

			rf.Input = shared
		}

		// Category-shared policy bundle: policy/<cat>/*.rego
		// Engaged only when a shared input exists; otherwise per-rule policy
		// resolution below handles the standard layout.
		bundled := false

		if rf.Input != "" && rf.Policy == "" {
			bundle, err := bundleRegoFromFS(rules.Embedded, "policy/"+baseName)
			if err != nil {
				return fmt.Errorf("bundling shared policy for %s: %w", path, err)
			}

			if bundle != "" {
				rf.Policy = bundle
				bundled = true
			}
		}

		// Per-rule resolution.
		for i := range rf.Rules {
			rule := &rf.Rules[i]

			if err := validateRuleID(rule.ID); err != nil {
				return fmt.Errorf("in %s: %w", path, err)
			}

			if rule.Input == "" {
				perRule, err := resolveInputFromFS(rules.Embedded, "inputs/"+baseName, rule.ID)
				if err != nil {
					return fmt.Errorf("resolving input for rule %s: %w", rule.ID, err)
				}

				rule.Input = perRule
			}

			if !bundled && rule.Policy == "" {
				perRule, err := resolveRegoFromFS(rules.Embedded, "policy/"+baseName, rule.ID)
				if err != nil {
					return fmt.Errorf("resolving policy for rule %s: %w", rule.ID, err)
				}

				rule.Policy = perRule
			}
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

		baseName := strings.TrimSuffix(entry.Name(), ".yaml")

		// Category-shared input: <dir>/inputs/<cat>.sh
		if rf.Input == "" {
			sharedPath := filepath.Join(dir, "inputs", baseName+".sh")

			data, err := os.ReadFile(sharedPath)
			if err == nil {
				rf.Input = string(data)
			} else if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("reading shared input %s: %w", sharedPath, err)
			}
		}

		// Category-shared policy bundle: <dir>/policy/<cat>/*.rego
		bundled := false

		if rf.Input != "" && rf.Policy == "" {
			bundle, err := bundleRegoFromFS(os.DirFS(dir), "policy/"+baseName)
			if err != nil {
				return nil, fmt.Errorf("bundling shared policy for %s: %w", path, err)
			}

			if bundle != "" {
				rf.Policy = bundle
				bundled = true
			}
		}

		// Per-rule resolution from inputs/<cat>/<RULEID>.sh and policy/<cat>/<RULEID>.rego.
		for i := range rf.Rules {
			rule := &rf.Rules[i]

			if err := validateRuleID(rule.ID); err != nil {
				return nil, fmt.Errorf("in %s: %w", path, err)
			}

			if rule.Input == "" {
				scriptPath := filepath.Join(dir, "inputs", baseName, rule.ID+".sh")

				data, err := os.ReadFile(scriptPath)
				if err == nil {
					rule.Input = string(data)
				} else if !errors.Is(err, os.ErrNotExist) {
					return nil, fmt.Errorf("reading per-rule input %s: %w", scriptPath, err)
				}
			}

			if !bundled && rule.Policy == "" {
				regoPath := filepath.Join(dir, "policy", baseName, rule.ID+".rego")

				data, err := os.ReadFile(regoPath)
				if err == nil {
					rule.Policy = string(data)
				} else if !errors.Is(err, os.ErrNotExist) {
					return nil, fmt.Errorf("reading per-rule rego %s: %w", regoPath, err)
				}
			}
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

// materializeBumblebeeCatalog writes the embedded bumblebee exposure catalog
// JSONs into a stable subdirectory of the user cache dir and returns its
// absolute path. Each scan rewrites the files and prunes any stale entries
// left from a previous obacht version. A persistent cache avoids the
// defer-vs-os.Exit hazard of a per-run temp dir.
func materializeBumblebeeCatalog() (string, error) {
	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("locating user cache dir: %w", err)
	}

	dir := filepath.Join(cacheRoot, "obacht", "bumblebee-catalog")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("creating cache dir: %w", err)
	}

	entries, err := fs.ReadDir(rules.BumblebeeCatalog, "catalogs/bumblebee")
	if err != nil {
		return "", fmt.Errorf("reading embedded catalog: %w", err)
	}

	want := make(map[string]struct{}, len(entries))

	wrote := 0

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}

		data, err := fs.ReadFile(rules.BumblebeeCatalog, "catalogs/bumblebee/"+e.Name())
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", e.Name(), err)
		}

		if err := os.WriteFile(filepath.Join(dir, e.Name()), data, 0o600); err != nil {
			return "", fmt.Errorf("writing %s: %w", e.Name(), err)
		}

		want[e.Name()] = struct{}{}
		wrote++
	}

	if wrote == 0 {
		return "", errors.New("no catalog files embedded")
	}

	if existing, err := os.ReadDir(dir); err == nil {
		for _, e := range existing {
			if e.IsDir() {
				continue
			}

			if _, keep := want[e.Name()]; !keep {
				_ = os.Remove(filepath.Join(dir, e.Name()))
			}
		}
	}

	return dir, nil
}

// bundleRegoFromFS reads every <dir>/*.rego in fsys and concatenates them into
// a single valid rego module. The first file is kept verbatim (including its
// `package` + `import` header); subsequent files have their leading
// package/import/comment block stripped so the result has exactly one header.
// Files are processed in lexicographic order for determinism — name shared
// helpers with a leading underscore (e.g. `_helpers.rego`) to anchor them
// first. Returns "" when the directory is absent or contains no .rego files.
func bundleRegoFromFS(fsys fs.FS, dir string) (string, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil
		}

		return "", err
	}

	var files []string

	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".rego") {
			files = append(files, e.Name())
		}
	}

	if len(files) == 0 {
		return "", nil
	}

	sort.Strings(files)

	parts := make([]string, 0, len(files))

	for i, name := range files {
		data, err := fs.ReadFile(fsys, dir+"/"+name)
		if err != nil {
			return "", err
		}

		body := strings.TrimRight(string(data), "\n")

		if i > 0 {
			body = stripRegoHeader(body)
		}

		if body != "" {
			parts = append(parts, body)
		}
	}

	return strings.Join(parts, "\n\n") + "\n", nil
}

// stripRegoHeader removes the leading package/import/comment/blank block from
// a rego source so it can be appended after another rego module that already
// carries the canonical header.
func stripRegoHeader(src string) string {
	lines := strings.Split(src, "\n")

	start := 0
	for ; start < len(lines); start++ {
		t := strings.TrimSpace(lines[start])
		if t == "" || strings.HasPrefix(t, "#") || strings.HasPrefix(t, "package ") || strings.HasPrefix(t, "import ") {
			continue
		}

		break
	}

	return strings.TrimRight(strings.Join(lines[start:], "\n"), "\n")
}

// resolveRegoFromFS checks for <dir>/<RULEID>.rego in fsys and returns its
// content if found. Returns empty string if not found.
func resolveRegoFromFS(fsys fs.FS, dir, ruleID string) (string, error) {
	path := dir + "/" + ruleID + ".rego"

	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil
		}

		return "", err
	}

	return string(data), nil
}
