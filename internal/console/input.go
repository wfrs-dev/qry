package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

type ValidatorFunc func(string) error

var ValidNullable ValidatorFunc = func(s string) error { return nil }

func Text(prompt string) string {
	sb := strings.Builder{}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(Colorize("<green+bold>%s:</>", prompt))
	fmt.Println(Colorize("<green+bold>ENTER + ; + ENTER</> <green>para finalizar</>"))
	for {
		fmt.Print(Colorize("<blue>┃</> "))
		scanner.Scan()
		line := scanner.Text()
		if line == ";" {
			break
		}
		sb.WriteString(line + "\n")
	}

	return sb.String()
}

func Confirm(prompt string, defTrue bool) bool {
	res := "n"
	if defTrue {
		res = "s"
	}
	yn := Input(prompt, res, func(s string) error {
		lower := strings.TrimSpace(strings.ToLower(s))
		if lower == "y" || lower == "s" || lower == "si" || lower == "yes" || lower == "n" || lower == "no" {
			return nil
		}
		return fmt.Errorf("opción no válida (sólo 'y' o 'n')")
	})
	yn = strings.TrimSpace(strings.ToLower(yn))

	return yn == "y" || yn == "yes" || yn == "si" || yn == "s"
}

func Input(prompt, defVal string, validator ...ValidatorFunc) string {
	var line string
	var fn ValidatorFunc = ValidNullable
	if len(validator) > 0 {
		fn = validator[0]
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(Colorize("<green+bold>%s</>", prompt))
		if defVal != "" {
			fmt.Print(Colorize(" <yellow+italic>[%s]</>", defVal))
		}
		fmt.Print(Colorize("<green+bold>:</> "))
		scanner.Scan()
		line = scanner.Text()
		if line == "" {
			line = defVal
		}
		err := fn(line)
		if err == nil {
			break
		}
		fmt.Println(Colorize("<red+italic>#! %s</>", err))
	}

	return line
}

func Select(prompt string, options []string, defpos ...int) int {
	fmt.Println(Colorize("<green+bold>%s:</>", prompt))
	sopt := make([]string, len(options))
	smax := 0
	for i, opt := range options {
		sopt[i] = fmt.Sprintf("(%d) %s", i+1, opt)
		l := utf8.RuneCountInString(sopt[i])
		if l > smax {
			smax = l
		}
	}

	w, _ := TerminalSize()
	printOptionsGrid(sopt, smax, w)
	dp := 0
	isdp := len(defpos) > 0
	if isdp {
		if defpos[0] < 0 || defpos[0] > len(options) {
			dp = defpos[0] - 1
		}
	}

	var sel int
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(Colorize("<green+bold>Opción</>"))
		if isdp {
			fmt.Print(Colorize(" <yellow+italic>[%d]</>", dp+1))
		}
		fmt.Print(Colorize("<green+bold>:</> "))
		in, err := reader.ReadString('\n')
		in = strings.TrimSpace(in)

		n, err := strconv.ParseInt(in, 10, 32)
		if err != nil {
			fmt.Println(Colorize("<red+italic>#! Opción `%s` no válida </>", in))
			continue
		}

		sel = int(n)
		if sel < 1 || sel > len(options) {
			fmt.Println(Colorize("<red+italic>#! Opción (%d) no válida </>", sel))
			continue
		}

		sel = sel - 1
		break
	}

	return int(sel)
}

func Password(prompt string, validator ValidatorFunc) (string, error) {
	var pwd string
	if validator == nil {
		return "", fmt.Errorf("la función de validación no puede ser nula")
	}
	for {
		fmt.Print(prompt + ": ")

		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", fmt.Errorf("error reading password: %v", err)
		}
		fmt.Println() // Nueva línea después de ingresar contraseña

		err = validator(string(passwordBytes))
		if err != nil {
			fmt.Println(Colorize("<red+italic>#! %s</>", err.Error()))
			continue
		}
		pwd = strings.TrimSpace(string(passwordBytes))
		break
	}

	return pwd, nil
}

func printOptionsGrid(options []string, max, termWidth int) {
	colWidth := max + 3 // padding
	cols := termWidth / colWidth
	if cols < 1 {
		cols = 1
	}

	for i, opt := range options {
		fmt.Printf("%-*s", colWidth, opt)
		if (i+1)%cols == 0 {
			fmt.Println()
		}
	}

	if len(options)%cols != 0 {
		fmt.Println()
	}
}
