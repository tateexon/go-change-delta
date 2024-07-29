package cmd

import (
	"bytes"
	"os/exec"
)

type Output struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func Execute(name string, args ...string) (*Output, error) {
	cmd := exec.Command(name, args...)
	out := &Output{}
	cmd.Stdout = &out.Stdout
	cmd.Stderr = &out.Stderr
	err := cmd.Run()
	return out, err
}
