package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/urfave/cli/v2"
)

const (
	DOCKER_IMAGE = "imwithye/jupyterlab"
)

var (
	dryrun   bool   = false
	pull     bool   = false
	port     int    = 8888
	detached bool   = false
	gpu      bool   = false
	token    string = ""
	tag      string = "latest"
	args     string = ""
)

func AutoMounts() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	args := []string{}

	for {
		pipConfigDir := filepath.Join(home, ".config", "pip")
		pipConfigDir, err := filepath.Abs(pipConfigDir)
		if err != nil {
			break
		}
		_, err = os.Stat(pipConfigDir)
		if err != nil {
			break
		}
		args = append(args, "-v", fmt.Sprintf("%s:/home/jupyter/.config/pip", pipConfigDir))
		break
	}

	for {
		gitConfig := filepath.Join(home, ".gitconfig")
		gitConfig, err := filepath.Abs(gitConfig)
		if err != nil {
			break
		}
		_, err = os.Stat(gitConfig)
		if err != nil {
			break
		}
		args = append(args, "-v", fmt.Sprintf("%s:/home/jupyter/.gitconfig", gitConfig))
		break
	}

	for {
		gitIgnore := filepath.Join(home, ".gitignore")
		gitIgnore, err := filepath.Abs(gitIgnore)
		if err != nil {
			break
		}
		_, err = os.Stat(gitIgnore)
		if err != nil {
			break
		}
		args = append(args, "-v", fmt.Sprintf("%s:/home/jupyter/.gitignore", gitIgnore))
		break
	}

	return args
}

func DockerPullCmd(ctx *cli.Context) error {
	if dryrun {
		fmt.Println("docker", "pull", fmt.Sprintf("%s:%s", DOCKER_IMAGE, tag))
		return nil
	}

	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", DOCKER_IMAGE, tag))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DockerRunCmd(ctx *cli.Context) error {
	params := []string{"run"}
	if detached {
		params = append(params, "-d", "--restart", "on-failure")
	} else {
		params = append(params, "-it", "--rm")
	}
	if gpu {
		params = append(params, "--gpus", "all")
	}
	params = append(params, "-p", fmt.Sprintf("%d:%d", port, port))
	workingDir := ctx.Args().First()
	if workingDir == "" {
		workingDir = "."
	}
	params = append(params, "-v", fmt.Sprintf("%s:%s", workingDir, "/home/jupyter/Workspace"))
	params = append(params, AutoMounts()...)
	if args != "" {
		if argsSplit, err := shlex.Split(args); err == nil {
			params = append(params, argsSplit...)
		}
	}
	params = append(params, fmt.Sprintf("%s:%s", DOCKER_IMAGE, tag))
	params = append(params, "jupyter", "lab", "--ip=0.0.0.0", fmt.Sprintf("--port=%d", port), "--no-browser")
	if token != "" {
		params = append(params, fmt.Sprintf("--NotebookApp.token=%s", token))
	}
	params = append(params, "--NotebookApp.terminado_settings={'shell_command': ['/usr/bin/bash']}")

	// if dry run
	if dryrun {
		fmt.Println("docker", strings.Join(params, " "))
		return nil
	}

	cmd := exec.Command("docker", params...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Run(ctx *cli.Context) error {
	if pull {
		err := DockerPullCmd(ctx)
		if err != nil {
			return err
		}
	}
	return DockerRunCmd(ctx)
}

func main() {
	err := (&cli.App{
		Name:  filepath.Base(os.Args[0]),
		Usage: "Run JupyterLab in the container.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "dryrun",
				Value:       dryrun,
				Usage:       "dryrun",
				Destination: &dryrun,
			},
			&cli.BoolFlag{
				Name:        "pull",
				Value:       pull,
				Usage:       "pull the docker image",
				Destination: &pull,
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       port,
				Usage:       "port to expose",
				Destination: &port,
			},
			&cli.BoolFlag{
				Name:        "detached",
				Aliases:     []string{"d"},
				Value:       detached,
				Usage:       "run in detached mode",
				Destination: &detached,
			},
			&cli.BoolFlag{
				Name:        "gpu",
				Aliases:     []string{"g"},
				Value:       gpu,
				Usage:       "enable gpu",
				Destination: &gpu,
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				Value:       token,
				Usage:       "jupyterlab token",
				Destination: &token,
			},
			&cli.StringFlag{
				Name:        "tag",
				Value:       tag,
				Usage:       "docker image tag",
				Destination: &tag,
			},
			&cli.StringFlag{
				Name:        "args",
				Value:       args,
				Usage:       "additional docker arguments",
				Destination: &args,
			},
		},
		Action: Run,
	}).Run(os.Args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
