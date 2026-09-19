package main
 
import "strconv"

func allSpells() []Spell {
	return []Spell{
		// --- Phénix (Feu) ---
		{
			Name: "Picorage enflammé", Element: ElementFeu,
			MinStage: 2, ManaCost: 8,
			Damage: 10, DotDamage: 3, DotDuration: 2,
		},
		{
			Name: "Plumage chaud", Element: ElementFeu,
			MinStage: 2, ManaCost: 6,
			Heal: 12,
		},
		{
			Name: "Coup de bec ardent", Element: ElementFeu,
			MinStage: 3, ManaCost: 12,
			Damage: 22,
		},
		{
			Name: "Cocorico enflammé", Element: ElementFeu,
			MinStage: 3, ManaCost: 10,
			AttackBuff: 8, BuffDuration: 3,
		},
		{
			Name: "Griffes de braise", Element: ElementFeu,
			MinStage: 4, ManaCost: 15,
			Damage: 32,
		},
		{
			Name: "Nova solaire", Element: ElementFeu,
			MinStage: 4, ManaCost: 20,
			Damage: 18, DotDamage: 6, DotDuration: 4,
		},
		{
			Name: SortRenaissance, Element: ElementFeu,
			MinStage: 4, ManaCost: 0,
			IsUltimate: true,
			Description: "Déclenché automatiquement à 0 PV, pas via le menu de sorts : " +
				"voir combat.go, gestion de la mort du joueur.",
		},
 
		// --- Kraken (Eau) ---
		{
			Name: "Éclaboussure", Element: ElementEau,
			MinStage: 2, ManaCost: 6,
			Damage: 8, InitiativeDebuff: 5, DebuffDuration: 2,
		},
		{
			Name: "Écailles naissantes", Element: ElementEau,
			MinStage: 2, ManaCost: 8,
			DamageReduction: 20, ReductionDuration: 3,
		},
		{
			Name: "Coup de nageoire lourd", Element: ElementEau,
			MinStage: 3, ManaCost: 10,
			Damage: 16, AttackDebuff: 3, DebuffDuration: 1,
		},
		{
			Name: "Bulle protectrice", Element: ElementEau,
			MinStage: 3, ManaCost: 12,
			ShieldValue: 25,
		},
		{
			Name: "Étreinte tentaculaire", Element: ElementEau,
			MinStage: 4, ManaCost: 18,
			Damage: 20, SkipEnemyTurn: true,
		},
		{
			Name: "Raz-de-marée", Element: ElementEau,
			MinStage: 4, ManaCost: 16,
			Damage: 26, InitiativeDebuff: 8, DebuffDuration: 3,
		},
		{
			Name: "Abysse", Element: ElementEau,
			MinStage: 4, ManaCost: 25,
			DotDamage: 10, DotDuration: 5,
			IsUltimate: true,
		},
 
		// --- Caméléon Spectral (Spectre) ---
		{
			Name: "Morsure furtive", Element: ElementSpectre,
			MinStage: 2, ManaCost: 7,
			Damage: 14,
		},
		{
			Name: "Camouflage", Element: ElementSpectre,
			MinStage: 2, ManaCost: 8,
			Invisible: true, InvisibleBonus: 50,
		},
		{
			Name: "Frappe chromatique", Element: ElementSpectre,
			MinStage: 3, ManaCost: 12,
			Damage: 24,
		},
		{
			Name: "Voile chromatique", Element: ElementSpectre,
			MinStage: 3, ManaCost: 10,
			Invisible: true, InvisibleBonus: 100,
		},
		{
			Name: "Lame spectrale", Element: ElementSpectre,
			MinStage: 4, ManaCost: 15,
			Damage: 28, IgnoreResist: true,
		},
		{
			Name: "Dissolution", Element: ElementSpectre,
			MinStage: 4, ManaCost: 12,
			Intangible: true,
		},
		{
			Name: SortExecutionSpectrale, Element: ElementSpectre,
			MinStage: 4, ManaCost: 30,
			Damage:     25,
			IsUltimate: true,
		},
 
		// --- Sort acheté, indépendant de toute espèce ---
		{
			Name: SortSoufflePrimordial, Element: ElementNeutre,
			MinStage: 0, ManaCost: 14,
			Damage: 20,
		},
	}
}

