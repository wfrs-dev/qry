package formatter

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"gitlab.com/wfrsgo/qry/internal/db"
)

//  ──╢ FILA 0000 ╟──────────

// PrintResult formatea el resultado de la consulta según el formato solicitado ("table" o "list").
func PrintResult(res *db.QueryResult, format string, out io.Writer) {
	if len(res.Rows) == 0 {
		fmt.Fprintln(out, res.Summary)
		return
	}

	if format == "list" {
		printList(res, out)
	} else {
		// por defecto se muestra como tabla
		printTable(res, out)
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, res.Summary)
}

func printTable(res *db.QueryResult, out io.Writer) {
	// Inicializar tabwriter con un padding razonable
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	// Imprimir encabezados
	fmt.Fprintln(w, strings.Join(res.Columns, "\t"))

	// Imprimir separador
	sep := make([]string, len(res.Columns))
	for i, col := range res.Columns {
		sepLen := len(col)
		if sepLen < 3 {
			sepLen = 3
		}
		sep[i] = strings.Repeat("-", sepLen)
	}
	fmt.Fprintln(w, strings.Join(sep, "\t"))

	// Imprimir filas
	for _, row := range res.Rows {
		strRow := make([]string, len(row))
		for i, v := range row {
			strRow[i] = formatValue(v)
		}
		fmt.Fprintln(w, strings.Join(strRow, "\t"))
	}
	w.Flush()
}

func printList(res *db.QueryResult, out io.Writer) {
	for i, row := range res.Rows {
		if i > 0 {
			fmt.Fprintf(out, "\n\n──╢ FILA %04d ╟──────────\n", i+1)
		}
		for j, col := range res.Columns {
			fmt.Fprintf(out, "%15s : %s\n", col, formatValue(row[j]))
		}
	}
}

func formatValue(v interface{}) string {
	if v == nil {
		return "<nil>"
	}
	switch val := v.(type) {
	case []byte:
		// Algunos drivers devuelven datos de texto como []byte
		return string(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
