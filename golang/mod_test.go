package golang

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tateexon/go-change-delta/utils"
)

func TestParsePackages(t *testing.T) {
	input, err := utils.FileToString("../testdata/golist.txt")
	require.NoError(t, err, "Failed to load diff file")

	packages, err := ParsePackages(utils.StringToBytesBuffer(input))
	require.NoError(t, err)
	require.Equal(t, 5, len(packages), "Incorrect number of packages found. Packages found: %+v", packages)
	require.Equal(t, "github.com/tateexon/go-change-delta", packages[0].ImportPath)
	require.Equal(t, "github.com/tateexon/go-change-delta/utils", packages[4].ImportPath)
}

func TestParseInvalidPackages(t *testing.T) {
	_, err := ParsePackages(utils.StringToBytesBuffer("abc\n}"))
	require.Error(t, err)
}

func TestGetGoDepMap(t *testing.T) {
	packages := []Package{
		{
			Dir:        "./test",
			ImportPath: "github.com/tateexon/go-change-delta",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"bytes",
				"cmp",
			},
			TestImports: []string{
				"unicode",
				"unicode/utf16",
				"unicode/utf8",
				"unsafe",
			},
			XTestImports: []string{
				"fmt",
				"github.com/stretchr/testify/require",
				"github.com/tateexon/go-change-delta/utils",
				"testing",
			},
		},
		{
			Dir:        "./test",
			ImportPath: "github.com/tateexon/go-change-delta/utils",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"bytes",
			},
		},
	}
	depMap := GetGoDepMap(packages)
	require.Equal(t, 10, len(depMap))
	require.Equal(t, 2, len(depMap["bytes"]))
	require.Equal(t, "github.com/tateexon/go-change-delta", depMap["bytes"][0])
	require.Equal(t, "github.com/tateexon/go-change-delta/utils", depMap["bytes"][1])
}

func TestGetGoFileMap(t *testing.T) {
	packages := []Package{
		{
			Dir:        "C:\\User\\t\\git\\go-change-delta\\test",
			ImportPath: "github.com/tateexon/go-change-delta/test",
			Root:       "C:\\User\\t\\git\\go-change-delta",
			Deps: []string{
				"bytes",
				"cmp",
			},

			GoFiles: []string{
				"cmd.go",
			},
			TestGoFiles: []string{
				"cmd_test.go",
			},
			XTestGoFiles: []string{
				"golang_test.go",
			},
		},
		{
			Dir:        "/User/t/git/go-change-delta/utils",
			ImportPath: "github.com/tateexon/go-change-delta/utils",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"bytes",
			},
			GoFiles: []string{
				"cmd.go",
			},
			EmbedFiles: []string{
				"child/embed.json",
			},
		},
		{
			Dir:        "/User/t/git/go-change-delta/utils/child",
			ImportPath: "github.com/tateexon/go-change-delta/utils/child",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"arg",
			},
			GoFiles: []string{
				"child.go",
			},
			EmbedFiles: []string{
				"embed.json",
			},
		},
	}

	t.Run("include test files", func(t *testing.T) {
		fileMap := GetGoFileMap(packages, true)
		require.Equal(t, 6, len(fileMap))
		require.Equal(t, "github.com/tateexon/go-change-delta/test", fileMap["test/cmd.go"][0], fmt.Sprintf("%+v", fileMap))
		require.Equal(t, "github.com/tateexon/go-change-delta/utils", fileMap["utils/cmd.go"][0], fmt.Sprintf("%+v", fileMap))
		require.Equal(t, "github.com/tateexon/go-change-delta/utils", fileMap["utils/child/embed.json"][0])
		require.Equal(t, "github.com/tateexon/go-change-delta/utils/child", fileMap["utils/child/embed.json"][1])
	})
	t.Run("exclude test files", func(t *testing.T) {
		fileMap := GetGoFileMap(packages, false)
		require.Equal(t, 4, len(fileMap))
		require.Equal(t, "github.com/tateexon/go-change-delta/test", fileMap["test/cmd.go"][0], fmt.Sprintf("%+v", fileMap))
		require.Equal(t, "github.com/tateexon/go-change-delta/utils", fileMap["utils/cmd.go"][0], fmt.Sprintf("%+v", fileMap))
	})
}

