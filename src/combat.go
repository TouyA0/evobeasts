package main
 
import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var combatMonsters = []menuOption{
	{MonstreVer, themeVert},
	{MonstreSanglier, themeOrange},
	{MonstreLoupGarou, themeViolet},
	{MonstreDragonAncien, themeRouge},
}

func combatMenu(p *Character) {
	banner, _ := renderImageString(bannerCombat, bannerWidth)
	selected := 0

	for {
		unlocked := min(p.MaxTierDefeated+1, len(combatMonsters))

		options := []menuOption{}
		for i := 0; i < unlocked; i++ {
			options = append(options, combatMonsters[i])
		}
		options = append(options, menuOption{"Retour", themeGris})

		choice := selectMenuFrom(len(options), selected, func(sel int) string {
			return banner + "\n" + boxify(menuLines(options, sel))
		})
		if choice == unlocked {
			return
		}
		selected = choice

		tier := choice + 1
		monster := monsterByTier(tier)
		storyPreFight(p, tier)
		trainingFight(p, monster)

		if isGameComplete(p) {
			storyFinale(p)
			return
		}
	}
}

func trainingFight(p *Character, m *Monster) {
	p.Shield = 0
	p.Effects = nil
	p.Invisible = false
	p.InvisibleBonus = 0
	p.Intangible = false
	p.UltimateUsed = false
	p.SkipNextTurn = false

	if m.Name == MonstreDragonAncien {
		playMusic(musicBoss)
	} else {
		playMusic(musicCombat)
	}
	defer playMusic(musicAmbiance)
 
	turn := 1
 
	for {
		playerFirst := effectiveInitiative(p.Initiative, p.Effects) >=
			effectiveInitiative(m.Initiative, m.Effects)

		combatScreen(p, m, turn, playerFirst)
 
		if playerFirst {
			runCharacterTurn(p, m)
			if m.CurrentHP <= 0 {
				endFight(p, m, true)
				pauseForInput()
				storyPostFight(p, m.Tier)
				return
			}
			runMonsterTurn(p, m, turn)
			if checkPlayerDefeat(p) {
				pauseForInput()
				return
			}
		} else {
			runMonsterTurn(p, m, turn)
			if checkPlayerDefeat(p) {
				pauseForInput()
				return
			}
			runCharacterTurn(p, m)
			if m.CurrentHP <= 0 {
				endFight(p, m, true)
				pauseForInput()
				storyPostFight(p, m.Tier)
				return
			}
		}
 
		resolveEffects(p, m)
		if m.CurrentHP <= 0 {
			endFight(p, m, true)
			pauseForInput()
			storyPostFight(p, m.Tier)
			return
		}
		if checkPlayerDefeat(p) {
			pauseForInput()
			return
		}
		autoAdvance(3)
		turn++
	}
}

func pauseForInput() {
	button("Continuer", themeVert)
}

func printSequential(lines ...string) {
	const charDelay = 25 * time.Millisecond

	for _, line := range lines {
		for _, r := range line {
			fmt.Print(string(r))
			time.Sleep(charDelay)
		}
		fmt.Println()
		time.Sleep(200 * time.Millisecond)
	}
}

func autoAdvance(seconds int) {
	for i := seconds; i > 0; i-- {
		fmt.Print("\r\033[K")

		if i == 1 {
			fmt.Printf("(Tour suivant dans %d seconde...)", i)
		} else {
			fmt.Printf("(Tour suivant dans %d secondes...)", i)
		}

		time.Sleep(1 * time.Second)
	}

	fmt.Print("\r\033[K")
}

