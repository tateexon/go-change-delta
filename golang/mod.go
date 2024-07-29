package golang

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tateexon/go-change-delta/cmd"
)

type Package struct {
	Dir          string   `json:"Dir"`
	ImportPath   string   `json:"ImportPath"`
	Root         string   `json:"Root"`
	Deps         []string `json:"Deps"`
	TestImports  []string `json:"TestImports"`
	XTestImports []string `json:"XTestImports"`
	GoFiles      []string `json:"GoFiles"`
	TestGoFiles  []string `json:"TestGoFiles"`
	XTestGoFiles []string `json:"XTestGoFiles"`
	EmbedFiles   []string `json:"EmbedFiles"`
}

type DepGraphItem struct {
	ImportPath string
	Root       string
	GoFiles    []string
}

func GoList() (*cmd.Output, error) {
	return cmd.Execute("go", "list", "-json", "./...")
}

func ParsePackages(goList bytes.Buffer) ([]Package, error) {
	var packages []Package
	scanner := bufio.NewScanner(&goList)
	var buffer bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		buffer.WriteString(line)

		if line == "}" {
			var pkg Package
			if err := json.Unmarshal(buffer.Bytes(), &pkg); err != nil {
				return nil, err
			}
			packages = append(packages, pkg)
			buffer.Reset()
		}
	}

	err := scanner.Err()
	return packages, err
}

func GetGoDepMap(packages []Package) map[string][]string {
	depGraph := make(map[string][]string)
	for _, pkg := range packages {
		for _, dep := range pkg.Deps {
			depGraph[dep] = append(depGraph[dep], pkg.ImportPath)
		}
		for _, dep := range pkg.TestImports {
			depGraph[dep] = append(depGraph[dep], pkg.ImportPath)
		}
		for _, dep := range pkg.XTestImports {
			depGraph[dep] = append(depGraph[dep], pkg.ImportPath)
		}
	}
	return depGraph
}

//nolint:revive
func GetGoFileMap(packages []Package, includeTestFiles bool) map[string]string {
	// Build dependency graph
	fileGraph := make(map[string]string)
	for _, pkg := range packages {
		addToGraph(pkg, pkg.GoFiles, fileGraph)
		addToGraph(pkg, pkg.EmbedFiles, fileGraph)
		if includeTestFiles {
			addToGraph(pkg, pkg.TestGoFiles, fileGraph)
			addToGraph(pkg, pkg.XTestGoFiles, fileGraph)
		}

	}
	return fileGraph
}

func addToGraph(pkg Package, files []string, fileGraph map[string]string) {
	for _, file := range files {
		path := strings.Replace(pkg.Dir, fmt.Sprintf("%s/", pkg.Root), "", 1)
		key := fmt.Sprintf("%s/%s", path, file)
		if _, exists := fileGraph[key]; exists {
			log.Fatalf("why did this happen, duplicate key %s\nfile a bug\n", key)
		}
		fileGraph[key] = pkg.ImportPath
	}
}

//nolint:revive
func FindAffectedPackages(pkg string, depGraph map[string][]string, externalPackage bool, maxDepth int) []string {
	visited := make(map[string]bool)
	var affected []string

	var dfs func(string, int)
	dfs = func(p string, depthLeft int) {
		if visited[p] {
			return
		}

		visited[p] = true
		// exclude the package itself if it is an external package
		if !(externalPackage && p == pkg) {
			affected = append(affected, p)
		}
		d := depthLeft - 1
		if d != 0 {
			for _, dep := range depGraph[p] {
				dfs(dep, d)
			}
		}
	}

	depth := maxDepth
	// depth is zero then we want infinite recursion, set this to -1 to enable this
	if depth <= 0 {
		depth = -1
	}
	dfs(pkg, depth)
	return affected
}
