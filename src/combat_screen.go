package main

import (
	"fmt"
	"regexp"
	"strings"
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
func runeWidth(r rune) int {
	switch {
	case r >= 0x1F300 && r <= 0x1FAFF: // émojis pictographiques (💥🌀💀👻🌫️🔥💚)
		return 2
	case r >= 0x2600 && r <= 0x27BF: // symboles/dingbats (⚡✨☠️)
		return 2
	default:
		return 1
	}
}

func visibleWidth(s string) int {
	stripped := ansiEscape.ReplaceAllString(s, "")
	w := 0
	for _, r := range stripped {
		w += runeWidth(r)
	}
	return w
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

func healthBar(current, max, width int, thresholds []float64) string {
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

	markerCol := make(map[int]bool)
	for _, t := range thresholds {
		col := int(t * float64(width))
		if col > 0 && col < width {
			markerCol[col] = true
		}
	}
	const markerColor = "\x1b[38;2;190;150;230m" // violet, seuils de phase

	var b strings.Builder
	b.WriteString(color + "[")
	for i := 0; i < width; i++ {
		switch {
		case markerCol[i]:
			b.WriteString(markerColor + "┆" + color)
		case i < filled:
			b.WriteString("█")
		default:
			b.WriteString("░")
		}
	}
	b.WriteString("]\x1b[0m")
	b.WriteString(fmt.Sprintf(" %d/%d", current, max))
	return b.String()
}

func phaseThresholds(m *Monster) []float64 {
	if len(m.Phases) <= 1 {
		return nil
	}
	var t []float64
	for i := 1; i < len(m.Phases); i++ {
		t = append(t, m.Phases[i].HPThreshold)
	}
	return t
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

func resistanceLines(res map[Element]float64) (weakLine, strongLine string) {
	var weak, strong []string
	for _, e := range []Element{ElementFeu, ElementEau, ElementSpectre} {
		v, ok := res[e]
		if !ok || v == 1.0 {
			continue
		}
		tag := fmt.Sprintf("%s x%.1f", e, v)
		if v > 1.0 {
			weak = append(weak, tag)
		} else {
			strong = append(strong, tag)
		}
	}
	if len(weak) > 0 {
		weakLine = "Faible: " + strings.Join(weak, ", ")
	}
	if len(strong) > 0 {
		strongLine = "Résiste: " + strings.Join(strong, ", ")
	}
	return
}

func initiativeLine(base int, effects []StatusEffect) string {
	eff := effectiveInitiative(base, effects)
	if eff != base {
		return fmt.Sprintf("INIT %d (base %d)", eff, base)
	}
	return fmt.Sprintf("INIT %d", eff)
}

func enemySkills(m *Monster) []MonsterSkill {
	if len(m.Phases) > 0 {
		return m.Phases[m.CurrentPhase].Attacks
	}
	return m.Skills
}

func predictMonsterSkillName(m *Monster, turn int) string {
	switch m.Name {
	case MonstreVer:
		return "Morsure"

	case MonstreSanglier:
		if turn%3 == 0 {
			return "Charge"
		}
		return "Attaque"

	case MonstreLoupGarou:
		if !m.HowlUsed && turn == 3 {
			return "Hurlement"
		}
		if turn%2 == 0 {
			return "Griffure enragée"
		}
		return "Attaque"

	case MonstreDragonAncien:
		switch m.CurrentPhase {
		case 0:
			if turn%2 == 0 {
				return "Souffle ardent"
			}
			return "Griffe"
		case 1:
			if turn%3 == 0 {
				return "Étreinte liquide"
			}
			return "Raz-de-marée"
		case 2:
			if turn%4 == 0 {
				return "Dissolution"
			}
			return "Frappe hors du plan"
		case 3:
			element := ElementFeu
			switch turn % 3 {
			case 1:
				element = ElementEau
			case 2:
				element = ElementSpectre
			}
			return "Convergence élémentaire (" + string(element) + ")"
		}
	}
	return "Attaque"
}

func findSkill(skills []MonsterSkill, name string) MonsterSkill {
	for _, s := range skills {
		if s.Name == name || strings.HasPrefix(name, s.Name+" (") {
			return s
		}
	}
	return MonsterSkill{
		Name: "Attaque", Kind: SkillAttack, DamagePct: 100,
		Description: "Attaque simple (100% ATK).",
	}
}

func enemyIntentText(m *Monster, turn int) string {
	if m.SkipNextTurn {
		return "Intention : Étourdi — ne peut pas agir"
	}
	name := predictMonsterSkillName(m, turn)
	skill := findSkill(enemySkills(m), name)
	return fmt.Sprintf("Intention : %s %s", skillIcon(skill.Kind), name)
}

func skillIcon(k SkillKind) string {
	switch k {
	case SkillAttack:
		return "💥"
	case SkillBuff:
		return "⚡"
	case SkillDebuff:
		return "🌀"
	case SkillHeal:
		return "💚"
	default:
		return "✨"
	}
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

func tierLabel(tier int) string {
	switch tier {
	case 1:
		return "Novice"
	case 2:
		return "Intermédiaire"
	case 3:
		return "Avancé"
	case 4:
		return "Boss"
	default:
		return ""
	}
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
	enemyHeaderText += fmt.Sprintf(" · Tier %d (%s)", m.Tier, tierLabel(m.Tier))
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
		"PV   " + healthBar(m.CurrentHP, m.MaxHP, 20, phaseThresholds(m)),
		fmt.Sprintf("ATK %d", m.Attack),
		initiativeLine(m.Initiative, m.Effects),
	}
		if status := statusLine(m.Shield, m.Effects, false, 0, m.Intangible, m.SkipNextTurn); status != "" {
		enemyInfoLines = append(enemyInfoLines, status)
	}
	if intent := enemyIntentText(m, turn); intent != "" {
		enemyInfoLines = append(enemyInfoLines,
			"\x1b[38;2;220;170;60m"+truncateVisible(intent, creatureArtWidth)+"\x1b[0m")
	}
	if weakLine, strongLine := resistanceLines(m.Resistances); weakLine != "" || strongLine != "" {
		if weakLine != "" {
			enemyInfoLines = append(enemyInfoLines,
				"\x1b[38;2;220;90;90m"+truncateVisible(weakLine, creatureArtWidth)+"\x1b[0m")
		}
		if strongLine != "" {
			enemyInfoLines = append(enemyInfoLines,
				"\x1b[38;2;140;170;220m"+truncateVisible(strongLine, creatureArtWidth)+"\x1b[0m")
		}
	}
	if hint := nextPhaseHint(m); hint != "" {
		enemyInfoLines = append(enemyInfoLines,
			"\x1b[38;2;180;140;220m"+truncateVisible(hint, creatureArtWidth)+"\x1b[0m")
	}
	if skills := enemySkills(m); len(skills) > 0 {
		enemyInfoLines = append(enemyInfoLines, "Compétences :")
		for _, s := range skills {
			line := fmt.Sprintf(" %s %s — %s", skillIcon(s.Kind), s.Name, s.Description)
			enemyInfoLines = append(enemyInfoLines,
				"\x1b[38;2;150;150;170m"+truncateVisible(line, creatureArtWidth)+"\x1b[0m")
		}
	}

	// --- Bloc d'infos joueur (droite) ---
	playerInfoLines := []string{
		"PV   " + healthBar(p.CurrentHP, p.MaxHP, 20, nil),
		"Mana " + healthBar(p.CurrentMana, p.MaxMana, 20, nil),
		fmt.Sprintf("ATK %d", p.Attack),
		initiativeLine(p.Initiative, p.Effects),
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