func checkPlayerDefeat(p *Character) bool {
	if p.CurrentHP > 0 {
		return false
	}

	if hasSpell(p, SortRenaissance) && !p.UltimateUsed {
		p.UltimateUsed = true
		p.CurrentHP = p.MaxHP * RenaissanceHPPct / 100
		p.Effects = applyEffect(p.Effects, StatusEffect{
			Name: SortRenaissance, TurnsLeft: RenaissanceBuffTurns, AttackBuff: RenaissanceAtkBuff,
		})
		fmt.Println("\n🔥 Renaissance ! " + p.Name + " renaît de ses cendres avec " +
			strconv.Itoa(p.CurrentHP) + " PV et frappe plus fort pour un temps.")
		return false
	}

	fmt.Println("\n" + p.Name + " est vaincu...")
	isDead(p)
	return true
}

func runCharacterTurn(p *Character, m *Monster) {
	p.Intangible = false

	if p.SkipNextTurn {
		fmt.Println(p.Name + " ne peut pas agir ce tour-ci.")
		p.SkipNextTurn = false
		return
	}

	actions := []menuOption{
		{"Attaquer", themeRouge},
		{"Sorts", themeViolet},
		{"Inventaire", themeBleu},
	}
	message := ""
	selected := 0

	for {
		choice := inlineSelect(len(actions), selected, menuHintSide, func(sel int) string {
			bar := make([]string, len(actions))
			for i, a := range actions {
				bar[i] = menuLine(a.label, 0, i == sel, a.theme)
			}
			text := "\n" + strings.Join(bar, "  ")
			if message != "" {
				text += "\n\n" + message
			}
			return text
		})
		if choice == -1 {
			continue // Échap : il n'y a rien à annuler à ce niveau
		}
		selected = choice
		message = ""

		switch choice {
		case 0:
			basicAttack(p, m)
			return
		case 1:
			used, msg := spellMenu(p, m)
			if used {
				return
			}
			message = msg
		case 2:
			used, msg := combatInventoryMenu(p, m)
			if used {
				return
			}
			message = msg
		}
	}
}

func basicAttack(p *Character, m *Monster) {
	base := p.Attack + effectiveAttackBonus(p.Effects)
	if base < 0 {
		base = 0
	}
	if p.Invisible {
		base += base * p.InvisibleBonus / 100
		p.Invisible = false
		p.InvisibleBonus = 0
	}
	dmg := dealDamage(base, p.Species.Element, m, false)
	printSequential(
		fmt.Sprintf("%s inflige %d dégâts à %s", p.Name, dmg, m.Name),
		fmt.Sprintf("%s  PV : %d/%d", m.Name, m.CurrentHP, m.MaxHP),
	)
	announcePhaseChange(m)
}

func spellIcon(s Spell) string {
	switch {
	case s.IsUltimate:
		return "💀"
	case s.Heal > 0:
		return "💚"
	case s.ShieldValue > 0:
		return "🛡️"
	case s.Invisible:
		return "👻"
	case s.Intangible:
		return "🌫️"
	case s.DotDamage > 0:
		return "🔥"
	case s.AttackBuff > 0:
		return "⚡"
	case s.Damage > 0:
		return "💥"
	default:
		return "✨"
	}
}

func spellMenu(p *Character, m *Monster) (bool, string) {
	castable := make([]Spell, 0, len(p.Spells))
	for _, s := range p.Spells {
		if s.Name != SortRenaissance {
			castable = append(castable, s)
		}
	}

	if len(castable) == 0 {
		return false, "Aucun sort connu."
	}

	message := ""
	selected := 0

	for {
		choice := inlineSelect(len(castable)+1, selected, menuHintBack, func(sel int) string {
			return spellMenuText(p, castable, sel, message)
		})
		if choice == -1 || choice == len(castable) {
			return false, ""
		}
		selected = choice

		spell := castable[choice]
		if ok, reason := canCastSpell(p, spell); !ok {
			message = "Sort inutilisable : " + reason
			continue
		}

		result, ok := castSpell(p, m, spell)
		if !ok {
			message = result
			continue
		}
		if m != nil {
			printSequential(result, fmt.Sprintf("%s  PV : %d/%d", m.Name, m.CurrentHP, m.MaxHP))
		} else {
			fmt.Println(result)
		}
		announcePhaseChange(m)
		return true, ""
	}
}

