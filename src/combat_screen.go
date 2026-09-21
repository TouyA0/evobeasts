package main

import (
	"fmt"
	"regexp"
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

var ansiEscape = regexp.MustCompile("\x1b\\[[0-9;]*m") // ignore les codes couleurs ANSI
func visibleWidth(s string) int {
	stripped := ansiEscape.ReplaceAllString(s, "")
	return utf8.RuneCountInString(stripped)
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

	fraction := float64(current) / float64(max)
	var color string
	switch {
	case fraction > 0.5:
		color = "\x1b[38;2;90;200;90m" // vert
	case fraction > 0.25:
		color = "\x1b[38;2;230;170;50m" // orange
	default:
		color = "\x1b[38;2;220;70;70m" // rouge
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return color + "[" + bar + "]" + "\x1b[0m" + fmt.Sprintf(" %d/%d", current, max)
}

func statusLine(shield int, effects []StatusEffect, invisible, intangible bool) string {
	var parts []string

	if shield > 0 {
		parts = append(parts, fmt.Sprintf("Bouclier %d", shield))
	}
	if invisible {
		parts = append(parts, "Invisible")
	}
	if intangible {
		parts = append(parts, "Intangible")
	}
	for _, e := range effects {
		switch {
		case e.DamagePerTurn > 0:
			parts = append(parts, fmt.Sprintf("%s (%dt)", e.Name, e.TurnsLeft))
		case e.HealPerTurn > 0:
			parts = append(parts, fmt.Sprintf("%s (%dt)", e.Name, e.TurnsLeft))
		case e.AttackBuff > 0:
			parts = append(parts, fmt.Sprintf("ATK +%d (%dt)", e.AttackBuff, e.TurnsLeft))
		case e.AttackBuff < 0:
			parts = append(parts, fmt.Sprintf("ATK %d (%dt)", e.AttackBuff, e.TurnsLeft))
		case e.DamageReduction > 0:
			parts = append(parts, fmt.Sprintf("DEF +%d%% (%dt)", e.DamageReduction, e.TurnsLeft))
		case e.InitiativeDebuff > 0:
			parts = append(parts, fmt.Sprintf("INIT -%d (%dt)", e.InitiativeDebuff, e.TurnsLeft))
		}
	}

	if len(parts) == 0 {
		return ""
	}
	return "\x1b[38;2;150;150;170m" + strings.Join(parts, "  ·  ") + "\x1b[0m"
}

func separator() string {
	return strings.Repeat("─", ScreenWidth)
}

// fitArt force un dessin ASCI dans une largeur et une hauteur max
func fitArt(art string, maxWidth, maxHeight int) string {
	lines := strings.Split(art, "\n")
	if len(lines) > maxHeight {
		lines = lines[:maxHeight]
	}
	for i, line := range lines {
		runes := []rune(line)
		if len(runes) > maxWidth {
			lines[i] = string(runes[:maxWidth])
		}
	}
	return strings.Join(lines, "\n")
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

func boxify(lines []string) string {
	width := 0
	for _, l := range lines {
		if w := visibleWidth(l); w > width {
			width = w
		}
	}

	var b strings.Builder
	b.WriteString("┌" + strings.Repeat("─", width+2) + "┐\n")
	for _, l := range lines {
		b.WriteString("│ " + padRight(l, width) + " │\n")
	}
	b.WriteString("└" + strings.Repeat("─", width+2) + "┘\n")
	return b.String()
}

func combatScreen(p *Character, m *Monster, turn int) {
	clearScreen()

	var lines []string
	lines = append(lines, fmt.Sprintf(" Tour %d", turn), "")

	if ca, ok := monsterArts[m.Name]; ok {
		fitted := fitArt(ca.Art, CombatArtWidth, CombatArtHeight)
		for _, line := range strings.Split(artToString(fitted, ca.Colors, CombatArtWidth), "\n") {
			lines = append(lines, center(line, ScreenWidth))
		}
	}
	lines = append(lines, "")
	monsterLabel := m.Name
	if len(m.Phases) > 0 {
		monsterLabel = fmt.Sprintf("%s (%s)", m.Name, m.Phases[m.CurrentPhase].Name)
	}
	lines = append(lines, center(fmt.Sprintf("%s  %s", monsterLabel,
		healthBar(m.CurrentHP, m.MaxHP, 12)), ScreenWidth))
	if status := statusLine(m.Shield, m.Effects, false, m.Intangible); status != "" {
		lines = append(lines, center(status, ScreenWidth))
	}
	lines = append(lines, "", separator())

	playerArt := ""
	key := fmt.Sprintf("%s-%d", p.Species.Name, p.Stage)
	if ca, ok := speciesArts[key]; ok {
		fitted := fitArt(ca.Art, PlayerArtWidth, PlayerArtHeight)
		playerArt = artToString(fitted, ca.Colors, PlayerArtWidth)
	}

	infos := fmt.Sprintf("%s — %s\nPV   %s\nMana %s",
		p.Name,
		p.Species.StageNames[p.Stage-1],
		p.Species.Element,
		healthBar(p.CurrentHP, p.MaxHP, 10),
		healthBar(p.CurrentMana, p.MaxMana, 10),
	)
	if status := statusLine(p.Shield, p.Effects, p.Invisible, p.Intangible); status != "" {
		infos += "\n" + status
	}
	infos += fmt.Sprintf("\n\nTête  : %s\nTorse : %s\nPieds : %s",
		slotName(p.Equipment.Head),
		slotName(p.Equipment.Chest),
		slotName(p.Equipment.Feet),
	)

	for _, l := range strings.Split(sideBySide(playerArt, infos, PlayerArtWidth, 3), "\n") {
		lines = append(lines, l)
	}
	lines = append(lines, separator())

	fmt.Print(boxify(lines))
}

func slotName(item *Item) string {
	if item == nil {
		return "—"
	}
	return item.Name
}