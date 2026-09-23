package main
 
import "strconv"

func allSpells() []Spell {
	return []Spell{
		// --- Phénix (Feu) : dégâts sur la durée et auto-soin ---
		{
			Name: "Picorage enflammé", Element: ElementFeu,
			Description: "Dégâts directs et laisse une brûlure qui continue de faire mal.",
			MinStage: 2, ManaCost: 8,
			Damage: 10, DotDamage: 4, DotDuration: 2,
		},
		{
			Name: "Plumage chaud", Element: ElementFeu,
			Description: "Restaure des PV en se réchauffant.",
			MinStage: 2, ManaCost: 8,
			Heal: 24,
		},
		{
			Name: "Coup de bec ardent", Element: ElementFeu,
			Description: "Un coup de bec chargé de flammes, gros dégâts directs.",
			MinStage: 3, ManaCost: 11,
			Damage: 24,
		},
		{
			Name: "Cocorico enflammé", Element: ElementFeu,
			Description: "Booste l'attaque pour les prochains tours.",
			MinStage: 3, ManaCost: 9,
			AttackBuff: 7, BuffDuration: 4,
		},
		{
			Name: "Griffes de braise", Element: ElementFeu,
			Description: "Des serres incandescentes, dégâts très élevés sur une cible unique.",
			MinStage: 4, ManaCost: 15,
			Damage: 34,
		},
		{
			Name: "Nova solaire", Element: ElementFeu,
			Description: "Explosion de chaleur intense, dégâts et brûlure longue durée.",
			MinStage: 4, ManaCost: 19,
			Damage: 16, DotDamage: 7, DotDuration: 4,
		},
		{
			Name: SortRenaissance, Element: ElementFeu,
			Description: "Renaît de ses cendres à 40 % de ses PV avec un regain de puissance. Une fois par combat.",
			MinStage: 4, ManaCost: 0,
			IsUltimate: true,
		},

		// --- Kraken (Eau) : encaissement et contrôle ---
		{
			Name: "Éclaboussure", Element: ElementEau,
			Description: "Une vague vive qui ralentit l'adversaire.",
			MinStage: 2, ManaCost: 6,
			Damage: 12, InitiativeDebuff: 6, DebuffDuration: 3,
		},
		{
			Name: "Écailles naissantes", Element: ElementEau,
			Description: "Durcit la peau, réduit les dégâts subis pendant quelques tours.",
			MinStage: 2, ManaCost: 7,
			DamageReduction: 30, ReductionDuration: 3,
		},
		{
			Name: "Coup de nageoire lourd", Element: ElementEau,
			Description: "Frappe puissante qui affaiblit l'attaque adverse.",
			MinStage: 3, ManaCost: 11,
			Damage: 22, AttackDebuff: 4, DebuffDuration: 3,
		},
		{
			Name: "Bulle protectrice", Element: ElementEau,
			Description: "Un bouclier d'eau qui absorbe les prochains dégâts.",
			MinStage: 3, ManaCost: 11,
			ShieldValue: 35,
		},
		{
			Name: "Étreinte tentaculaire", Element: ElementEau,
			Description: "Immobilise la cible, qui perd son prochain tour.",
			MinStage: 4, ManaCost: 17,
			Damage: 22, SkipEnemyTurn: true,
		},
		{
			Name: "Raz-de-marée", Element: ElementEau,
			Description: "Une vague massive, dégâts élevés et ralentit fortement la cible.",
			MinStage: 4, ManaCost: 16,
			Damage: 28, InitiativeDebuff: 10, DebuffDuration: 3,
		},
		{
			Name: "Abysse", Element: ElementEau,
			Description: "Aspire l'adversaire vers les profondeurs, dégâts continus sur la durée.",
			MinStage: 4, ManaCost: 24,
			DotDamage: 13, DotDuration: 5,
			IsUltimate: true,
		},

		// --- Caméléon Spectral (Spectre) : burst depuis l'invisibilité ---
		{
			Name: "Morsure furtive", Element: ElementSpectre,
			Description: "Une attaque discrète, bien plus forte lancée depuis l'invisibilité.",
			MinStage: 2, ManaCost: 7,
			Damage: 15,
		},
		{
			Name: "Camouflage", Element: ElementSpectre,
			Description: "Devient invisible : esquive accrue, et la prochaine attaque frappe plus de deux fois plus fort.",
			MinStage: 2, ManaCost: 11,
			Invisible: true, InvisibleBonus: 130,
		},
		{
			Name: "Frappe chromatique", Element: ElementSpectre,
			Description: "Un éclair de couleurs, dégâts élevés.",
			MinStage: 3, ManaCost: 13,
			Damage: 30,
		},
		{
			Name: "Voile chromatique", Element: ElementSpectre,
			Description: "Invisibilité renforcée : esquive accrue, et la prochaine attaque frappe près de trois fois plus fort.",
			MinStage: 3, ManaCost: 15,
			Invisible: true, InvisibleBonus: 170,
		},
		{
			Name: "Lame spectrale", Element: ElementSpectre,
			Description: "Frappe à moitié hors du plan, ignore les résistances de la cible.",
			MinStage: 4, ManaCost: 15,
			Damage: 32, IgnoreResist: true,
		},
		{
			Name: SortExecutionSpectrale, Element: ElementSpectre,
			Description: "Frappe fatale qui ignore les résistances, dévastatrice depuis l'invisibilité.",
			MinStage: 4, ManaCost: 20,
			Damage: 35, IgnoreResist: true,
			IsUltimate: true,
		},

		// --- Sort acheté, indépendant de toute espèce ---
		{
			Name: SortSoufflePrimordial, Element: ElementNeutre,
			Description: "Un souffle d'énergie brute : jamais résisté. Le recours quand l'élément ne passe pas.",
			MinStage: 0, ManaCost: 13,
			Damage: 24,
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

	// Cas particulier : Exécution spectrale frappe bien plus fort depuis l'invisibilité.
	if spell.Name == SortExecutionSpectrale && p.Invisible {
		damage = ExecutionFurtiveDmg
	}

	// Bonus/malus d'attaque en cours (Cocorico enflammé, Renaissance...).
	if damage > 0 {
		damage += effectiveAttackBonus(p.Effects)
		if damage < 0 {
			damage = 0
		}
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
		applied := dealDamage(damage, spell.Element, m, spell.IgnoreResist)
		message += ": " + strconv.Itoa(applied) + " dégâts"
	}
	if spell.DotDamage > 0 && m != nil {
		m.Effects = applyEffect(m.Effects, StatusEffect{
			Name: spell.Name, TurnsLeft: spell.DotDuration,
			DamagePerTurn: spell.DotDamage,
		})
	}
	if spell.SkipEnemyTurn && m != nil {
		if m.StunImmune {
			message += " (" + m.Name + " résiste à l'immobilisation)"
		} else {
			m.SkipNextTurn = true
		}
	}
	if spell.InitiativeDebuff > 0 && m != nil {
		m.Effects = applyEffect(m.Effects, StatusEffect{
			Name: spell.Name, TurnsLeft: spell.DebuffDuration,
			InitiativeDebuff: spell.InitiativeDebuff,
		})
	}
	if spell.AttackDebuff > 0 && m != nil {
		m.Effects = applyEffect(m.Effects, StatusEffect{
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
		p.Effects = applyEffect(p.Effects, StatusEffect{
			Name: spell.Name, TurnsLeft: spell.ReductionDuration,
			DamageReduction: spell.DamageReduction,
		})
	}
	if spell.AttackBuff > 0 {
		p.Effects = applyEffect(p.Effects, StatusEffect{
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