func spellMenuText(p *Character, castable []Spell, sel int, message string) string {
	labels := make([]string, len(castable))
	for i, s := range castable {
		tag := ""
		if s.Element != p.Species.Element {
			tag += fmt.Sprintf(" [%s]", s.Element)
		}
		if s.IsUltimate {
			tag += " (ULTIME)"
		}
		labels[i] = fmt.Sprintf("%s %s%s (%d mana)", spellIcon(s), s.Name, tag, s.ManaCost)
	}
	width := maxWidth(append(labels, "Retour")...)

	lines := []string{"Sorts :"}
	for i, s := range castable {
		ok, reason := canCastSpell(p, s)
		theme := elementTheme(s.Element)
		if s.IsUltimate {
			theme = themeOr
		}
		if !ok {
			theme = themeGris
		}

		line := menuLine(labels[i], width, i == sel, theme)
		if !ok {
			if i != sel {
				line = "\x1b[38;2;220;70;70m" + line + "\x1b[0m"
			}
			line += "  \x1b[38;2;220;70;70m— " + reason + "\x1b[0m"
		}
		lines = append(lines, line, "      \x1b[90m"+s.Description+"\x1b[0m")
	}

	lines = append(lines, "", menuLine("Retour", width, sel == len(castable), themeGris))
	if message != "" {
		lines = append(lines, "", message)
	}
	return strings.Join(lines, "\n")
}

func combatInventoryMenu(p *Character, m *Monster) (bool, string) {
	var usable []inventoryEntry
	for _, e := range inventoryEntries(p.Inventory) {
		it := p.Inventory[e.index]
		if it.HealValue > 0 || it.ManaValue > 0 || it.PoisonDamage > 0 {
			usable = append(usable, e)
		}
	}
	if len(usable) == 0 {
		return false, "Aucun objet utilisable en combat."
	}

	choice := inlineSelect(len(usable)+1, 0, menuHintBack, func(sel int) string {
		labels := []string{"Retour"}
		for _, e := range usable {
			labels = append(labels, e.label)
		}
		width := maxWidth(labels...)

		lines := []string{"Inventaire :"}
		for i, e := range usable {
			line := menuLine(e.label, width, i == sel, e.theme)
			if e.desc != "" {
				line += "  \x1b[90m" + e.desc + "\x1b[0m"
			}
			lines = append(lines, line)
		}
		lines = append(lines, "", menuLine("Retour", width, sel == len(usable), themeGris))
		return strings.Join(lines, "\n")
	})
	if choice == -1 || choice == len(usable) {
		return false, ""
	}

	item := p.Inventory[usable[choice].index]
	fmt.Println("Vous utilisez " + item.Name)

	switch {
	case item.HealValue > 0:
		p.CurrentHP += item.HealValue
		if p.CurrentHP > p.MaxHP {
			p.CurrentHP = p.MaxHP
		}
		fmt.Printf("%s  PV : %d/%d\n", p.Name, p.CurrentHP, p.MaxHP)
	case item.ManaValue > 0:
		p.CurrentMana += item.ManaValue
		if p.CurrentMana > p.MaxMana {
			p.CurrentMana = p.MaxMana
		}
		fmt.Printf("%s  Mana : %d/%d\n", p.Name, p.CurrentMana, p.MaxMana)
	case item.PoisonDamage > 0:
		m.Effects = applyEffect(m.Effects, StatusEffect{
			Name: item.Name, TurnsLeft: item.PoisonTicks,
			DamagePerTurn: item.PoisonDamage,
		})
		fmt.Println(item.Name + " est jetée sur " + m.Name + " !")
	}

	removeInventory(p, item.Name)
	return true, ""
}

func announcePhaseChange(m *Monster) {
	if m == nil {
		return
	}
	if msg := updatePhase(m); msg != "" {
		fmt.Println(msg)
		if line, ok := storyPhaseLines[m.Phases[m.CurrentPhase].Name]; ok {
			fmt.Println("\x1b[38;2;200;180;255m" + line + "\x1b[0m")
		}
	}
}

