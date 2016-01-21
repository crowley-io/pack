package main

import (
	cli "github.com/jawher/mow.cli"
)

func getPathOption(app *cli.Cli) (path *string) {

	path = app.String(cli.StringOpt{
		Name:  "f file",
		Value: "packer.yml",
		Desc:  "Configuration file",
	})

	return
}
