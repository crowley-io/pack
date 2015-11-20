package main

import (
	"fmt"
	"os"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
	"github.com/crowley-io/pack/install"
	cli "github.com/jawher/mow.cli"
)

func main() {

	app := cli.App("crowley-pack", "Docker build system.")

	path := getPathOption(app)

	app.Action = func() {

		c, err := configuration.Parse(*path)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(255)
			return
		}

		d, err := docker.New(c)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(255)
			return
		}

		if err = install.Install(d, c); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// TODO run pack build step.

	}

	app.Run(os.Args)

}
