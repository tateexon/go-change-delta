package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tateexon/go-change-delta/golang"
	"github.com/tateexon/go-change-delta/utils"
)

func TestSetConfig(t *testing.T) {
	config := setConfig(utils.Ptr("main"), utils.Ptr("abc"), utils.Ptr("abc,123"), utils.Ptr(1), utils.Ptr(true))
	require.NotNil(t, config)
	require.Equal(t, "main", config.Branch)
	require.Equal(t, "abc", config.ProjectPath)
	require.Equal(t, 2, len(config.Excludes))
	require.Equal(t, "abc", config.Excludes[0])
	require.Equal(t, "123", config.Excludes[1])
	require.Equal(t, 1, config.Levels)
	require.True(t, config.IncludeTestDeps)
}

func TestFindAllAffectedPackages(t *testing.T) {
	config := &Config{
		Branch:          "main",
		ProjectPath:     "",
		Excludes:        []string{},
		Levels:          0,
		IncludeTestDeps: true,
	}
	packages := []golang.Package{
		{
			Dir:        "/User/t/git/go-change-delta/test",
			ImportPath: "github.com/tateexon/go-change-delta/test",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"bytes",
			},
			TestImports: []string{
				"unicode",
			},
			GoFiles: []string{
				"test.go",
			},
			TestGoFiles: []string{
				"cmd_test.go",
			},
			EmbedFiles: []string{
				"testdata/blarg.json",
			},
		},
		{
			Dir:        "/User/t/git/go-change-delta/two",
			ImportPath: "github.com/tateexon/go-change-delta/two",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"github.com/tateexon/go-change-delta/test",
			},
			TestImports: []string{
				"fmt",
			},
			GoFiles: []string{
				"test.go",
			},
			TestGoFiles: []string{
				"cmd_test.go",
			},
			EmbedFiles: []string{
				"testdata/blarg.json",
			},
		},
	}
	depMap := golang.GetGoDepMap(packages)

	t.Run("only changed files", func(t *testing.T) {
		changedPackages := []string{
			"github.com/tateexon/go-change-delta/test",
		}
		changedModPackages := []string{}

		pkgs := findAllAffectedPackages(config, changedPackages, changedModPackages, depMap)
		require.Equal(t, 2, len(pkgs))
		require.Equal(t, "github.com/tateexon/go-change-delta/test", pkgs[0])
		require.Equal(t, "github.com/tateexon/go-change-delta/two", pkgs[1])
	})

	t.Run("only changed go mod file", func(t *testing.T) {
		changedPackages := []string{}
		changedModPackages := []string{
			"fmt",
		}

		pkgs := findAllAffectedPackages(config, changedPackages, changedModPackages, depMap)
		require.Equal(t, 1, len(pkgs))
		require.Equal(t, "github.com/tateexon/go-change-delta/two", pkgs[0])
	})
}
