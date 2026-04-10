package console

import (
	"bufio"
	"os"
	"strings"

	"golang.org/x/term"
)

func GetStdin() (string, bool) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		sb := strings.Builder{}
		for scanner.Scan() {
			sb.WriteString(scanner.Text() + "\n")
		}
		return sb.String(), true
	}

	return "", false
}

func TerminalSize() (int, int) {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))

	return w, h
}
