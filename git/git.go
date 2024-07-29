package git

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/tateexon/go-change-delta/cmd"
)

func Diff(baseBranch string) (*cmd.Output, error) {
	return cmd.Execute("git", "diff", "--name-only", baseBranch)
}

func ModDiff(baseBranch, projectPath string) (*cmd.Output, error) {
	return cmd.Execute("git", "diff", baseBranch, "--unified=0", "--", filepath.Join(projectPath, "go.mod"))
}

func GetGoModChangesFromDiff(lines bytes.Buffer) ([]string, error) {
	changedLines := strings.Split(lines.String(), "\n")

	// Filter out lines that do not indicate package changes
	var packages []string
	for _, line := range changedLines {
		if strings.HasPrefix(line, "+") {
			// ignore comments or empty lines (e.g., not relevant)
			if strings.HasPrefix(line, "+ ") || strings.HasPrefix(line, "+++ ") {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) > 1 {
				// The second field should contains the module path
				packages = append(packages, fields[1])
			}
		}
	}

	return packages, nil
}

func GetChangedGoPackagesFromDiff(out bytes.Buffer, projectPath string, excludes []string, fileMap map[string][]string) ([]string, error) {
	changedFiles := strings.Split(out.String(), "\n")

	// Filter out non-Go files and directories and embeds
	changedPackages := make(map[string]bool)
	for _, file := range changedFiles {
		if file == "" || !strings.HasPrefix(file, projectPath) || shouldExclude(excludes, file) {
			continue
		}

		// if the changed file is in the fileMap then we add it to the changed packages
		for _, importPath := range fileMap[file] {
			changedPackages[importPath] = true
		}
	}

	// Convert map keys to slice
	var packages []string
	for pkg := range changedPackages {
		packages = append(packages, pkg)
	}

	return packages, nil
}

func shouldExclude(excludes []string, item string) bool {
	for _, v := range excludes {
		if strings.HasPrefix(item, v) {
			return true
		}
	}
	return false
}
