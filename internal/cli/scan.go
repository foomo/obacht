package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/franklinkim/bouncer/internal/collector"
	"github.com/franklinkim/bouncer/internal/preflight"
	"github.com/franklinkim/bouncer/internal/reporter"
	"github.com/franklinkim/bouncer/pkg/engine"
	"github.com/franklinkim/bouncer/pkg/schema"
	"github.com/franklinkim/bouncer/policies"
)

var (
	category string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan the local development environment for security issues",
	RunE:  runScan,
}

func init() {
	scanCmd.Flags().StringVar(&category, "category", "", "comma-separated list of categories to scan (e.g. ssh,git,env)")
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

// filterCollectors returns only collectors whose Name matches the category set.
func filterCollectors(collectors []collector.Collector, cats map[string]bool) []collector.Collector {
	if cats == nil {
		return collectors
	}

	var filtered []collector.Collector

	for _, c := range collectors {
		if cats[c.Name()] {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

// filterRules returns only rules whose Category matches the category set.
func filterRules(rules []schema.Rule, cats map[string]bool) []schema.Rule {
	if cats == nil {
		return rules
	}

	var filtered []schema.Rule

	for _, r := range rules {
		if cats[r.Category] {
			filtered = append(filtered, r)
		}
	}

	return filtered
}

func runScan(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Step 1: Preflight check for OPA.
	if err := preflight.CheckOPA(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "preflight: %v\n", err)
		os.Exit(Error)
	}

	// Step 2: Load built-in rules from embedded policies.
	rules, err := loadEmbeddedRules()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading embedded rules: %v\n", err)
		os.Exit(Error)
	}

	// Step 3: Load Rego files from embedded policies.
	regoFiles, err := loadEmbeddedRego()
	if err != nil {
		fmt.Fprintf(os.Stderr, "loading embedded rego: %v\n", err)
		os.Exit(Error)
	}

	// Step 4: Optionally load external rules from --rules-dir.
	if rulesDir != "" {
		extRules, extRego, err := loadExternalRules(rulesDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "loading external rules: %v\n", err)
			os.Exit(Error)
		}

		rules = mergeRules(rules, extRules)

		regoFiles = append(regoFiles, extRego...)
	}

	// Step 5: Parse category filter.
	cats := parseCategories()

	// Step 6: Filter rules by category.
	rules = filterRules(rules, cats)

	// Step 7: Create collectors.
	collectors := []collector.Collector{
		collector.NewSSHCollector(),
		collector.NewGitCollector(),
		collector.NewDockerCollector(),
		collector.NewKubeCollector(),
		collector.NewEnvCollector(),
		collector.NewShellCollector(),
		collector.NewToolsCollector(),
		collector.NewPathCollector(),
		collector.NewOSCollector(),
	}
	collectors = filterCollectors(collectors, cats)

	// Step 8: Run CollectAll.
	facts, collectorResults, err := collector.CollectAll(ctx, collectors)
	if err != nil {
		fmt.Fprintf(os.Stderr, "collecting facts: %v\n", err)
		os.Exit(Error)
	}

	// Step 7: Create engine and evaluate.
	eng, err := engine.NewEngine(regoFiles, rules)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating engine: %v\n", err)
		os.Exit(Error)
	}

	scanResult, err := eng.Evaluate(ctx, facts, collectorResults)
	if err != nil {
		fmt.Fprintf(os.Stderr, "evaluating policies: %v\n", err)
		os.Exit(Error)
	}

	// Step 8: Report results.
	rep := reporter.ForFormat(format)
	if err := rep.Report(os.Stdout, scanResult); err != nil {
		fmt.Fprintf(os.Stderr, "reporting results: %v\n", err)
		os.Exit(Error)
	}

	// Step 9: Exit with appropriate code.
	if scanResult.Summary.Failed > 0 {
		os.Exit(Findings)
	}

	return nil
}

// loadEmbeddedRules reads all YAML rule files from the embedded filesystem.
func loadEmbeddedRules() ([]schema.Rule, error) {
	var rules []schema.Rule

	err := fs.WalkDir(policies.Embedded, "rules", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}

		data, err := policies.Embedded.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		var rf schema.RulesFile
		if err := yaml.Unmarshal(data, &rf); err != nil {
			return fmt.Errorf("parsing %s: %w", path, err)
		}

		rules = append(rules, rf.Rules...)

		return nil
	})

	return rules, err
}

// loadEmbeddedRego reads all .rego files (excluding tests) from the embedded filesystem.
func loadEmbeddedRego() ([][]byte, error) {
	var regoFiles [][]byte

	err := fs.WalkDir(policies.Embedded, "rego", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".rego") || strings.HasSuffix(path, "_test.rego") {
			return nil
		}

		data, err := policies.Embedded.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		regoFiles = append(regoFiles, data)

		return nil
	})

	return regoFiles, err
}

// loadExternalRules loads rule YAML and Rego files from an external directory.
func loadExternalRules(dir string) ([]schema.Rule, [][]byte, error) {
	var (
		rules     []schema.Rule
		regoFiles [][]byte
	)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		switch {
		case strings.HasSuffix(entry.Name(), ".yaml"):
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, nil, fmt.Errorf("reading %s: %w", path, err)
			}

			var rf schema.RulesFile
			if err := yaml.Unmarshal(data, &rf); err != nil {
				return nil, nil, fmt.Errorf("parsing %s: %w", path, err)
			}

			rules = append(rules, rf.Rules...)

		case strings.HasSuffix(entry.Name(), ".rego") && !strings.HasSuffix(entry.Name(), "_test.rego"):
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, nil, fmt.Errorf("reading %s: %w", path, err)
			}

			regoFiles = append(regoFiles, data)
		}
	}

	return rules, regoFiles, nil
}

// mergeRules merges external rules into built-in rules. External rules with the
// same ID override built-in ones.
func mergeRules(builtIn, external []schema.Rule) []schema.Rule {
	// Build a map of external rules by ID.
	extMap := make(map[string]schema.Rule, len(external))
	for _, r := range external {
		extMap[r.ID] = r
	}

	// Override built-in rules where external ones exist.
	merged := make([]schema.Rule, 0, len(builtIn)+len(external))
	seen := make(map[string]bool)

	for _, r := range builtIn {
		if ext, ok := extMap[r.ID]; ok {
			merged = append(merged, ext)
		} else {
			merged = append(merged, r)
		}

		seen[r.ID] = true
	}

	// Append external rules that are entirely new.
	for _, r := range external {
		if !seen[r.ID] {
			merged = append(merged, r)
		}
	}

	return merged
}