func runMonsterTurn(p *Character, m *Monster, turn int) {
	if m.SkipNextTurn {
		fmt.Println(m.Name + " ne peut pas agir ce tour-ci.")
		m.SkipNextTurn = false
		m.StunImmune = true
		return
	}
	m.StunImmune = false
 
	switch m.Name {
	case MonstreVer:
		verPattern(p, m)
	case MonstreSanglier:
		sanglierPattern(p, m, turn)
	case MonstreLoupGarou:
		loupGarouPattern(p, m, turn)
	case MonstreDragonAncien:
		dragonPattern(p, m, turn)
	default:
		monsterBasicAttack(p, m, 100)
	}
}

func monsterBasicAttack(p *Character, m *Monster, percent int) bool {
	base := monsterDamage(m, percent)
	dmg, dodged := dealDamageToPlayer(base, m.Element, p)
	if dodged {
		fmt.Printf("%s attaque, %s esquive !\n", m.Name, p.Name)
		return false
	}
	printSequential(
		fmt.Sprintf("%s inflige %d dégâts à %s", m.Name, dmg, p.Name),
		fmt.Sprintf("%s  PV : %d/%d", p.Name, p.CurrentHP, p.MaxHP),
	)
	return true
}

func verPattern(p *Character, m *Monster) {
	monsterBasicAttack(p, m, 100)
}

func sanglierPattern(p *Character, m *Monster, turn int) {
	percent := 100
	charge := turn%3 == 0
	if charge {
		percent = sanglierChargePct
		fmt.Println(m.Name + " charge !")
	}

	fraction := float64(m.CurrentHP) / float64(m.MaxHP)
	if fraction <= 0.25 {
		percent = percent * sanglierEnrage25Pct / 100
	} else if fraction <= 0.50 {
		percent = percent * sanglierEnrage50Pct / 100
	}

	if monsterBasicAttack(p, m, percent) && charge && rand.Intn(100) < sanglierStunChance {
		p.SkipNextTurn = true
		fmt.Println(p.Name + " est étourdi pour le prochain tour !")
	}
}

func loupGarouPattern(p *Character, m *Monster, turn int) {
	if !m.HowlUsed && turn >= 3 {
		m.Attack += loupGarouHowlBonusAtk
		m.HowlUsed = true
		fmt.Println(m.Name + " hurle à la lune ! Son attaque augmente.")
	}
 
	// Régénération lunaire posée une seule fois, en continu.
	hasRegen := false
	for _, e := range m.Effects {
		if e.Name == "Régénération lunaire" {
			hasRegen = true
			break
		}
	}
	if !hasRegen {
		m.Effects = applyEffect(m.Effects, StatusEffect{
			Name: "Régénération lunaire", TurnsLeft: PermanentEffectDuration, HealPerTurn: loupGarouRegenHeal,
		})
	}
 
	if turn%2 == 0 {
		fmt.Println(m.Name + " lacère avec fureur !")
		monsterBasicAttack(p, m, loupGarouClawPct)
		if p.CurrentHP > 0 {
			monsterBasicAttack(p, m, loupGarouClawPct)
		}
	} else {
		monsterBasicAttack(p, m, 100)
	}
}

