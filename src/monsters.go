package main
 
import "math/rand"

func initVer() *Monster {
	return &Monster{
		Name:       MonstreVer,
		Tier:       1,
		MaxHP:      20,
		CurrentHP:  20,
		Attack:     3,
		Initiative: 3,
		Element:    ElementNeutre,
		Skills:      []string{"Morsure"},
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
		Skills:      []string{"Charge", "Enragé"},
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
		Skills:      []string{"Griffure enragée", "Hurlement", "Régénération lunaire"},
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
				Attacks: []string{"Souffle ardent", "Griffe"},
			},
			{
				Name: "Vague Ancestrale", HPThreshold: 0.75, Element: ElementEau,
				Resistances: map[Element]float64{
					ElementFeu: 1.5, ElementEau: 0.5, ElementSpectre: 1.0, ElementNeutre: 1.0,
				},
				Attacks: []string{"Raz-de-marée", "Étreinte liquide"},
			},
			{
				Name: "Voile Spectral", HPThreshold: 0.50, Element: ElementSpectre,
				Resistances: map[Element]float64{
					ElementFeu: 1.0, ElementEau: 1.5, ElementSpectre: 0.5, ElementNeutre: 1.0,
				},
				Attacks: []string{"Frappe hors du plan", "Dissolution"},
			},
			{
				Name: "Convergence", HPThreshold: 0.25, Element: ElementNeutre,
				Resistances: map[Element]float64{
					ElementFeu: 1.0, ElementEau: 1.0, ElementSpectre: 1.0, ElementNeutre: 1.0,
				},
				DamageBonus: 50,
				Attacks:     []string{"Souffle ardent", "Raz-de-marée", "Frappe hors du plan"},
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
