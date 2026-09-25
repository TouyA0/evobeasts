package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type menuTheme struct {
	bg, fg string // codes RGB "r;g;b"
}

var (
	themeVert      = menuTheme{bg: "100;220;120", fg: "20;30;20"}
	themeBleu      = menuTheme{bg: "90;160;230", fg: "10;20;35"}
	themeOr        = menuTheme{bg: "230;190;80", fg: "40;30;10"}
	themeOrange    = menuTheme{bg: "220;120;50", fg: "30;15;5"}
	themeRouge     = menuTheme{bg: "200;60;60", fg: "255;240;240"}
	themeViolet    = menuTheme{bg: "160;110;220", fg: "20;10;30"}
	themeTurquoise = menuTheme{bg: "80;200;190", fg: "10;30;30"}
	themeGris      = menuTheme{bg: "120;120;120", fg: "245;245;245"}
)

type menuOption struct {
	label string
	theme menuTheme
}

const (
	menuHint     = "↑/↓ pour naviguer, Entrée pour valider"
	menuHintSide = "←/→ pour choisir, Entrée pour valider"
	menuHintBack = "↑/↓ pour naviguer, Entrée pour valider, Échap pour revenir"
)

func elementTheme(e Element) menuTheme {
	switch e {
	case ElementFeu:
		return themeOrange
	case ElementEau:
		return themeBleu
	case ElementSpectre:
		return themeViolet
	}
	return themeVert
}

var pendingKeys []byte

var arrowKeys = []struct{ seq, key string }{
	{"\x1b[A", "prev"}, {"\x1bOA", "prev"}, {"\x1b[D", "prev"}, {"\x1bOD", "prev"},
	{"\x1b[B", "next"}, {"\x1bOB", "next"}, {"\x1b[C", "next"}, {"\x1bOC", "next"},
}

func readKey() string {
	if len(pendingKeys) == 0 {
		buf := make([]byte, 64)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return "enter"
		}
		pendingKeys = buf[:n]
	}
	k := pendingKeys
	for _, a := range arrowKeys {
		if strings.HasPrefix(string(k), a.seq) {
			pendingKeys = k[len(a.seq):]
			return a.key
		}
	}
	if k[0] == 0x1b && len(k) > 1 && (k[1] == '[' || k[1] == 'O') {
		pendingKeys = nil
		return ""
	}
	pendingKeys = k[1:]
	switch k[0] {
	case '\r', '\n':
		return "enter"
	case 0x1b:
		return "esc"
	case 0x03:
		return "ctrlc"
	}
	return string(k[:1])
}

func isShortcut(key string, count int) bool {
	return len(key) == 1 && key[0] >= '1' && key[0] <= '9' && int(key[0]-'1') < count
}

func quitGame(fd int, oldState *term.State) {
	term.Restore(fd, oldState)
	fmt.Print("\x1b[?25h\n")
	os.Exit(0)
}

func selectMenu(count int, render func(selected int) string) int {
	return selectMenuHint(count, 0, menuHint, render)
}

func selectMenuFrom(count, start int, render func(selected int) string) int {
	return selectMenuHint(count, start, menuHint, render)
}

func selectMenuHint(count, start int, hint string, render func(selected int) string) int {
	draw := func(sel int) []string {
		screen := strings.TrimRight(render(sel), "\n") + "\n\n" + hint
		return fitScreen(strings.Split(screen, "\n"))
	}

	selected := min(max(start, 0), count-1)
	prev := draw(selected)

	clearScreen()
	fmt.Print("\x1b[?25l")
	fmt.Print(strings.Join(prev, "\n"))
	defer func() { fmt.Printf("\x1b[%d;1H\x1b[J\x1b[?25h", len(prev)) }()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return selected
	}
	defer term.Restore(fd, oldState)
	flushInput()

	for {
		switch key := readKey(); {
		case key == "prev":
			selected = (selected - 1 + count) % count
		case key == "next":
			selected = (selected + 1) % count
		case key == "enter":
			return selected
		case key == "ctrlc":
			quitGame(fd, oldState)
		case isShortcut(key, count):
			return int(key[0] - '1')
		default:
			continue
		}

		// Redessine uniquement les lignes qui ont changé.
		next := draw(selected)
		var b strings.Builder
		for i, line := range next {
			if i >= len(prev) || line != prev[i] {
				fmt.Fprintf(&b, "\x1b[%d;1H%s\x1b[K", i+1, line)
			}
		}
		if len(next) < len(prev) {
			fmt.Fprintf(&b, "\x1b[%d;1H\x1b[J", len(next)+1)
		}
		os.Stdout.WriteString(b.String())
		prev = next
	}
}

