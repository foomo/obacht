package runner

import (
	"bufio"
	"fmt"
	"io/fs"
	"strings"
)

const libPrefix = "_lib/"

// Preprocess expands `# include: <path>` directives at the top of a rule
// script into inlined library content from fsys. scriptPath is the path of
// the rule script being processed (used for error messages only).
//
// Rules:
//   - includes must appear before the first non-comment, non-shebang, non-blank line
//   - include path must be under `_lib/` (no traversal)
//   - includes are deduplicated by path
//   - libraries may include other libraries (DAG; cycles are rejected)
//   - the resolved output preserves declared order: lib bodies first, then rule body
//   - include directive lines are stripped from the resulting body
func Preprocess(fsys fs.FS, scriptPath, script string) (string, error) {
	includes, body, err := parseIncludes(script)
	if err != nil {
		return "", fmt.Errorf("parsing includes in %s: %w", scriptPath, err)
	}

	if len(includes) == 0 {
		return script, nil
	}

	visited := map[string]bool{}
	stack := map[string]bool{}

	var libs []string

	for _, inc := range includes {
		expanded, err := expandInclude(fsys, inc, visited, stack)
		if err != nil {
			return "", err
		}

		libs = append(libs, expanded...)
	}

	var b strings.Builder

	b.WriteString(body)

	out := stitch(libs, b.String())

	return out, nil
}

// parseIncludes scans the top of the script for `# include: <path>` lines.
// Returns the list of paths (in order) and the script body with directive
// lines removed. The shebang line is preserved at the top.
func parseIncludes(script string) ([]string, string, error) {
	var (
		includes  []string
		bodyLines []string
		inHeader  = true
		shebang   string
	)

	scanner := bufio.NewScanner(strings.NewReader(script))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()

		if inHeader {
			trimmed := strings.TrimSpace(line)

			switch {
			case strings.HasPrefix(line, "#!") && shebang == "":
				shebang = line
				continue
			case trimmed == "":
				bodyLines = append(bodyLines, line)
				continue
			case strings.HasPrefix(line, "# include:"):
				path := strings.TrimSpace(strings.TrimPrefix(line, "# include:"))
				if path == "" {
					return nil, "", fmt.Errorf("empty include directive")
				}

				includes = append(includes, path)

				continue
			case strings.HasPrefix(trimmed, "#"):
				bodyLines = append(bodyLines, line)
				continue
			default:
				inHeader = false
			}
		}

		bodyLines = append(bodyLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, "", err
	}

	body := strings.Join(bodyLines, "\n")
	if shebang != "" {
		body = shebang + "\n" + body
	}

	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}

	return includes, body, nil
}

// expandInclude resolves a single include path recursively, returning the
// lib bodies in declared order. visited tracks already-included paths
// (for dedup). stack tracks the current resolution path (for cycle detection).
func expandInclude(fsys fs.FS, path string, visited, stack map[string]bool) ([]string, error) {
	if !strings.HasPrefix(path, libPrefix) {
		return nil, fmt.Errorf("include path must be under %s: %q", libPrefix, path)
	}

	if strings.Contains(path, "..") {
		return nil, fmt.Errorf("include path must not contain '..': %q", path)
	}

	if stack[path] {
		return nil, fmt.Errorf("include cycle detected at %q", path)
	}

	if visited[path] {
		return nil, nil
	}

	fullPath := "inputs/" + path

	data, err := fs.ReadFile(fsys, fullPath)
	if err != nil {
		return nil, fmt.Errorf("reading include %q: %w", path, err)
	}

	stack[path] = true

	subIncludes, body, err := parseIncludes(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing include %q: %w", path, err)
	}

	var out []string

	for _, sub := range subIncludes {
		expanded, err := expandInclude(fsys, sub, visited, stack)
		if err != nil {
			return nil, err
		}

		out = append(out, expanded...)
	}

	out = append(out, stripShebang(body))

	delete(stack, path)

	visited[path] = true

	return out, nil
}

// stripShebang removes a leading #! line from a library body.
func stripShebang(body string) string {
	if strings.HasPrefix(body, "#!") {
		_, after, ok := strings.Cut(body, "\n")
		if ok {
			return after
		}
	}

	return body
}

// stitch concatenates lib bodies in order, then appends the rule body.
// The rule body's shebang stays on top.
func stitch(libs []string, ruleBody string) string {
	var b strings.Builder

	shebang := ""
	rest := ruleBody

	if strings.HasPrefix(ruleBody, "#!") {
		i := strings.IndexByte(ruleBody, '\n')
		if i >= 0 {
			shebang = ruleBody[:i+1]
			rest = ruleBody[i+1:]
		}
	}

	if shebang != "" {
		b.WriteString(shebang)
	}

	for _, lib := range libs {
		b.WriteString(lib)

		if !strings.HasSuffix(lib, "\n") {
			b.WriteString("\n")
		}
	}

	b.WriteString(rest)

	return b.String()
}