func spellsForStage(s Species, stage int) []Spell {
	var result []Spell
	for _, spell := range allSpells() {
		if spell.Element == s.Element && spell.MinStage == stage {
			result = append(result, spell)
		}
	}
	return result
}

func spellByName(name string) (Spell, bool) {
	for _, spell := range allSpells() {
		if spell.Name == name {
			return spell, true
		}
	}
	return Spell{}, false
}

func hasSpell(p *Character, name string) bool {
	for _, s := range p.Spells {
		if s.Name == name {
			return true
		}
	}
	return false
}

func learnSpell(p *Character, name string) bool {
	if hasSpell(p, name) {
		return false
	}
	spell, ok := spellByName(name)
	if !ok {
		return false
	}
	p.Spells = append(p.Spells, spell)
	return true
}

func canCastSpell(p *Character, spell Spell) (bool, string) {
	if spell.IsUltimate && p.UltimateUsed {
		return false, "Vous avez déjà utilisé votre ultime ce combat."
	}
	if p.CurrentMana < spell.ManaCost {
		return false, "Mana insuffisant."
	}
	return true, ""
}

func castSpell(p *Character, m *Monster, spell Spell) (string, bool) {
	ok, reason := canCastSpell(p, spell)
	if !ok {
		return reason, false
	}
 
	p.CurrentMana -= spell.ManaCost
	if spell.IsUltimate {
		p.UltimateUsed = true
	}
 
	damage := spell.Damage
 
	// Cas particulier : Exécution spectrale fait 60 dégâts si lancé depuis l'invisibilité
	if spell.Name == SortExecutionSpectrale && p.Invisible {
		damage = 60
	}
 
	// Bonus d'invisibilité générique : consommé par le premier sort
	// offensif lancé après Camouflage ou Voile chromatique.
	if p.Invisible && damage > 0 && spell.Name != SortExecutionSpectrale {
		damage += damage * p.InvisibleBonus / 100
	}
	if damage > 0 {
		p.Invisible = false
		p.InvisibleBonus = 0
	}
 
	message := spell.Name
 
	// Effets offensifs sur le monstre
	if damage > 0 && m != nil {
		dealDamage(damage, spell.Element, m, spell.IgnoreResist)
		message += ": " + strconv.Itoa(damage) + " dégâts"
	}
	if spell.DotDamage > 0 && m != nil {
		m.Effects = append(m.Effects, StatusEffect{
			Name: spell.Name, TurnsLeft: spell.DotDuration,
			DamagePerTurn: spell.DotDamage,
		})
	}
	if spell.SkipEnemyTurn && m != nil {
		m.SkipNextTurn = true
	}
	if spell.InitiativeDebuff > 0 && m != nil {
		m.Effects = append(m.Effects, StatusEffect{
			Name: spell.Name, TurnsLeft: spell.DebuffDuration,
			InitiativeDebuff: spell.InitiativeDebuff,
		})
	}
	if spell.AttackDebuff > 0 && m != nil {
		m.Effects = append(m.Effects, StatusEffect{
			Name: spell.Name, TurnsLeft: spell.DebuffDuration,
			AttackBuff: -spell.AttackDebuff,
		})
	}
 
	// Effets sur le joueur lui-même
	if spell.Heal > 0 {
		p.CurrentHP += spell.Heal
		if p.CurrentHP > p.MaxHP {
			p.CurrentHP = p.MaxHP
		}
	}
	if spell.ShieldValue > 0 {
		p.Shield += spell.ShieldValue
	}
	if spell.DamageReduction > 0 {
		p.Effects = append(p.Effects, StatusEffect{
			Name: spell.Name, TurnsLeft: spell.ReductionDuration,
			DamageReduction: spell.DamageReduction,
		})
	}
	if spell.AttackBuff > 0 {
		p.Effects = append(p.Effects, StatusEffect{
			Name: spell.Name, TurnsLeft: spell.BuffDuration,
			AttackBuff: spell.AttackBuff,
		})
	}
	if spell.Invisible {
		p.Invisible = true
		p.InvisibleBonus = spell.InvisibleBonus
	}
	if spell.Intangible {
		p.Intangible = true
	}
 
	return message, true
}
