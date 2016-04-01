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

func start(module string) {
	fmt.Printf("\n [crowley-pack] -> %s\n\n", module)
}

func exit(err error, exit int) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exit)
	}
}

func run(dck docker.Docker, cnf *configuration.Configuration) error {

	ls := docker.NewLogStream()

	start("install")
	if err := install.Install(dck, ls, cnf); err != nil {
		return err
	}

	start("compose")
	if err := compose.Compose(dck, ls, cnf); err != nil {
		return err
	}

	start("publish")
	if err := publish.Publish(dck, ls, cnf); err != nil {
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

	if err := app.Run(os.Args); err != nil {
		exit(err, 1)
	}

}
