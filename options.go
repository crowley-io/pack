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

func getPushOptions(app *cli.Cmd) (remote *string, nocache, nopull *bool) {

	remote = app.String(cli.StringArg{
		Name:      "REMOTE",
		Value:     "",
		Desc:      "Image's remote identifier",
		HideValue: true,
	})

	nocache = app.Bool(cli.BoolOpt{
		Name:      "no-cache",
		Value:     false,
		Desc:      "Disable cache when composing the image",
		HideValue: true,
	})

	nopull = app.Bool(cli.BoolOpt{
		Name:      "no-pull",
		Value:     false,
		Desc:      "Disable pull when composing the image",
		HideValue: true,
	})

	return
}
