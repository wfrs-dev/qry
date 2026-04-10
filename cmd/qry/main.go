package main

import (
	"fmt"
	"os"

	"gitlab.com/wfrsgo/qry/internal/command"
	"gitlab.com/wfrsgo/qry/internal/console"
)

func main() {
	app := console.NewApplication("qry", "0.1.0")
	app.Description("Cliente mímimo para consultar bases de datos")

	app.AddCommand(&command.AddCommand{})
	app.AddCommand(&command.ListCommand{})
	app.AddCommand(&command.RunCommand{})
	app.Setup()

	ctx := console.GetContext()
	if err := app.Run(ctx); err != nil {
		fmt.Println(console.Colorize("<italic+red> %s</>", err))
		os.Exit(99)
	}
}
