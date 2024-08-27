package process

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLookPath(t *testing.T) {
	cases := []struct {
		i string
		o bool
	}{
		{"go", true},
		{"git", true},
		{"wontfound", false},
	}

	for _, elem := range cases {
		found := LookPath(elem.i)
		assert.Equal(t, elem.o, found)
	}
}

func TestStartProcess(t *testing.T) {
	cmd := StartProcess("mkdir", []string{"test"}...)

	// is cmd.Process is not nil, then process was started successfully
	assert.NotEqual(t, nil, cmd.Process)
}

func TestRunProcess(t *testing.T) {
	cmd := RunProcess("mkdir", []string{"test"}...)

	// is cmd.Process is not nil, then process was run successfully
	assert.NotEqual(t, nil, cmd.Process)
}

func TestRunProcessWithCommandContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := RunProcessWithContext(ctx, "sleep", []string{"20"}...)
	pid := cmd.Process.Pid
	exist := DoesProcessExist(pid)
	assert.Equal(t, true, exist)
}

func TestStartProcessDecoupled(t *testing.T) {
	StartProcessDecoupled("mkdir", []string{"test"}...)
	fmt.Println("finishing parent process")
}
