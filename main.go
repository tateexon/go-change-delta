package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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

	config := SetConfig(branch, projectPath, excludes, levels, includeTestDeps)

	goList, gitDiff, gitModDiff := MakeExecCalls(config)

	Run(config, goList, gitDiff, gitModDiff)
}

func SetConfig(branch, projectPath, excludes *string, levels *int, includeTestDeps *bool) *Config {
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

func Run(config *Config, goList, gitDiff, gitModDiff *cmd.Output) {
	packages, err := golang.ParsePackages(goList.Stdout)
	if err != nil {
		log.Fatalf("Error parsing packages: %v", err)
	}

	fileGraph := golang.GetGoFileMap(packages, config.IncludeTestDeps)

	var changedPackages []string
	changedPackages, err = git.GetChangedGoPackagesFromDiff(gitDiff.Stdout, config.ProjectPath, config.Excludes, fileGraph)
	if err != nil {
		log.Fatalf("Error getting changed packages: %v", err)
	}

	modChangedPackages, err := git.GetGoModChangesFromDiff(gitModDiff.Stdout)
	if err != nil {
		log.Fatalf("Error getting go.mod changes: %v", err)
	}

	depGraph := golang.GetGoDepMap(packages)

	// Find affected packages
	affectedPkgs := map[string]bool{}
	for _, pkg := range changedPackages {
		p := golang.FindAffectedPackages(pkg, depGraph, false, config.Levels)
		for _, p := range p {
			affectedPkgs[p] = true
		}
	}

	for _, pkg := range modChangedPackages {
		p := golang.FindAffectedPackages(pkg, depGraph, true, config.Levels)
		for _, p := range p {
			affectedPkgs[p] = true
		}
	}

	o := ""
	for k := range affectedPkgs {
		o = fmt.Sprintf("%s %s", o, k)
	}

	if len(o) > 0 {
		fmt.Println(o)
	}
}

func MakeExecCalls(config *Config) (*cmd.Output, *cmd.Output, *cmd.Output) {
	goList, err := golang.GoList()
	if err != nil {
		log.Fatalf("Error getting go list: %v", err)
	}
	gitDiff, err := git.Diff(config.Branch)
	if err != nil {
		log.Fatalf("Error getting the git diff: %v", err)
	}
	gitModDiff, err := git.ModDiff(config.Branch, config.ProjectPath)
	if err != nil {
		log.Fatalf("Error getting the git mod diff")
	}

	return goList, gitDiff, gitModDiff
}
