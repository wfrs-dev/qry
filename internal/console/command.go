package console

import (
	"fmt"
	"os"

	"github.com/integrii/flaggy"
)

type Commander interface {
	Setup() *flaggy.Subcommand
	Execute(ctx *Context) error
	IsUsed() bool
	Name() string
}

type Application struct {
	commands []Commander
	name     string
	version  string
}

func NewApplication(name, version string) *Application {
	flaggy.SetName(name)
	flaggy.SetVersion(version)
	return &Application{
		commands: make([]Commander, 0),
	}
}

func (app *Application) Description(desc string) *Application {
	flaggy.SetDescription(desc)

	return app
}

func (app *Application) AddCommand(cmd Commander) *Application {
	app.commands = append(app.commands, cmd)

	return app
}

func (app *Application) Setup() *Application {
	for _, cmd := range app.commands {
		flaggy.AttachSubcommand(cmd.Setup(), 1)
	}

	flaggy.Parse()

	return app
}

func (app *Application) Run(ctx *Context) error {
	for _, cmd := range app.commands {
		if cmd.IsUsed() {
			return cmd.Execute(ctx)
		}
	}

	var msg string
	if len(os.Args) >= 2 {
		msg = fmt.Sprintf("Comando %q no encontrado", os.Args[1])
	} else {
		msg = "No se especificó ningún comando"
	}
	flaggy.ShowHelpAndExit(msg)

	return nil
}