func TestFindAffectedPackages(t *testing.T) {
	packages := []Package{
		{
			Dir:        "/User/t/git/go-change-delta/test",
			ImportPath: "github.com/tateexon/go-change-delta/test",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"bytes",
				"cmp",
			},
			TestImports: []string{
				"unicode",
				"unicode/utf16",
				"unicode/utf8",
				"unsafe",
			},
			XTestImports: []string{
				"fmt",
				"github.com/stretchr/testify/require",
				"github.com/tateexon/go-change-delta/utils",
				"testing",
			},
			GoFiles: []string{
				"test.go",
			},
			TestGoFiles: []string{
				"cmd_test.go",
			},
			XTestGoFiles: []string{
				"cmd_test_test.go",
			},
			EmbedFiles: []string{
				"testdata/blarg.json",
			},
		},
		{
			Dir:        "/User/t/git/go-change-delta/utils",
			ImportPath: "github.com/tateexon/go-change-delta/utils",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"bytes",
			},
			GoFiles: []string{
				"utils.go",
			},
		},
		{
			Dir:        "/User/t/git/go-change-delta/three",
			ImportPath: "github.com/tateexon/go-change-delta/three",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"github.com/tateexon/go-change-delta/test",
			},
			GoFiles: []string{
				"three.go",
			},
		},
		{
			Dir:        "/User/t/git/go-change-delta/four",
			ImportPath: "github.com/tateexon/go-change-delta/four",
			Root:       "/User/t/git/go-change-delta",
			Deps: []string{
				"github.com/tateexon/go-change-delta/three",
			},
			GoFiles: []string{
				"four.go",
			},
		},
	}
	depMap := GetGoDepMap(packages)
	require.Equal(t, 12, len(depMap))

	t.Run("depth 1 should only find self", func(t *testing.T) {
		found := FindAffectedPackages("github.com/tateexon/go-change-delta/utils", depMap, false, 1)
		require.Equal(t, 1, len(found))
		require.Equal(t, "github.com/tateexon/go-change-delta/utils", found[0])
	})

	t.Run("depth 2 should find only children with dependencies", func(t *testing.T) {
		found := FindAffectedPackages("github.com/tateexon/go-change-delta/utils", depMap, false, 2)
		require.Equal(t, 2, len(found))
		require.Equal(t, "github.com/tateexon/go-change-delta/utils", found[0])
		require.Equal(t, "github.com/tateexon/go-change-delta/test", found[1])
	})

	t.Run("go mod lib in both modules", func(t *testing.T) {
		found := FindAffectedPackages("bytes", depMap, true, 2)
		require.Equal(t, 2, len(found))
		require.Equal(t, "github.com/tateexon/go-change-delta/test", found[0])
		require.Equal(t, "github.com/tateexon/go-change-delta/utils", found[1])
	})

	t.Run("go mod lib in one lib", func(t *testing.T) {
		found := FindAffectedPackages("cmp", depMap, true, 2)
		require.Equal(t, 1, len(found))
		require.Equal(t, "github.com/tateexon/go-change-delta/test", found[0])
	})

	t.Run("find affected of affected, recursion level 3", func(t *testing.T) {
		found := FindAffectedPackages("bytes", depMap, true, 3)
		require.Equal(t, 3, len(found))
		require.Equal(t, "github.com/tateexon/go-change-delta/test", found[0])
		require.Equal(t, "github.com/tateexon/go-change-delta/three", found[1])
		require.Equal(t, "github.com/tateexon/go-change-delta/utils", found[2])
	})

	t.Run("infinite recursion test", func(t *testing.T) {
		found := FindAffectedPackages("bytes", depMap, true, -1)
		require.Equal(t, 4, len(found))
		require.Equal(t, "github.com/tateexon/go-change-delta/test", found[0])
		require.Equal(t, "github.com/tateexon/go-change-delta/three", found[1])
		require.Equal(t, "github.com/tateexon/go-change-delta/four", found[2])
		require.Equal(t, "github.com/tateexon/go-change-delta/utils", found[3])
	})
}
