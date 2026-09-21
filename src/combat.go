package main
 
import (
	"fmt"
	"math/rand"
	"strconv"
)

func combatMenu(p *Character) {
	for {
		fmt.Println("\n--- Combat ---")
		unlocked := p.MaxTierDefeated + 1
		if unlocked > 4 {
			unlocked = 4
		}
		names := []string{MonstreVer, MonstreSanglier, MonstreLoupGarou, MonstreDragonAncien}
		for i := 1; i <= unlocked; i++ {
			fmt.Printf("%d. %s\n", i, names[i-1])
		}
		fmt.Println("0. Retour")
		fmt.Print("Votre choix : ")
 
		choice := readInput()
		if choice == "0" {
			return
		}
 
		tier, err := strconv.Atoi(choice)
		if err != nil || tier < 1 || tier > unlocked {
			fmt.Println("Choix invalide.")
			continue
		}
 
		monster := monsterByTier(tier)
		trainingFight(p, monster)
 
		if isGameComplete(p) {
			fmt.Println("\n*** Le Dragon Ancien est vaincu. EvoBeasts est terminé. ***")
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
 
	turn := 1
 
	for {
		combatScreen(p, m, turn)
 
		playerFirst := effectiveInitiative(p.Initiative, p.Effects) >=
			effectiveInitiative(m.Initiative, m.Effects)
 
		if playerFirst {
			runCharacterTurn(p, m)
			if m.CurrentHP <= 0 {
				endFight(p, m, true)
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
				return
			}
		}
 
		resolveEffects(p, m)
		pauseForInput()
		turn++
	}
}

func pauseForInput() {
	fmt.Print("\n(Appuyez sur Entrée pour continuer...)")
	readInput()
}

func checkPlayerDefeat(p *Character) bool {
	if p.CurrentHP > 0 {
		return false
	}
 
	if hasSpell(p, SortRenaissance) && !p.UltimateUsed {
		p.UltimateUsed = true
		p.CurrentHP = p.MaxHP / 2
		p.Effects = append(p.Effects, StatusEffect{
			Name: SortRenaissance, TurnsLeft: 3, AttackBuff: 10,
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
	if p.SkipNextTurn {
		fmt.Println(p.Name + " ne peut pas agir ce tour-ci.")
		p.SkipNextTurn = false
		return
	}
 
	fmt.Println("\n1. Attaquer   2. Sorts   3. Inventaire")
	fmt.Print("Votre choix : ")
	choice := readInput()
 
	switch choice {
	case "1":
		basicAttack(p, m)
	case "2":
		spellMenu(p, m)
	case "3":
		combatInventoryMenu(p, m)
	default:
		fmt.Println("Choix invalide, tour passé.")
	}
}

func basicAttack(p *Character, m *Monster) {
	dmg := dealDamage(p.Attack, p.Species.Element, m, false)
	fmt.Printf("%s inflige %d dégâts à %s\n", p.Name, dmg, m.Name)
	fmt.Printf("%s  PV : %d/%d\n", m.Name, m.CurrentHP, m.MaxHP)
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

func spellMenu(p *Character, m *Monster) {
	castable := make([]Spell, 0, len(p.Spells))
	for _, s := range p.Spells {
		if s.Name != SortRenaissance {
			castable = append(castable, s)
		}
	}

	if len(castable) == 0 {
		fmt.Println("Aucun sort connu.")
		return
	}

	for {
		for i, s := range castable {
			tag := ""
			if s.Element != p.Species.Element {
				tag += fmt.Sprintf(" [%s]", s.Element)
			}
			if s.IsUltimate {
				tag += " (ULTIME)"
			}

			ok, reason := canCastSpell(p, s)
			if ok {
				fmt.Printf("%d. %s %s%s (%d mana)\n", i+1, spellIcon(s), s.Name, tag, s.ManaCost)
			} else {
				fmt.Printf("\x1b[38;2;220;70;70m%d. %s %s%s (%d mana) — %s\x1b[0m\n",
					i+1, spellIcon(s), s.Name, tag, s.ManaCost, reason)
			}
			fmt.Printf("   %s\n", s.Description)
		}
		fmt.Println("0. Retour")
		fmt.Print("Votre choix : ")

		choice := readInput()
		if choice == "0" {
			return
		}
		idx, err := strconv.Atoi(choice)
		if err != nil || idx < 1 || idx > len(castable) {
			fmt.Println("Choix invalide, tour passé.")
			return
		}

		spell := castable[idx-1]
		if ok, reason := canCastSpell(p, spell); !ok {
			fmt.Println("Sort inutilisable : " + reason)
			continue
		}

		message, ok := castSpell(p, m, spell)
		fmt.Println(message)
		if !ok {
			continue
		}
		if m != nil {
			fmt.Printf("%s  PV : %d/%d\n", m.Name, m.CurrentHP, m.MaxHP)
		}
		announcePhaseChange(m)
		return
	}
}

func combatInventoryMenu(p *Character, m *Monster) {
	if len(p.Inventory) == 0 {
		fmt.Println("Inventaire vide.")
		return
	}
	for i, item := range p.Inventory {
		fmt.Printf("%d. %s\n", i+1, item.Name)
	}
	fmt.Println("0. Retour")
	fmt.Print("Votre choix : ")
	choice := readInput()
	if choice == "0" {
		return
	}
	idx, err := strconv.Atoi(choice)
	if err != nil || idx < 1 || idx > len(p.Inventory) {
		fmt.Println("Choix invalide, tour passé.")
		return
	}
 
	item := p.Inventory[idx-1]
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
		m.Effects = append(m.Effects, StatusEffect{
			Name: item.Name, TurnsLeft: item.PoisonTicks,
			DamagePerTurn: item.PoisonDamage,
		})
		fmt.Println(item.Name + " est jetée sur " + m.Name + " !")
	}
 
	removeInventory(p, item.Name)
}

func announcePhaseChange(m *Monster) {
	if m == nil {
		return
	}
	if msg := updatePhase(m); msg != "" {
		fmt.Println(msg)
	}
}

func runMonsterTurn(p *Character, m *Monster, turn int) {
	if m.SkipNextTurn {
		fmt.Println(m.Name + " ne peut pas agir ce tour-ci.")
		m.SkipNextTurn = false
		return
	}
 
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

func monsterBasicAttack(p *Character, m *Monster, percent int) {
	base := m.Attack * percent / 100
	dmg, dodged := dealDamageToPlayer(base, m.Element, p)
	if dodged {
		fmt.Printf("%s attaque, %s esquive !\n", m.Name, p.Name)
		return
	}
	fmt.Printf("%s inflige %d dégâts à %s\n", m.Name, dmg, p.Name)
	fmt.Printf("%s  PV : %d/%d\n", p.Name, p.CurrentHP, p.MaxHP)
}

func verPattern(p *Character, m *Monster) {
	monsterBasicAttack(p, m, 100)
}

func sanglierPattern(p *Character, m *Monster, turn int) {
	percent := 100
	if turn%3 == 0 {
		percent = 130
		fmt.Println(m.Name + " charge !")
		if rand.Intn(100) < 30 {
			p.SkipNextTurn = true
			fmt.Println(p.Name + " est étourdi pour le prochain tour !")
		}
	}
 
	fraction := float64(m.CurrentHP) / float64(m.MaxHP)
	if fraction <= 0.25 {
		percent = percent * 160 / 100
	} else if fraction <= 0.50 {
		percent = percent * 130 / 100
	}
 
	monsterBasicAttack(p, m, percent)
}

func loupGarouPattern(p *Character, m *Monster, turn int) {
	if !m.HowlUsed && turn == 3 {
		m.Attack += 5
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
		m.Effects = append(m.Effects, StatusEffect{
			Name: "Régénération lunaire", TurnsLeft: 9999, HealPerTurn: 6,
		})
	}
 
	if turn%2 == 0 {
		fmt.Println(m.Name + " lacère avec fureur !")
		monsterBasicAttack(p, m, 70)
		if p.CurrentHP > 0 {
			monsterBasicAttack(p, m, 70)
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
			monsterBasicAttack(p, m, 120)
			p.Effects = append(p.Effects, StatusEffect{
				Name: "Brûlure du Dragon", TurnsLeft: 2, DamagePerTurn: 5,
			})
		} else {
			monsterBasicAttack(p, m, 100)
		}
 
	case 1: // Vague Ancestrale
		hasRegen := false
		for _, e := range m.Effects {
			if e.Name == "Régénération abyssale" {
				hasRegen = true
			}
		}
		if !hasRegen {
			m.Effects = append(m.Effects, StatusEffect{
				Name: "Régénération abyssale", TurnsLeft: 9999, HealPerTurn: 8,
			})
		}
		if turn%3 == 0 {
			fmt.Println(m.Name + " vous enveloppe d'une étreinte liquide !")
			monsterBasicAttack(p, m, 80)
			p.SkipNextTurn = true
		} else {
			fmt.Println(m.Name + " déchaîne un raz-de-marée !")
			monsterBasicAttack(p, m, 120)
		}
 
	case 2: // Voile Spectral
		if turn%4 == 0 {
			fmt.Println(m.Name + " se dissout hors du plan et esquive tout ce tour.")
			m.Intangible = true
		} else {
			fmt.Println(m.Name + " frappe hors du plan, ignorant tout bouclier !")
			base := m.Attack * 130 / 100
			p.Shield = 0 // la frappe ignore le bouclier
			dmg, dodged := dealDamageToPlayer(base, m.Element, p)
			if dodged {
				fmt.Println(p.Name + " esquive !")
			} else {
				fmt.Printf("%s inflige %d dégâts à %s\n", m.Name, dmg, p.Name)
				fmt.Printf("%s  PV : %d/%d\n", p.Name, p.CurrentHP, p.MaxHP)
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
		base := m.Attack * 150 / 100
		dmg, dodged := dealDamageToPlayerElement(base, element, p)
		if dodged {
			fmt.Println(p.Name + " esquive !")
		} else {
			fmt.Printf("%s inflige %d dégâts à %s\n", m.Name, dmg, p.Name)
			fmt.Printf("%s  PV : %d/%d\n", p.Name, p.CurrentHP, p.MaxHP)
		}
	}
}

func monsterResistance(m *Monster, e Element) float64 {
	if mult, ok := m.Resistances[e]; ok {
		return mult
	}
	return 1.0
}

func dealDamage(base int, element Element, target *Monster, ignoreResist bool) int {
	mult := 1.0
	if !ignoreResist {
		mult = monsterResistance(target, element)
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
	if rand.Intn(100) < p.DodgeRate {
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
		if e.TurnsLeft > 1 {
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
	}
 
	if p.MaxTierDefeated < m.Tier {
		p.MaxTierDefeated = m.Tier
	}
 
	p.CurrentHP = p.MaxHP
	p.CurrentMana = p.MaxMana
}
