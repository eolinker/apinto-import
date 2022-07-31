package main

import (
	"os"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/apinto-import/cli"
)

func init() {
	InitCLILog()
}

func main() {
	app := cli.NewApp()
	app.Default()
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
		return
	}
}
