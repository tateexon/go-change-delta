package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/tateexon/go-change-delta/cmd"
	"github.com/tateexon/go-change-delta/git"
	"github.com/tateexon/go-change-delta/golang"
)

type Config struct {
	Branch          string
	ProjectPath     string
	Excludes        []string
	Levels          int
	IncludeTestDeps bool
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		stop() // restore default exit behavior
		log.Println("Cancelling... interrupt again to exit")
	}()

	branch := flag.String("b", "", "The base git branch to compare current changes with. Required.")
	projectPath := flag.String("p", "", "The path to the project. Default is the current directory. Useful for subprojects.")
	excludes := flag.String("e", "", "The comma separated list of paths to exclude. Useful for repositories with multiple go projects within.")
	levels := flag.Int("l", 2, "The number of levels of recursion to search for affected packages. Default is 2. 0 is unlimited.")
	includeTestDeps := flag.Bool("t", true, "Should we include test dependencies. Default is true")
	flag.Parse()

	checkTools()

	config := setConfig(branch, projectPath, excludes, levels, includeTestDeps)

	goList, gitDiff, gitModDiff := makeExecCalls(config)

	run(config, goList, gitDiff, gitModDiff)
}

func setConfig(branch, projectPath, excludes *string, levels *int, includeTestDeps *bool) *Config {
	if *branch == "" {
		log.Fatalf("Branch is required")
	}

	parsedExcludes := []string{}
	if *excludes != "" {
		parsedExcludes = strings.Split(*excludes, ",")
		for i, e := range parsedExcludes {
			parsedExcludes[i] = strings.TrimSpace(e)
		}
	}
	return &Config{
		Branch:          *branch,
		ProjectPath:     *projectPath,
		Excludes:        parsedExcludes,
		Levels:          *levels,
		IncludeTestDeps: *includeTestDeps,
	}
}

// checkTools check if the required tools are installed in the path
func checkTools() {
	_, goErr := exec.LookPath("go")
	_, gitErr := exec.LookPath("git")
	toolsMissing := []string{}
	if goErr != nil {
		toolsMissing = append(toolsMissing, "go")
	}
	if gitErr != nil {
		toolsMissing = append(toolsMissing, "git")
	}
	if len(toolsMissing) > 0 {
		log.Fatalf("tools are missing from your path and may not be installed that are required to run this tool: %+v", toolsMissing)
	}
}

func run(config *Config, goList, gitDiff, gitModDiff *cmd.Output) {
	packages, err := golang.ParsePackages(goList.Stdout)
	if err != nil {
		log.Fatalf("Error parsing packages: %v", err)
	}

	fileMap := golang.GetGoFileMap(packages, config.IncludeTestDeps)

	var changedPackages []string
	changedPackages, err = git.GetChangedGoPackagesFromDiff(gitDiff.Stdout, config.ProjectPath, config.Excludes, fileMap)
	if err != nil {
		log.Fatalf("Error getting changed packages: %v", err)
	}

	changedModPackages, err := git.GetGoModChangesFromDiff(gitModDiff.Stdout)
	if err != nil {
		log.Fatalf("Error getting go.mod changes: %v", err)
	}

	depMap := golang.GetGoDepMap(packages)

	affectedPkgs := findAllAffectedPackages(config, changedPackages, changedModPackages, depMap)

	printAffectedPackages(affectedPkgs)
}

func findAllAffectedPackages(config *Config, changedPackages, changedModPackages []string, depMap golang.DepMap) []string {
	// Find affected packages
	// use map to make handling duplicates simpler
	affectedPkgs := map[string]bool{}

	// loop through packages changed via file changes
	for _, pkg := range changedPackages {
		p := golang.FindAffectedPackages(pkg, depMap, false, config.Levels)
		for _, p := range p {
			affectedPkgs[p] = true
		}
	}

	// loop through packages changed via go.mod changes
	for _, pkg := range changedModPackages {
		p := golang.FindAffectedPackages(pkg, depMap, true, config.Levels)
		for _, p := range p {
			affectedPkgs[p] = true
		}
	}

	// convert map to array
	pkgs := []string{}
	for k := range affectedPkgs {
		pkgs = append(pkgs, k)
	}

	return pkgs
}

func makeExecCalls(config *Config) (*cmd.Output, *cmd.Output, *cmd.Output) {
	goList, err := golang.GoList()
	if err != nil {
		log.Fatalf("Error getting go list: %v\nStdErr: %s", err, goList.Stderr.String())
	}
	gitDiff, err := git.Diff(config.Branch)
	if err != nil {
		log.Fatalf("Error getting the git diff: %v\nStdErr: %s", err, gitDiff.Stderr.String())
	}
	gitModDiff, err := git.ModDiff(config.Branch, config.ProjectPath)
	if err != nil {
		log.Fatalf("Error getting the git mod diff: %v\nStdErr: %s", err, gitModDiff.Stderr.String())
	}

	return goList, gitDiff, gitModDiff
}

func printAffectedPackages(pkgs []string) {
	o := ""
	for _, k := range pkgs {
		o = fmt.Sprintf("%s %s", o, k)
	}

	if len(o) > 0 {
		fmt.Println(o)
	}
}
