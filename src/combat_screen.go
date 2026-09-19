package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	CombatArtWidth   = 34
	CombatArtHeight  = 12
	PlayerArtWidth   = 28
	PlayerArtHeight  = 10
	ScreenWidth      = 60
)

func clearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

func visibleWidth(s string) int {
	return utf8.RuneCountInString(s)
}

func padRight(s string, width int) string {
	w := visibleWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

func center(s string, width int) string {
	w := visibleWidth(s)
	if w >= width {
		return s
	}
	left := (width - w) / 2
	return strings.Repeat(" ", left) + s
}

func sideBySide(left, right string, leftWidth, gap int) string {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")

	height := len(leftLines)
	if len(rightLines) > height {
		height = len(rightLines)
	}

	var out strings.Builder
	for i := 0; i < height; i++ {
		l := ""
		if i < len(leftLines) {
			l = leftLines[i]
		}
		r := ""
		if i < len(rightLines) {
			r = rightLines[i]
		}
		out.WriteString(padRight(l, leftWidth))
		out.WriteString(strings.Repeat(" ", gap))
		out.WriteString(r)
		out.WriteString("\n")
	}
	return out.String()
}

func healthBar(current, max, width int) string {
	if max <= 0 {
		max = 1
	}
	if current < 0 {
		current = 0
	}
	filled := current * width / max
	if filled > width {
		filled = width
	}
	return fmt.Sprintf("[%s%s] %d/%d",
		strings.Repeat("█", filled),
		strings.Repeat("░", width-filled),
		current, max)
}

func separator() string {
	return strings.Repeat("─", ScreenWidth)
}

func artToString(art string, colors [][3]int, width int) string {
	lines := strings.Split(art, "\n")
	var out strings.Builder
	for i, line := range lines {
		c := colors[i%len(colors)]
		padded := padRight(line, width)
		out.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m",
			c[0], c[1], c[2], padded))
		if i < len(lines)-1 {
			out.WriteString("\n")
		}
	}
	return out.String()
}

func combatScreen(p *Character, m *Monster, turn int) {
	clearScreen()

	fmt.Printf("  Tour %d\n\n", turn)

	if ca, ok := monsterArts[m.Name]; ok {
		for _, line := range strings.Split(artToString(ca.Art, ca.Colors, CombatArtWidth), "\n") {
			fmt.Println(center(line, ScreenWidth))
		}
	}
	fmt.Println()
	fmt.Println(center(fmt.Sprintf("%s  %s", m.Name,
		healthBar(m.CurrentHP, m.MaxHP, 12)), ScreenWidth))
	fmt.Println()
	fmt.Println(separator())

	playerArt := ""
	key := fmt.Sprintf("%s-%d", p.Species.Name, p.Stage)
	if ca, ok := speciesArts[key]; ok {
		playerArt = artToString(ca.Art, ca.Colors, PlayerArtWidth)
	}

	infos := fmt.Sprintf("%s — %s\nPV   %s\nMana %s\n\nTête  : %s\nTorse : %s\nPieds : %s",
		p.Name,
		p.Species.StageNames[p.Stage-1],
		healthBar(p.CurrentHP, p.MaxHP, 10),
		healthBar(p.CurrentMana, p.MaxMana, 10),
		slotName(p.Equipment.Head),
		slotName(p.Equipment.Chest),
		slotName(p.Equipment.Feet),
	)

	fmt.Print(sideBySide(playerArt, infos, PlayerArtWidth, 3))
	fmt.Println(separator())
}

func slotName(item *Item) string {
	if item == nil {
		return "—"
	}
	return item.Name
}
