package collector

import "path/filepath"

// resolvePath returns the resolved path following symlinks.
// If resolution fails (e.g. path doesn't exist), it returns the original path.
func resolvePath(path string) string {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}

	return resolved
}
