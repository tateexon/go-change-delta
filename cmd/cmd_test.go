package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecuteBadCommand(t *testing.T) {
	_, err := Execute("zznotaprogram/\n")
	require.Error(t, err)
}

func TestExecuteGoodCommand(t *testing.T) {
	c, err := Execute("which", "which")
	require.NoError(t, err)
	require.Contains(t, c.Stdout.String(), "which")
	require.Equal(t, c.Stderr.String(), "")
}
