package main
 
import (
	"fmt"
	"math/rand"
)

const (
	sanglierChargePct   = 140
	sanglierStunChance  = 35
	sanglierEnrage50Pct = 130
	sanglierEnrage25Pct = 160

	loupGarouHowlBonusAtk = 4
	loupGarouRegenHeal    = 5
	loupGarouClawPct      = 65

	dragonSouffleArdentPct   = 125
	dragonBrulureDmg         = 6
	dragonBrulureDuree       = 2
	dragonGriffePct          = 105
	dragonRegenAbyssale      = 6
	dragonEtreinteLiquidePct = 90
	dragonRazDeMareePct      = 115
	dragonFrappeHorsPlanPct  = 135
	dragonConvergencePct     = 135
	dragonConvergenceBonus   = 30 // bonus de dégâts de la phase finale
)

func initVer() *Monster {
	return &Monster{
		Name:       MonstreVer,
		Tier:       1,
		MaxHP:      15,
		CurrentHP:  15,
		Attack:     4,
		Initiative: 3,
		Element:    ElementNeutre,
		Skills: []MonsterSkill{
			{
				Name: "Morsure", Kind: SkillAttack, DamagePct: 100,
				Description: "Attaque simple à chaque tour (100% ATK).",
			},
		},
		Resistances: map[Element]float64{
			ElementFeu: 1.0, ElementEau: 1.0, ElementSpectre: 1.0, ElementNeutre: 1.0,
		},
		XPReward:   40,
		GoldReward: 6,
		LootTable: []LootEntry{
			{Item: Item{Name: ResPlume, Type: ItemResource, Price: PrixPlume}, Chance: 70},
		},
	}
}
 
func initSanglier() *Monster {
	m := &Monster{
		Name:       MonstreSanglier,
		Tier:       2,
		MaxHP:      55,
		CurrentHP:  55,
		Attack:     7,
		Initiative: 10,
		Element:    ElementNeutre,
		Skills: []MonsterSkill{
			{
				Name: "Charge", Kind: SkillAttack, DamagePct: sanglierChargePct,
				Description: fmt.Sprintf("Tous les 3 tours : %d%% dégâts, %d%% de chance d'étourdir.",
					sanglierChargePct, sanglierStunChance),
			},
			{
				Name: "Enragé", Kind: SkillBuff,
				Description: fmt.Sprintf("Passif : dégâts x%.2f sous 50%% PV, x%.2f sous 25%% PV.",
					float64(sanglierEnrage50Pct)/100, float64(sanglierEnrage25Pct)/100),
			},
		},
		XPReward:   90,
		GoldReward: 15,
		LootTable: []LootEntry{
			{Item: Item{Name: ResEcaille, Type: ItemResource, Price: PrixEcaille}, Chance: 75},
		},
	}
	m.Resistances = map[Element]float64{
		ElementFeu: 1.0, ElementEau: 1.0, ElementSpectre: 1.0, ElementNeutre: 1.0,
	}
	return m
}
 
func initLoupGarou() *Monster {
	return &Monster{
		Name:       MonstreLoupGarou,
		Tier:       3,
		MaxHP:      90,
		CurrentHP:  90,
		Attack:     11,
		Initiative: 14,
		Element:    ElementNeutre,
		Skills: []MonsterSkill{
			{
				Name: "Griffure enragée", Kind: SkillAttack, DamagePct: loupGarouClawPct,
				Description: fmt.Sprintf("Tours pairs : 2 attaques à %d%% ATK chacune.", loupGarouClawPct),
			},
			{
				Name: "Hurlement", Kind: SkillBuff,
				Description: fmt.Sprintf("Au tour 3 (une fois) : +%d ATK permanent.", loupGarouHowlBonusAtk),
			},
			{
				Name: "Régénération lunaire", Kind: SkillHeal,
				Description: fmt.Sprintf("Passif dès le début : soigne %d PV/tour, indéfiniment.", loupGarouRegenHeal),
			},
		},
		Resistances: map[Element]float64{
			ElementFeu: 1.0, ElementEau: 1.0, ElementSpectre: 1.0, ElementNeutre: 1.0,
		},
		XPReward:   150,
		GoldReward: 30,
		LootTable: []LootEntry{
			{Item: Item{Name: ResSpectre, Type: ItemResource, Price: PrixSpectre}, Chance: 70},
			{Item: Item{Name: ResEcaille, Type: ItemResource, Price: PrixEcaille}, Chance: 50},
		},
	}
}