func inlineSelect(count, start int, hint string, render func(selected int) string) int {
	height := 0 // nombre de lignes du menu affiché, pour pouvoir remonter
	redraw := func(sel int) {
		var b strings.Builder
		if height > 0 {
			fmt.Fprintf(&b, "\r\x1b[%dA\x1b[J", height)
		}
		block := strings.Join(fitWidth(strings.Split(strings.TrimRight(render(sel), "\n")+"\n\n"+hint, "\n")), "\n")
		b.WriteString(block)
		height = strings.Count(block, "\n")
		fmt.Print(b.String())
	}
	erase := func() {
		fmt.Printf("\r\x1b[%dA\x1b[J", height)
	}

	selected := min(max(start, 0), count-1)
	fmt.Print("\x1b[?25l")
	redraw(selected)
	defer fmt.Print("\x1b[?25h")
	defer erase()

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return selected
	}
	defer term.Restore(fd, oldState)
	flushInput()

	for {
		switch key := readKey(); {
		case key == "prev":
			selected = (selected - 1 + count) % count
		case key == "next":
			selected = (selected + 1) % count
		case key == "enter":
			return selected
		case key == "esc":
			return -1
		case key == "ctrlc":
			quitGame(fd, oldState)
		case isShortcut(key, count):
			return int(key[0] - '1')
		default:
			continue
		}
		redraw(selected)
	}
}

func menuLine(label string, width int, selected bool, t menuTheme) string {
	pad := strings.Repeat(" ", max(0, width-visibleWidth(label)))
	if selected {
		return "\x1b[48;2;" + t.bg + "m\x1b[38;2;" + t.fg + "m\x1b[1m ▶ " + label + pad + " \x1b[0m"
	}
	return "   " + label + pad + " "
}

func maxWidth(labels ...string) int {
	width := 0
	for _, l := range labels {
		width = max(width, visibleWidth(l))
	}
	return width
}

func menuLines(options []menuOption, selected int) []string {
	labels := make([]string, len(options))
	for i, o := range options {
		labels[i] = o.label
	}
	width := maxWidth(labels...)

	lines := make([]string, len(options))
	for i, o := range options {
		lines[i] = menuLine(o.label, width, i == selected, o.theme)
	}
	return lines
}

// confirm affiche une question Oui / Non (Non sélectionné par défaut).
func confirm(lines ...string) bool {
	options := []menuOption{{"Oui", themeRouge}, {"Non", themeGris}}
	choice := selectMenuFrom(len(options), 1, func(sel int) string {
		return strings.Join(lines, "\n") + "\n\n" + strings.Join(menuLines(options, sel), "\n")
	})
	return choice == 0
}

// button affiche un bouton déjà sélectionné et attend Entrée.
func button(label string, t menuTheme) {
	fmt.Print("\x1b[?25l\n" + menuLine(label, 0, true, t))
	defer fmt.Print("\x1b[?25h\n")

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return
	}
	defer term.Restore(fd, oldState)
	flushInput()

	for {
		switch readKey() {
		case "enter":
			return
		case "ctrlc":
			quitGame(fd, oldState)
		}
	}
}

func backButton(t menuTheme) {
	button("Retour", t)
}

func fitWidth(lines []string) []string {
	cols, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return lines
	}
	for i, l := range lines {
		if visibleWidth(l) >= cols {
			lines[i] = truncateVisible(l, cols-1) + "\x1b[0m"
		}
	}
	return lines
}

func fitScreen(lines []string) []string {
	lines = fitWidth(lines)
	_, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || len(lines) < rows {
		return lines
	}
	body, hint := lines[:len(lines)-1], lines[len(lines)-1]
	h := rows - 2
	sel := 0
	for i, l := range body {
		if strings.Contains(l, " ▶ ") {
			sel = i
			break
		}
	}
	top := min(max(sel-h/2, 0), len(body)-h)
	return append(append([]string{}, body[top:top+h]...), hint)
}
