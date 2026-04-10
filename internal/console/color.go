package console

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

var reColor *regexp.Regexp = regexp.MustCompile(`(<[^>]+>)`)
var isColor = true

var colors = map[string]string{
	"black":     "\033[30m",
	"red":       "\033[31m",
	"green":     "\033[32m",
	"yellow":    "\033[33m",
	"blue":      "\033[34m",
	"purple":    "\033[35m",
	"cyan":      "\033[36m",
	"white":     "\033[37m",
	"hiblack":   "\033[90m",
	"hired":     "\033[91m",
	"higreen":   "\033[92m",
	"hiyellow":  "\033[93m",
	"hiblue":    "\033[94m",
	"hipurple":  "\033[95m",
	"hicyan":    "\033[96m",
	"hiwhite":   "\033[97m",
	"bgblack":   "\033[40m",
	"bgred":     "\033[41m",
	"bggreen":   "\033[42m",
	"bgyellow":  "\033[43m",
	"bgblue":    "\033[44m",
	"bgpurple":  "\033[45m",
	"bgcyan":    "\033[46m",
	"bgwhite":   "\033[47m",
	"hbblack":   "\033[100m",
	"hbred":     "\033[101m",
	"hbgreen":   "\033[102m",
	"hbyellow":  "\033[103m",
	"hbblue":    "\033[104m",
	"hbpurple":  "\033[105m",
	"hbcyan":    "\033[106m",
	"hbwhite":   "\033[107m",
	"bold":      "\033[1m",
	"italic":    "\033[3m",
	"underline": "\033[4m",
}

func init() {
	s, ok := os.LookupEnv("NO_COLOR")
	isColor = !(ok || s != "")
}

// Colorize devuelve una cadena con los colores definidos en el mapa colors
func Colorize(format string, a ...any) string {
	txt := format
	out := reColor.ReplaceAllStringFunc(
		txt,
		func(s string) string {
			if isColor {
				s := strings.ToLower(strings.Trim(s, " <>"))
				slog.Debug("s:", slog.String("s", s))
				if s == "/" || s == "reset" {
					return "\033[0m"
				} else {
					sb := strings.Builder{}
					for value := range strings.SplitSeq(s, "+") {
						slog.Debug("value:", slog.String("val", value))
						if color, ok := colors[value]; ok {
							sb.WriteString(color)
						}
					}
					return sb.String()
				}
			}

			return ""
		},
	)

	return fmt.Sprintf(out, a...)
}
