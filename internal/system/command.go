package system

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) error
}

type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s %s failed: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return nil
}

type DryRunner struct {
	Commands []string
}

func (r *DryRunner) Run(_ context.Context, name string, args ...string) error {
	r.Commands = append(r.Commands, strings.TrimSpace(name+" "+strings.Join(args, " ")))
	return nil
}
