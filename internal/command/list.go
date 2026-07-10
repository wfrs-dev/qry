package command

import (
	"errors"
	"fmt"

	"github.com/integrii/flaggy"
	"gitlab.com/wfrsgo/qry/internal/config"
	"gitlab.com/wfrsgo/qry/internal/console"
)

type ListCommand struct {
	*flaggy.Subcommand
	name bool
}

func (cmd *ListCommand) Setup() *flaggy.Subcommand {
	cmd.Subcommand = flaggy.NewSubcommand(cmd.Name())
	cmd.Subcommand.Description = "Lista las conexiones guardadas"
	cmd.Subcommand.Bool(&cmd.name, "n", "name", "Muestra solo los nombres de las conexiones")

	return cmd.Subcommand
}

func (cmd *ListCommand) Execute(ctx *console.Context) error {
	conns, err := config.ListConnections()
	if err != nil {
		return fmt.Errorf("no se pudo cargar las conexiones: %w", err)
	}

	if len(conns) == 0 {
		return errors.New("No hay conexiones guardadas")
	}

	if !cmd.name {
		fmt.Printf(" %-20s %-15s %s\n", "NOMBRE", "DRIVER", "DSN")
		fmt.Println("----------------------------------------------------------")
	}
	for name, conn := range conns {
		if cmd.name {
			fmt.Println(name)
			continue
		}

		fmt.Printf(" %-20s %-15s %s\n", name, conn.Driver, conn.MaskedDSN())
	}

	return nil
}

func (cmd *ListCommand) Name() string {
	return "list"
}

func (cmd *ListCommand) IsUsed() bool {
	return cmd.Subcommand.Used
}
