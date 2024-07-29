package git

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tateexon/go-change-delta/utils"
)

func verifySliceItemsPresent(t *testing.T, expected, actual []string) {
	require.Equal(t, len(expected), len(actual), fmt.Sprintf("Expected: %+v\nActual: %+v\n", expected, actual))
	for _, item := range expected {
		require.Contains(t, actual, item)
	}
}

func TestGetGoModChangesFromDiff(t *testing.T) {
	input, err := utils.FileToBytesBuffer("../testdata/gitmoddiff.txt")
	require.NoError(t, err, "Failed to load diff file")

	changes, err := GetGoModChangesFromDiff(input)
	require.NoError(t, err, "Failed to get diff changes")

	expected := []string{
		"github.com/tateexon/go-change-delta",
	}
	verifySliceItemsPresent(t, expected, changes)
}

func TestGetChangedGoPackagesFromDiff(t *testing.T) {
	input, err := utils.FileToBytesBuffer("../testdata/gitdiff.txt")
	require.NoError(t, err, "Failed to load diff file")

	fileGraph := map[string]string{
		"cmd/cmd.go":      "abc",
		"git/git.go":      "def",
		"git/git_test.go": "efg",
		"golang/mod.go":   "hij",
		"utils/utils.go":  "klm",
	}
	packages, err := GetChangedGoPackagesFromDiff(input, "", []string{"utils"}, fileGraph)
	require.NoError(t, err, "Failed to get diff changes")

	expected := []string{
		"abc",
		"def",
		"efg",
		"hij",
	}
	verifySliceItemsPresent(t, expected, packages)
}
