package main

func allSpecies() []Species {
	return []Species{
		{
			Name:    EspecePhenix,
			Element: ElementFeu,
			StageNames: []string{
				"Œuf ardent", "Poussin", "Poulet", "Phénix",
			},
			BaseMaxHP:      50,
			BaseMaxMana:    44,
			BaseAttack:     6,
			HPPerStage:     20,
			ManaPerStage:   18,
			AttackPerStage: 3,
			Initiative:     12,
			DodgeRate:      5,
			// Triangle : le Feu bat le Spectre, l'Eau bat le Feu.
			Resistances: map[Element]float64{
				ElementFeu:     0.6,
				ElementEau:     1.35,
				ElementSpectre: 1.0,
				ElementNeutre:  1.0,
			},
		},
		{
			Name:    EspeceKraken,
			Element: ElementEau,
			StageNames: []string{
				"Œuf immergé", "Alevin", "Poisson", "Kraken",
			},
			BaseMaxHP:      70,
			BaseMaxMana:    36,
			BaseAttack:     6,
			HPPerStage:     25,
			ManaPerStage:   17,
			AttackPerStage: 2,
			Initiative:     8,
			DodgeRate:      5,
			// Le Spectre bat l'Eau, l'Eau bat le Feu.
			Resistances: map[Element]float64{
				ElementFeu:     1.0,
				ElementEau:     0.6,
				ElementSpectre: 1.35,
				ElementNeutre:  1.0,
			},
		},
		{
			Name:    EspeceSpectre,
			Element: ElementSpectre,
			StageNames: []string{
				"Œuf voilé", "Caméléonneau", "Caméléon", "Caméléon Spectral",
			},
			BaseMaxHP:      48,
			BaseMaxMana:    32,
			BaseAttack:     6,
			HPPerStage:     16,
			ManaPerStage:   15,
			AttackPerStage: 4,
			Initiative:     18,
			DodgeRate:      20,
			// Le Feu bat le Spectre, le Spectre bat l'Eau.
			Resistances: map[Element]float64{
				ElementFeu:     1.35,
				ElementEau:     1.0,
				ElementSpectre: 0.6,
				ElementNeutre:  1.0,
			},
		},
	}
}

func speciesByName(name string) Species {
	for _, s := range allSpecies() {
		if s.Name == name {
			return s
		}
	}
	return allSpecies()[0]
}

func canEvolve(p *Character, defeatedTier int) bool {
	if defeatedTier <= p.MaxTierDefeated {
		return false
	}
	if p.Stage >= 4 {
		return false
	}
	return true
}

func evolve(p *Character) {
	p.Stage++
	recalculateStats(p)
	p.CurrentHP = p.MaxHP
	p.CurrentMana = p.MaxMana
	for _, spell := range spellsForStage(p.Species, p.Stage) {
		learnSpell(p, spell.Name)
	}
}

func currentStageName(p *Character) string {
	return p.Species.StageNames[p.Stage-1]
}

func getResistance(s Species, e Element) float64 {
	if mult, ok := s.Resistances[e]; ok {
		return mult
	}
	return 1.0
}

func isGameComplete(p *Character) bool {
	return p.MaxTierDefeated >= 4
}

func onEvolve(p *Character, previousStageName string) string {
	return previousStageName + " devient... " + currentStageName(p) + " !"
}
