package process

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func LookPath(name string) bool {
	_, err := exec.LookPath(name)
	if err != nil && errors.Is(err, exec.ErrNotFound) {
		return false
	}
	if err != nil {
		log.Fatalln(err)
	}
	return true
}

func StartProcess(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	// Start does not wait for process to complete
	// Is State() was successful, cmd.Process will be set
	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	// Start command needs to be followed up with Wait command
	// Wait will wait for comaand to finish and then realeases sytem recources
	err = cmd.Wait()
	if err != nil {
		log.Fatalln(err)
	}

	return cmd
}

func RunProcess(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	return cmd
}

func RunProcessWithContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	// This is done to clean up the child process
	// and to make cleanup non blocking so that
	// pid can be extracted
	go func(cmd *exec.Cmd) {
		err := cmd.Wait()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("child process cleaned up")
	}(cmd)

	return cmd
}

func DoesProcessExist(pid int) bool {
	p, err := FindProcess(pid)
	if err != nil {
		return false
	}

	// On unix find process always succeeds
	// The recommended way is to send a signal to the process and check if any error is reported
	// If signal is zero, no signal is sent but error checking is still performed
	err = p.Signal(syscall.Signal(0))
	return err == nil
}

func FindProcess(pid int) (*os.Process, error) {
	return os.FindProcess(pid)
}

// StartProcessDecoupled will start a child process such that when parent process is killed, child process still runs
func StartProcessDecoupled(name string, args ...string) {
	cmd := exec.Command(name, args...)
	err := cmd.Start()
	if err != nil {
		log.Fatalln(err)
	}

	// this will move child process to a different process group
	// Note: Typing ^C on the terminal kills all processes in the process group
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// start a go routine to wait for the child process
	go func(cmd *exec.Cmd) {
		time.Sleep(10 * time.Second)
		err := cmd.Wait()
		if err != nil {
			log.Fatalln(err)
		}
		err = os.Chdir("test")
		if err != nil {
			log.Fatalln(err)
		}
		// [TODO] - goroutine is not adding logs to the file - debug
		data := []byte("succesfully closed child process")
		err = os.WriteFile("logs.txt", data, 0o644)
		if err != nil {
			log.Fatalln(err)
		}
	}(cmd)
}
