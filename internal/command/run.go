package command

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/integrii/flaggy"
	"gitlab.com/wfrsgo/qry/internal/config"
	"gitlab.com/wfrsgo/qry/internal/console"
	"gitlab.com/wfrsgo/qry/internal/db"
	"gitlab.com/wfrsgo/qry/internal/formatter"
)

type RunCommand struct {
	*flaggy.Subcommand
	qry  string
	name string
	list bool
}

func (cmd *RunCommand) Setup() *flaggy.Subcommand {
	cmd.Subcommand = flaggy.NewSubcommand(cmd.Name())
	cmd.Subcommand.Description = "Ejecuta una consulta SQL dada por la opción `query` o por la entrada estándar si no se especifica `query`"
	cmd.Subcommand.String(&cmd.qry, "q", "query", "Consulta a ejecutar")
	cmd.Subcommand.Bool(&cmd.list, "l", "list", "Muestra el resultado en formato lista, por defecto se muestra en formato tabla")
	cmd.Subcommand.AddPositionalValue(&cmd.name, "name", 1, true, "Nombre de la conexión a utilizar")

	return cmd.Subcommand
}

func (cmd *RunCommand) Execute(ctx *console.Context) error {
	cnx, err := config.GetConnection(cmd.name)
	if err != nil {
		return err
	}

	if cmd.qry == "" {
		cmd.qry = fromStdin()
	}

	if cmd.qry == "" {
		return errors.New("no se especificó consulta")
	}

	driver, err := db.GetDriver(cnx.Driver)
	if err != nil {
		return err
	}

	dbConn, err := driver.Connect(cnx.DSN)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	result, err := driver.Query(dbConn, cmd.qry)
	if err != nil {
		return err
	}

	out := "table"
	if cmd.list {
		out = "list"
	}

	formatter.PrintResult(result, out, os.Stdout)

	return nil
}

func (cmd *RunCommand) Name() string {
	return "run"
}

func (cmd *RunCommand) IsUsed() bool {
	return cmd.Subcommand.Used
}

func fromStdin() string {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s", bytes)
}
