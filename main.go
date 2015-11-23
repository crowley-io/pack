package main

import (
	"fmt"
	"os"

	"github.com/crowley-io/pack/compose"
	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
	"github.com/crowley-io/pack/install"
	"github.com/crowley-io/pack/publish"
	cli "github.com/jawher/mow.cli"
)

func exit(err error, exit int) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exit)
	}
}

func run(dck docker.Docker, cnf *configuration.Configuration) error {

	if err := install.Install(dck, cnf); err != nil {
		return err
	}

	if err := compose.Compose(dck, cnf); err != nil {
		return err
	}

	if err := publish.Publish(dck, cnf); err != nil {
		return err
	}

	return nil
}

func main() {

	app := cli.App("crowley-pack", "Docker build system.")

	path := getPathOption(app)

	app.Action = func() {

		c, err := configuration.Parse(*path)
		exit(err, 253)

		d, err := docker.New(c)
		exit(err, 254)

		if err = run(d, c); err != nil {
			exit(err, 255)
		}

	}

	app.Run(os.Args)

}
