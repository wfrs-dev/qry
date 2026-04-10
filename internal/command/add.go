package command

import (
	"fmt"
	"sort"
	"strings"

	"github.com/integrii/flaggy"
	"gitlab.com/wfrsgo/qry/internal/config"
	"gitlab.com/wfrsgo/qry/internal/console"
	"gitlab.com/wfrsgo/qry/internal/db"
)

type AddCommand struct {
	*flaggy.Subcommand
	drivers []string
	name    string
}

// 󱘳
func (cmd *AddCommand) Setup() *flaggy.Subcommand {
	cmd.Subcommand = flaggy.NewSubcommand(cmd.Name())
	cmd.Subcommand.Description = "Agrega una nueva conexión"
	cmd.Subcommand.AddPositionalValue(&cmd.name, "nombre", 1, true, "Nombre de la conexión a agregar")

	return cmd.Subcommand
}

func (cmd *AddCommand) Execute(ctx *console.Context) error {
	drivers := make([]string, 0)
	instances := db.GetDriverInstances()
	for name := range instances {
		drivers = append(drivers, name)
	}
	sort.Strings(drivers)

	tstr := fmt.Sprintf("NUEVA CONEXIÓN `%s`:", cmd.name)

	fmt.Println(console.Colorize("<blue+bold>󱘳 %s\n%s</>", tstr, strings.Repeat("=", len([]rune(tstr))+1)))

	driver := console.Select("󰮆 Driver de base de datos", drivers)

	sdriver := drivers[driver]

	fmt.Println(console.Colorize("<green>󰪩 Driver seleccionado: <bold>%s</>\n<hipurple+italic>󰡦 Ejemplo de DSN: %s</>", sdriver, db.DSNDrivers[sdriver]))
	dsn := console.Input("󰮆 DSN de la conexión, siguiendo el formato de ejemplo", "")
	fmt.Println(console.Colorize("<yellow>󱘖 Validando DSN...</>"))
	d, err := db.GetDriver(sdriver)
	if err != nil {
		return err
	}

	cnx, err := d.Connect(dsn)
	if err != nil {
		return err
	}
	defer cnx.Close()

	fmt.Println(console.Colorize("<green> DSN validado correctamente</>"))

	err = config.AddConnection(cmd.name, sdriver, dsn)
	if err != nil {
		return err
	}

	fmt.Println(console.Colorize("<green>󱍭 Conexión agregada correctamente</>"))

	return nil
}

func (cmd *AddCommand) IsUsed() bool {
	return cmd.Subcommand.Used
}

func (cmd *AddCommand) Name() string {
	return "add"
}
