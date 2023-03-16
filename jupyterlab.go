package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	port     int
	detached bool
	token    string
	tag      string
)

func DockerCmd() *exec.Cmd {
	args := []string{"run"}
	args = append(args, "-p", fmt.Sprintf("%d:8888", port))
	if detached {
		args = append(args, "-d")
	}
	if token != "" {
		args = append(args, "-e", fmt.Sprintf("TOKEN=%s", token))
	}
	args = append(args, fmt.Sprintf("imwithye/jupyterlab:%s", tag))
	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func Run(ctx *cli.Context) error {
	cmd := DockerCmd()
	return cmd.Run()
}

func main() {
	err := (&cli.App{
		Name:  "jupyterlab.exe",
		Usage: "Run JupyterLab in the container.",
		Flags: []cli.Flag{
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
