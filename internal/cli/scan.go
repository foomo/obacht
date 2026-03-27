package cli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

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

	// Filter by category.
	cats := parseCategories()
	ruleFiles = filterRuleFiles(ruleFiles, cats)

	// Evaluate all rule files.
	scanResult, err := engine.Evaluate(ctx, ruleFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "evaluating policies: %v\n", err)
		os.Exit(Error)
	}

	// Report results.
	rep := reporter.ForFormat(format)
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

	err := fs.WalkDir(policies.Embedded, ".", func(path string, d fs.DirEntry, err error) error {
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

		ruleFiles = append(ruleFiles, rf)

		return nil
	})

	return ruleFiles, err
}

// loadExternalRuleFiles loads rule YAML files from an external directory.
// It also resolves policy file references (non-inline rego) relative to the directory.
func loadExternalRuleFiles(dir string) ([]schema.RulesFile, error) {
	var ruleFiles []schema.RulesFile

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}

		var rf schema.RulesFile
		if err := yaml.Unmarshal(data, &rf); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}

		// Resolve policy file references.
		rf.Policy, err = resolvePolicy(dir, rf.Policy)
		if err != nil {
			return nil, fmt.Errorf("resolving policy in %s: %w", path, err)
		}

		for i := range rf.Rules {
			rf.Rules[i].Policy, err = resolvePolicy(dir, rf.Rules[i].Policy)
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
