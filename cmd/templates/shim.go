package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/term"
)

type Shim struct {
	RunName string
	Image   string
	Workdir string
	Env     []string
	Volumes []string
	Ports   []string
	Stdout  io.Writer
	Stderr  io.Writer
}

func (shim *Shim) Exists() bool {
	cmd := shim.docker("inspect", shim.Image)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

func (shim *Shim) Pull() error {
	cmd := shim.docker("pull", shim.Image)
	cmd.Stdout = os.Stdout // make this output visible
	return cmd.Run()
}

func (shim *Shim) Load(file io.Reader) error {
	cmd := shim.docker("load")
	cmd.Stdout = os.Stdout // make this output visible

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	_, err = io.Copy(stdin, file)
	if err != nil {
		return err
	}

	err = stdin.Close()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (shim *Shim) Exec(containerArgs []string) error {

	args, err := shim.assembleRunArgs() // aquí es donde tengo que hacer la modificación
	if err != nil {
		return err
	}

	// args = append([]string{cmd.Path}, args...)
	args = append(args, containerArgs...)

	cmd := shim.docker(args...)
	cmd.Stdout = os.Stdout // make this output visible
	cmd.Stdout = os.Stderr // make this output visible

	return cmd.Run()
}

func (shim *Shim) docker(arg ...string) *exec.Cmd {
	name := os.Getenv("DOCKER")
	if name == "" {
		name = "docker"
	}

	cmd := exec.Command(name, arg...)
	cmd.Stdout = shim.Stdout
	cmd.Stderr = shim.Stderr
	return cmd
}

func (shim *Shim) assembleRunArgs() ([]string, error) {
	args := []string{"run", "--rm"}

	args = append(args, "--name", shim.RunName)

	if term.IsTerminal(int(os.Stdout.Fd())) {
		args = append(args, "-t")
	}

	for _, env := range shim.Env {
		args = append(args, "-e", env)
	}

	for _, volume := range shim.Volumes {
		args = append(args, "-v", volume)
	}

	for _, bind_port := range shim.Ports {
		args = append(args, "-p", bind_port)
	}

	if shim.Workdir != "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		args = append(args, "-w", shim.Workdir)
		args = append(args, "-v", fmt.Sprintf("%s:%s", cwd, shim.Workdir))
	}

	return append(args, shim.Image), nil
}
