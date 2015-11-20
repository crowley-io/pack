package main

import (
	"fmt"
	"os"

	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
	"github.com/crowley-io/pack/install"
	cli "github.com/jawher/mow.cli"
)

func exit(err error, exit int) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exit)
	}
}

func main() {

	app := cli.App("crowley-pack", "Docker build system.")

	path := getPathOption(app)

	app.Action = func() {

		c, err := configuration.Parse(*path)
		exit(err, 253)

		d, err := docker.New(c)
		exit(err, 254)

		if err = install.Install(d, c); err != nil {
			exit(err, 1)
		}

		// TODO run pack build step.

	}

	app.Run(os.Args)

}
