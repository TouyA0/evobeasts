package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
	"os"
	"os/exec"
	"runtime"
)

const (
	creatureArtWidth  = 55
)

func clearScreen() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
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

func truncateVisible(s string, maxWidth int) string {
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	if maxWidth <= 1 {
		return "…"
	}
	return string(runes[:maxWidth-1]) + "…"
}

func centerText(s string, width int) string {
	w := visibleWidth(s)
	if w >= width {
		return s
	}
	left := (width - w) / 2
	right := width - w - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func dividedColumns(left, right string, leftWidth int) string {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")
 
	height := len(leftLines)
	if len(rightLines) > height {
		height = len(rightLines)
	}
 
	var out strings.Builder
	for i := 0; i < height; i++ {
		l, r := "", ""
		if i < len(leftLines) {
			l = leftLines[i]
		}
		if i < len(rightLines) {
			r = rightLines[i]
		}
		out.WriteString(padRight(l, leftWidth))
		out.WriteString(" │ ")
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

func statusLine(shield int, effects []StatusEffect, invisible bool, invisibleBonus int, intangible bool, skipNextTurn bool) string {
	var parts []string

	if shield > 0 {
		parts = append(parts, fmt.Sprintf("Bouclier %d", shield))
	}
	if invisible {
		if invisibleBonus > 0 {
			parts = append(parts, fmt.Sprintf("Invisible (+%d%%)", invisibleBonus))
		} else {
			parts = append(parts, "Invisible")
		}
	}
	if intangible {
		parts = append(parts, "Intangible")
	}
	if skipNextTurn {
		parts = append(parts, "Étourdi")
	}
	for _, e := range effects {
		switch {
		case e.DamagePerTurn > 0:
			parts = append(parts, fmt.Sprintf("%s (%s)", e.Name, effectDuration(e.TurnsLeft)))
		case e.HealPerTurn > 0:
			parts = append(parts, fmt.Sprintf("%s (%s)", e.Name, effectDuration(e.TurnsLeft)))
		case e.AttackBuff > 0:
			parts = append(parts, fmt.Sprintf("ATK +%d (%s)", e.AttackBuff, effectDuration(e.TurnsLeft)))
		case e.AttackBuff < 0:
			parts = append(parts, fmt.Sprintf("ATK %d (%s)", e.AttackBuff, effectDuration(e.TurnsLeft)))
		case e.DamageReduction > 0:
			parts = append(parts, fmt.Sprintf("DEF +%d%% (%s)", e.DamageReduction, effectDuration(e.TurnsLeft)))
		case e.InitiativeDebuff > 0:
			parts = append(parts, fmt.Sprintf("INIT -%d (%s)", e.InitiativeDebuff, effectDuration(e.TurnsLeft)))
		}
	}

	if len(parts) == 0 {
		return ""
	}
	return "\x1b[38;2;150;150;170m" + strings.Join(parts, "  ·  ") + "\x1b[0m"
}

func effectDuration(turnsLeft int) string {
	if turnsLeft >= PermanentEffectDuration {
		return "permanent"
	}
	return fmt.Sprintf("%dt", turnsLeft)
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

func vulnerabilityHint(res map[Element]float64) string {
	var weak []string
	for _, e := range []Element{ElementFeu, ElementEau, ElementSpectre} {
		if v, ok := res[e]; ok && v > 1.0 {
			weak = append(weak, string(e))
		}
	}
	if len(weak) == 0 {
		return ""
	}
	return "Vulnérable: " + strings.Join(weak, "/")
}

func enemySkills(m *Monster) []string {
	if len(m.Phases) > 0 {
		return m.Phases[m.CurrentPhase].Attacks
	}
	return m.Skills
}

func nextPhaseHint(m *Monster) string {
	if len(m.Phases) == 0 || m.CurrentPhase >= len(m.Phases)-1 {
		return ""
	}
	nextThreshold := m.Phases[m.CurrentPhase+1].HPThreshold
	pvAtNext := int(nextThreshold * float64(m.MaxHP))
	remaining := m.CurrentHP - pvAtNext
	if remaining < 0 {
		remaining = 0
	}
	return fmt.Sprintf("Phase suivante dans %d PV", remaining)
}

func hasCastableUltimate(p *Character) bool {
	for _, s := range p.Spells {
		if s.IsUltimate && s.Name != SortRenaissance {
			return true
		}
	}
	return false
}

func combatScreen(p *Character, m *Monster, turn int, playerFirst bool) {
	clearScreen()

	var lines []string
	lines = append(lines, fmt.Sprintf(" Tour %d", turn))
	if playerFirst {
		lines = append(lines, "\x1b[38;2;100;220;120m▶ Vous jouez en premier\x1b[0m")
	} else {
		lines = append(lines, "\x1b[38;2;220;80;80m▶ L'ennemi joue en premier\x1b[0m")
	}
	lines = append(lines, "")

	enemyHeaderText := m.Name
	if len(m.Phases) > 0 {
		enemyHeaderText = fmt.Sprintf("%s (%s)", m.Name, m.Phases[m.CurrentPhase].Name)
	}
	if m.Element != ElementNeutre {
		enemyHeaderText += " · " + string(m.Element)
	}
	enemyHeaderText = truncateVisible(enemyHeaderText, creatureArtWidth)
	enemyHeader := "\x1b[38;2;220;80;80m" + centerText(enemyHeaderText, creatureArtWidth) + "\x1b[0m"
 
	playerHeaderText := truncateVisible(
		fmt.Sprintf("%s — %s (%s)", p.Name, p.Species.StageNames[p.Stage-1], p.Species.Element),
		creatureArtWidth)
	playerHeader := "\x1b[38;2;100;220;120m" + centerText(playerHeaderText, creatureArtWidth) + "\x1b[0m"

	enemyArt := ""
	if imgPath := monsterImagePath(m.Name); imgPath != "" {
		enemyArt, _ = renderImageString(imgPath, creatureArtWidth)
	}

	playerArt := ""
	if imgPath := speciesImagePath(p.Species.Name, p.Stage); imgPath != "" {
		playerArt, _ = renderImageString(imgPath, creatureArtWidth)
	}
 
		// --- Bloc d'infos ennemi (gauche) : PV, ATK, vulnérabilité, phase, compétences ---
	enemyInfoLines := []string{
		"PV   " + healthBar(m.CurrentHP, m.MaxHP, 20),
		fmt.Sprintf("ATK %d", m.Attack),
	}
	if status := statusLine(m.Shield, m.Effects, false, 0, m.Intangible, m.SkipNextTurn); status != "" {
		enemyInfoLines = append(enemyInfoLines, status)
	}
	if hint := vulnerabilityHint(m.Resistances); hint != "" {
		enemyInfoLines = append(enemyInfoLines,
			"\x1b[38;2;140;170;220m"+truncateVisible(hint, creatureArtWidth)+"\x1b[0m")
	}
	if hint := nextPhaseHint(m); hint != "" {
		enemyInfoLines = append(enemyInfoLines,
			"\x1b[38;2;180;140;220m"+truncateVisible(hint, creatureArtWidth)+"\x1b[0m")
	}
	if skills := enemySkills(m); len(skills) > 0 {
		skillLine := "Compétences: " + strings.Join(skills, ", ")
		enemyInfoLines = append(enemyInfoLines,
			"\x1b[38;2;150;150;170m"+truncateVisible(skillLine, creatureArtWidth)+"\x1b[0m")
	}

	// --- Bloc d'infos joueur (droite) ---
	playerInfoLines := []string{
		"PV   " + healthBar(p.CurrentHP, p.MaxHP, 20),
		"Mana " + healthBar(p.CurrentMana, p.MaxMana, 20),
		fmt.Sprintf("ATK %d", p.Attack),
	}
	if p.DodgeRate > 0 {
		playerInfoLines = append(playerInfoLines, fmt.Sprintf("Esquive %d%%", p.DodgeRate))
	}
	if status := statusLine(p.Shield, p.Effects, p.Invisible, p.InvisibleBonus, p.Intangible, p.SkipNextTurn); status != "" {
		playerInfoLines = append(playerInfoLines, status)
	}
	if hasCastableUltimate(p) {
		if p.UltimateUsed {
			playerInfoLines = append(playerInfoLines, "\x1b[38;2;150;150;170mUltime : utilisé\x1b[0m")
		} else {
			playerInfoLines = append(playerInfoLines, "\x1b[38;2;220;190;90mUltime : disponible\x1b[0m")
		}
	}
	playerInfoLines = append(playerInfoLines, "",
		"Tête  : "+truncateVisible(slotName(p.Equipment.Head), creatureArtWidth-8),
		"Torse : "+truncateVisible(slotName(p.Equipment.Chest), creatureArtWidth-8),
		"Pieds : "+truncateVisible(slotName(p.Equipment.Feet), creatureArtWidth-8),
	)
 
	leftColumn := strings.Join(append([]string{enemyHeader, ""},
		append(strings.Split(enemyArt, "\n"), append([]string{""}, enemyInfoLines...)...)...), "\n")
 
	rightColumn := strings.Join(append([]string{playerHeader, ""},
		append(strings.Split(playerArt, "\n"), append([]string{""}, playerInfoLines...)...)...), "\n")
 
	merged := dividedColumns(leftColumn, rightColumn, creatureArtWidth)
	lines = append(lines, strings.Split(strings.TrimRight(merged, "\n"), "\n")...)
 
	lines = append(lines, "", strings.Repeat("─", creatureArtWidth*2+3))
 
	fmt.Print(boxify(lines))
}

func slotName(item *Item) string {
	if item == nil {
		return "—"
	}
	return item.Name
}
