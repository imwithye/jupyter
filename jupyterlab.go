package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	DOCKER_IMAGE = "imwithye/jupyterlab"
)

var (
	pull     bool
	port     int
	detached bool
	token    string
	tag      string
)

func DockerPullCmd() error {
	cmd := exec.Command("docker", "pull", fmt.Sprintf("%s:%s", DOCKER_IMAGE, tag))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DockerRunCmd() error {
	args := []string{"run"}
	args = append(args, "-p", fmt.Sprintf("%d:%d", port, port))
	if detached {
		args = append(args, "-d")
	}
	args = append(args, "-e", fmt.Sprintf("JUPYTERLAB_PORT=%d", port))
	args = append(args, "-e", fmt.Sprintf("JUPYTERLAB_TOKEN=%s", token))
	args = append(args, fmt.Sprintf("%s:%s", DOCKER_IMAGE, tag))
	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Run(ctx *cli.Context) error {
	if pull {
		err := DockerPullCmd()
		if err != nil {
			return err
		}
	}
	return DockerRunCmd()
}

func main() {
	err := (&cli.App{
		Name:  filepath.Base(os.Args[0]),
		Usage: "Run JupyterLab in the container.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "pull",
				Value:       false,
				Usage:       "pull the docker image",
				Destination: &pull,
			},
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       8888,
				Usage:       "port to expose",
				Destination: &port,
			},
			&cli.BoolFlag{
				Name:        "detached",
				Aliases:     []string{"d"},
				Value:       false,
				Usage:       "run in detached mode",
				Destination: &detached,
			},
			&cli.StringFlag{
				Name:        "token",
				Usage:       "jupyterlab token",
				Destination: &token,
				Action: func(c *cli.Context, t string) error {
					if t != strings.ToLower(t) {
						return fmt.Errorf("token must be lowercase")
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "tag",
				Value:       "latest",
				Usage:       "docker image tag",
				Destination: &tag,
			},
		},
		Action: Run,
	}).Run(os.Args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