func dragonPattern(p *Character, m *Monster, turn int) {
	switch m.CurrentPhase {
	case 0: // Éveil Igné
		if turn%2 == 0 {
			fmt.Println(m.Name + " crache un souffle ardent !")
			if monsterBasicAttack(p, m, dragonSouffleArdentPct) {
				p.Effects = applyEffect(p.Effects, StatusEffect{
					Name: "Brûlure du Dragon", TurnsLeft: dragonBrulureDuree, DamagePerTurn: dragonBrulureDmg,
				})
			}
		} else {
			monsterBasicAttack(p, m, dragonGriffePct)
		}
 
	case 1: // Vague Ancestrale
		hasRegen := false
		for _, e := range m.Effects {
			if e.Name == "Régénération abyssale" {
				hasRegen = true
			}
		}
		if !hasRegen {
			m.Effects = applyEffect(m.Effects, StatusEffect{
				Name: "Régénération abyssale", TurnsLeft: PermanentEffectDuration, HealPerTurn: dragonRegenAbyssale,
			})
		}
		if turn%3 == 0 {
			fmt.Println(m.Name + " vous enveloppe d'une étreinte liquide !")
			if monsterBasicAttack(p, m, dragonEtreinteLiquidePct) {
				p.SkipNextTurn = true
			}
		} else {
			fmt.Println(m.Name + " déchaîne un raz-de-marée !")
			monsterBasicAttack(p, m, dragonRazDeMareePct)
		}
 
	case 2: // Voile Spectral
		if turn%4 == 0 {
			fmt.Println(m.Name + " se dissout hors du plan et esquive tout ce tour.")
			m.Intangible = true
		} else {
			fmt.Println(m.Name + " frappe hors du plan, ignorant tout bouclier !")
			base := monsterDamage(m, dragonFrappeHorsPlanPct)
			shield := p.Shield
			p.Shield = 0 // la frappe ignore le bouclier
			dmg, dodged := dealDamageToPlayer(base, m.Element, p)
			if dodged {
				p.Shield = shield // esquivé : le bouclier est intact
				fmt.Println(p.Name + " esquive !")
			} else {
				printSequential(
					fmt.Sprintf("%s inflige %d dégâts à %s", m.Name, dmg, p.Name),
					fmt.Sprintf("%s  PV : %d/%d", p.Name, p.CurrentHP, p.MaxHP),
				)
			}
		}
 
	case 3: // Convergence
		fmt.Println(m.Name + " déchaîne la puissance des trois éléments !")
		pattern := turn % 3
		var element Element
		switch pattern {
		case 0:
			element = ElementFeu
		case 1:
			element = ElementEau
		default:
			element = ElementSpectre
		}
		base := monsterDamage(m, dragonConvergencePct)
		dmg, dodged := dealDamageToPlayerElement(base, element, p)
		if dodged {
			fmt.Println(p.Name + " esquive !")
		} else {
			printSequential(
				fmt.Sprintf("%s inflige %d dégâts à %s", m.Name, dmg, p.Name),
				fmt.Sprintf("%s  PV : %d/%d", p.Name, p.CurrentHP, p.MaxHP),
			)
		}
	}
}

func monsterResistance(m *Monster, e Element) float64 {
	if mult, ok := m.Resistances[e]; ok {
		return mult
	}
	return 1.0
}

func effectiveAttackBonus(effects []StatusEffect) int {
	total := 0
	for _, e := range effects {
		total += e.AttackBuff
	}
	return total
}

func monsterDamage(m *Monster, percent int) int {
	atk := m.Attack + effectiveAttackBonus(m.Effects)
	if atk < 0 {
		atk = 0
	}
	base := atk * percent / 100
	if len(m.Phases) > 0 {
		base += base * m.Phases[m.CurrentPhase].DamageBonus / 100
	}
	return base
}

func applyEffect(effects []StatusEffect, e StatusEffect) []StatusEffect {
	if e.DamagePerTurn == 0 {
		for i := range effects {
			if effects[i].Name == e.Name {
				effects[i] = e
				return effects
			}
		}
	}
	return append(effects, e)
}

func dealDamage(base int, element Element, target *Monster, ignoreResist bool) int {
	// La Dissolution du dragon lui fait esquiver la prochaine attaque.
	if target.Intangible {
		target.Intangible = false
		fmt.Println(target.Name + " est intangible : l'attaque le traverse sans effet.")
		return 0
	}

	mult := monsterResistance(target, element)
	if ignoreResist && mult < 1.0 {
		mult = 1.0
	}
	amount := int(float64(base) * mult)
	if amount < 0 {
		amount = 0
	}

	if target.Shield > 0 {
		if target.Shield >= amount {
			target.Shield -= amount
			amount = 0
		} else {
			amount -= target.Shield
			target.Shield = 0
		}
	}

	target.CurrentHP -= amount
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}
	return amount
}