func initDragonAncien() *Monster {
	return &Monster{
		Name:       MonstreDragonAncien,
		Tier:       4,
		MaxHP:      230,
		CurrentHP:  230,
		Attack:     17,
		Initiative: 16,
		Element:    ElementFeu,
		Resistances: map[Element]float64{
			ElementFeu: 1.0, ElementEau: 1.0, ElementSpectre: 1.0, ElementNeutre: 1.0,
		},
		XPReward:   300,
		GoldReward: 80,
		LootTable: []LootEntry{
			{Item: Item{Name: ResAbyssal, Type: ItemResource, Price: PrixAbyssal}, Chance: 100},
		},
		CurrentPhase: 0,
		Phases: []Phase{
			{
				// Le dragon est Feu : il résiste au Feu, il craint l'Eau.
				Name: "Éveil Igné", HPThreshold: 1.0, Element: ElementFeu,
				Resistances: map[Element]float64{
					ElementFeu: 0.5, ElementEau: 1.5, ElementSpectre: 1.0, ElementNeutre: 1.0,
				},
				Attacks: []MonsterSkill{
					{
						Name: "Souffle ardent", Kind: SkillAttack, Element: ElementFeu, DamagePct: dragonSouffleArdentPct,
						Description: fmt.Sprintf("Tours pairs : %d%% dégâts + brûlure %d dégâts/tour (%d tours).",
							dragonSouffleArdentPct, dragonBrulureDmg, dragonBrulureDuree),
					},
					{
						Name: "Griffe", Kind: SkillAttack, DamagePct: dragonGriffePct,
						Description: fmt.Sprintf("Tours impairs : %d%% dégâts.", dragonGriffePct),
					},
				},
			},
			{
				// Le dragon est Eau : il résiste à l'Eau, il craint le Spectre.
				Name: "Vague Ancestrale", HPThreshold: 0.75, Element: ElementEau,
				Resistances: map[Element]float64{
					ElementFeu: 1.0, ElementEau: 0.5, ElementSpectre: 1.5, ElementNeutre: 1.0,
				},
				Attacks: []MonsterSkill{
					{
						Name: "Étreinte liquide", Kind: SkillAttack, Element: ElementEau, DamagePct: dragonEtreinteLiquidePct,
						Description: fmt.Sprintf("Tous les 3 tours : %d%% dégâts + étourdissement.", dragonEtreinteLiquidePct),
					},
					{
						Name: "Raz-de-marée", Kind: SkillAttack, Element: ElementEau, DamagePct: dragonRazDeMareePct,
						Description: fmt.Sprintf("Sinon : %d%% dégâts.", dragonRazDeMareePct),
					},
					{
						Name: "Régénération abyssale", Kind: SkillHeal,
						Description: fmt.Sprintf("Passif : soigne %d PV/tour.", dragonRegenAbyssale),
					},
				},
			},
			{
				// Le dragon est Spectre : il résiste au Spectre, il craint le Feu.
				Name: "Voile Spectral", HPThreshold: 0.50, Element: ElementSpectre,
				Resistances: map[Element]float64{
					ElementFeu: 1.5, ElementEau: 1.0, ElementSpectre: 0.5, ElementNeutre: 1.0,
				},
				Attacks: []MonsterSkill{
					{
						Name: "Dissolution", Kind: SkillBuff, Element: ElementSpectre,
						Description: "Tous les 4 tours : esquive tout ce tour.",
					},
					{
						Name: "Frappe hors du plan", Kind: SkillAttack, Element: ElementSpectre, DamagePct: dragonFrappeHorsPlanPct,
						Description: fmt.Sprintf("Sinon : %d%% dégâts, ignore le bouclier.", dragonFrappeHorsPlanPct),
					},
				},
			},
			{
				Name: "Convergence", HPThreshold: 0.25, Element: ElementNeutre,
				Resistances: map[Element]float64{
					ElementFeu: 1.0, ElementEau: 1.0, ElementSpectre: 1.0, ElementNeutre: 1.0,
				},
				DamageBonus: dragonConvergenceBonus,
				Attacks: []MonsterSkill{
					{
						Name: "Convergence élémentaire", Kind: SkillAttack, DamagePct: dragonConvergencePct,
						Description: fmt.Sprintf("%d%% dégâts, élément change chaque tour (Feu→Eau→Spectre), +%d%% dégâts de phase.",
							dragonConvergencePct, dragonConvergenceBonus),
					},
				},
			},
		},
	}
}

func monsterByTier(tier int) *Monster {
	switch tier {
	case 1:
		return initVer()
	case 2:
		return initSanglier()
	case 3:
		return initLoupGarou()
	case 4:
		return initDragonAncien()
	default:
		return initVer()
	}
}

func updatePhase(m *Monster) string {
	if len(m.Phases) == 0 {
		return ""
	}
	fraction := float64(m.CurrentHP) / float64(m.MaxHP)
	newPhase := 0
	for i := len(m.Phases) - 1; i >= 0; i-- {
		if fraction <= m.Phases[i].HPThreshold {
			newPhase = i
			break
		}
	}
	if newPhase != m.CurrentPhase {
		m.CurrentPhase = newPhase
		m.Element = m.Phases[newPhase].Element
		m.Resistances = m.Phases[newPhase].Resistances
		return m.Name + " entre en phase : " + m.Phases[newPhase].Name + " !"
	}
	return ""
}

func rollLoot(m *Monster) []Item {
	var loot []Item
	for _, entry := range m.LootTable {
		if rand.Intn(100) < entry.Chance {
			loot = append(loot, entry.Item)
		}
	}
	return loot
}
