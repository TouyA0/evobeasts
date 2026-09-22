package main
 
import (
	"fmt"
	"math/rand"
)

const (
	sanglierChargePct   = 130
	sanglierStunChance  = 30
	sanglierEnrage50Pct = 130
	sanglierEnrage25Pct = 160

	loupGarouHowlBonusAtk = 5
	loupGarouRegenHeal    = 6
	loupGarouClawPct      = 70

	dragonSouffleArdentPct   = 120
	dragonBrulureDmg         = 5
	dragonBrulureDuree       = 2
	dragonGriffePct          = 100
	dragonRegenAbyssale      = 8
	dragonEtreinteLiquidePct = 80
	dragonRazDeMareePct      = 120
	dragonFrappeHorsPlanPct  = 130
	dragonConvergencePct     = 150
)

func initVer() *Monster {
	return &Monster{
		Name:       MonstreVer,
		Tier:       1,
		MaxHP:      20,
		CurrentHP:  20,
		Attack:     3,
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
		GoldReward: 5,
		LootTable: []LootEntry{
			{Item: Item{Name: ResPlume, Type: ItemResource, Price: 1}, Chance: 60},
		},
	}
}
 
func initSanglier() *Monster {
	return &Monster{
		Name:       MonstreSanglier,
		Tier:       2,
		MaxHP:      45,
		CurrentHP:  45,
		Attack:     7,
		Initiative: 9,
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
		Resistances: map[Element]float64{
			ElementFeu: 1.0, ElementEau: 1.0, ElementSpectre: 1.0, ElementNeutre: 1.0,
		},
		XPReward:   90,
		GoldReward: 12,
		LootTable: []LootEntry{
			{Item: Item{Name: ResEcaille, Type: ItemResource, Price: 4}, Chance: 60},
		},
	}
}
 
func initLoupGarou() *Monster {
	return &Monster{
		Name:       MonstreLoupGarou,
		Tier:       3,
		MaxHP:      80,
		CurrentHP:  80,
		Attack:     12,
		Initiative: 13,
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
		GoldReward: 25,
		LootTable: []LootEntry{
			{Item: Item{Name: ResSpectre, Type: ItemResource, Price: 3}, Chance: 50},
			{Item: Item{Name: ResEcaille, Type: ItemResource, Price: 4}, Chance: 40},
		},
	}
}

func initDragonAncien() *Monster {
	return &Monster{
		Name:       MonstreDragonAncien,
		Tier:       4,
		MaxHP:      150,
		CurrentHP:  150,
		Attack:     16,
		Initiative: 15,
		Element:    ElementFeu, // élément de la phase de départ
		Resistances: map[Element]float64{
			ElementFeu: 1.0, ElementEau: 1.0, ElementSpectre: 1.0, ElementNeutre: 1.0,
		},
		XPReward:   300,
		GoldReward: 80,
		LootTable: []LootEntry{
			{Item: Item{Name: ResAbyssal, Type: ItemResource, Price: 7}, Chance: 100},
		},
		CurrentPhase: 0,
		Phases: []Phase{
			{
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
				Name: "Vague Ancestrale", HPThreshold: 0.75, Element: ElementEau,
				Resistances: map[Element]float64{
					ElementFeu: 1.5, ElementEau: 0.5, ElementSpectre: 1.0, ElementNeutre: 1.0,
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
				Name: "Voile Spectral", HPThreshold: 0.50, Element: ElementSpectre,
				Resistances: map[Element]float64{
					ElementFeu: 1.0, ElementEau: 1.5, ElementSpectre: 0.5, ElementNeutre: 1.0,
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
				DamageBonus: 50,
				Attacks: []MonsterSkill{
					{
						Name: "Convergence élémentaire", Kind: SkillAttack, DamagePct: dragonConvergencePct,
						Description: fmt.Sprintf("%d%% dégâts, élément change chaque tour (Feu→Eau→Spectre), +50%% dégâts de phase.",
							dragonConvergencePct),
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