func dealDamageToPlayer(base int, element Element, p *Character) (int, bool) {
	return dealDamageToPlayerElement(base, element, p)
}

func dealDamageToPlayerElement(base int, element Element, p *Character) (int, bool) {
	dodge := p.DodgeRate
	if p.Invisible {
		dodge += InvisibleDodgeBonus
	}
	if rand.Intn(100) < dodge {
		return 0, true
	}
 
	mult := getResistance(p.Species, element)
	amount := int(float64(base) * mult)
 
	for _, e := range p.Effects {
		if e.DamageReduction > 0 {
			amount -= amount * e.DamageReduction / 100
		}
	}
	if amount < 0 {
		amount = 0
	}
 
	if p.Intangible {
		p.Intangible = false
		return 0, true
	}
 
	if p.Shield > 0 {
		if p.Shield >= amount {
			p.Shield -= amount
			amount = 0
		} else {
			amount -= p.Shield
			p.Shield = 0
		}
	}
 
	p.CurrentHP -= amount
	if p.CurrentHP < 0 {
		p.CurrentHP = 0
	}
	return amount, false
}

func effectiveInitiative(base int, effects []StatusEffect) int {
	total := base
	for _, e := range effects {
		total -= e.InitiativeDebuff
	}
	return total
}

func resolveEffects(p *Character, m *Monster) {
	p.Effects = tickEffects(p.Effects, func(dmg int) {
		p.CurrentHP -= dmg
		if p.CurrentHP < 0 {
			p.CurrentHP = 0
		}
	}, func(heal int) {
		p.CurrentHP += heal
		if p.CurrentHP > p.MaxHP {
			p.CurrentHP = p.MaxHP
		}
	})
 
	if m != nil && m.CurrentHP > 0 {
		m.Effects = tickEffects(m.Effects, func(dmg int) {
			m.CurrentHP -= dmg
			if m.CurrentHP < 0 {
				m.CurrentHP = 0
			}
			announcePhaseChange(m)
		}, func(heal int) {
			if m.CurrentHP <= 0 {
				return
			}
			m.CurrentHP += heal
			if m.CurrentHP > m.MaxHP {
				m.CurrentHP = m.MaxHP
			}
		})
	}
}

func tickEffects(effects []StatusEffect, applyDamage func(int), applyHeal func(int)) []StatusEffect {
	var remaining []StatusEffect
	for _, e := range effects {
		if e.DamagePerTurn > 0 {
			applyDamage(e.DamagePerTurn)
		}
		if e.HealPerTurn > 0 {
			applyHeal(e.HealPerTurn)
		}
		if e.TurnsLeft >= PermanentEffectDuration {
			remaining = append(remaining, e)
		} else if e.TurnsLeft > 1 {
			e.TurnsLeft--
			remaining = append(remaining, e)
		}
	}
	return remaining
}

func endFight(p *Character, m *Monster, won bool) {
	if !won {
		return
	}
 
	fmt.Println("\n" + m.Name + " est vaincu !")
 
	gainGold(p, m.GoldReward)
	gainXP(p, m.XPReward)
 
	for _, item := range rollLoot(m) {
		if addInventory(p, item) {
			fmt.Println("Vous récupérez : " + item.Name)
		} else {
			fmt.Println(item.Name + " tombe au sol, votre inventaire est plein.")
		}
	}
 
	if canEvolve(p, m.Tier) {
		before := currentStageName(p)
		evolve(p)
		fmt.Println(onEvolve(p, before))
		storyEvolve(p)
	}
 
	if p.MaxTierDefeated < m.Tier {
		p.MaxTierDefeated = m.Tier
	}
 
	p.CurrentHP = p.MaxHP
	p.CurrentMana = p.MaxMana
}
