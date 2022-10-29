package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func RunCommand(dir string, cmd string, timeoutSeconds int, env map[string]string, args ...string) error {
	var absDir string
	var err error
	if dir == "" {
		dir = "."
	}
	absDir, err = filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("ERROR|runCommand(%s %s)|Cannot get absolute path of '%s'|%s", cmd, args, dir, err.Error())
	}
	err = os.Chdir(absDir)
	if err != nil {
		return fmt.Errorf("ERROR|runCommand(%s %s)|Cannot cd '%s'|%s", cmd, args, dir, err.Error())
	}

	log.Printf("INFO|command.Execute()|Executing '%s %s' in '%s'", cmd, args, absDir)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, cmd, args...)
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Env = os.Environ()
	err = command.Run()
	if err != nil {
		return fmt.Errorf("ERROR|runCommand(%s %s)|Could not execute|%s", cmd, args, err.Error())
	}

	return nil
}
