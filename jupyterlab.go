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
	token    string = ""
	tag      string = "latest"
	args     string = ""
)

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
		params = append(params, "-d")
		params = append(params, "--restart", "on-failure")
	} else {
		params = append(params, "-it", "--rm")
	}
	params = append(params, "-p", fmt.Sprintf("%d:%d", port, port))
	workingDir := ctx.Args().First()
	if workingDir == "" {
		workingDir = "."
	}
	params = append(params, "-v", fmt.Sprintf("%s:%s", workingDir, "/home/jupyter/Workspace"))
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
			&cli.StringFlag{
				Name:        "token",
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
