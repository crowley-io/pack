package main

import (
	"fmt"
	"os"

	"github.com/crowley-io/pack/compose"
	"github.com/crowley-io/pack/configuration"
	"github.com/crowley-io/pack/docker"
	"github.com/crowley-io/pack/install"
	"github.com/crowley-io/pack/publish"
	"github.com/fatih/color"
	cli "github.com/jawher/mow.cli"
)

var (
	// Version export pack binary version
	Version string
)

func start(module string) {
	color.New(color.Bold).PrintfFunc()("\n [crowley-pack] -> %s\n\n", module)
}

func exit(err error, exit int) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		cli.Exit(exit)
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

func handle(parse func() (*configuration.Configuration, error)) {

	c, err := parse()
	exit(err, 253)

	d, err := docker.New(c)
	exit(err, 254)

	if err = run(d, c); err != nil {
		exit(err, 255)
	}

}

func main() {

	app := cli.App("crowley-pack", "Docker build system.")
	app.Version("v version", fmt.Sprintf("crowley-pack %s", Version))

	app.Command("push", "Compose an image and publish on a registry", func(cmd *cli.Cmd) {
		remote, nocache, nopull := getPushOptions(cmd)
		cmd.Action = func() {
			handle(createRemoteConfigurationParser(remote, nocache, nopull))
		}
	})

	path := getPathOption(app)

	app.Action = func() {
		handle(createFileConfigurationParser(path))
	}

	if err := app.Run(os.Args); err != nil {
		exit(err, 1)
	}

}
