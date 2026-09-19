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
			BaseMaxMana:    40,
			BaseAttack:     5,
			HPPerStage:     18,
			ManaPerStage:   15,
			AttackPerStage: 3,
			Initiative:     12,
			DodgeRate:      5,
			Resistances: map[Element]float64{
				ElementFeu:     0.6,
				ElementEau:     1.4,
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
			BaseMaxHP:      75,
			BaseMaxMana:    30,
			BaseAttack:     5,
			HPPerStage:     25,
			ManaPerStage:   12,
			AttackPerStage: 2,
			Initiative:     6,
			DodgeRate:      0,
			Resistances: map[Element]float64{
				ElementFeu:     0.7,
				ElementEau:     0.6,
				ElementSpectre: 1.2,
				ElementNeutre:  0.9,
			},
		},
		{
			Name:    EspeceSpectre,
			Element: ElementSpectre,
			StageNames: []string{
				"Œuf voilé", "Caméléonneau", "Caméléon", "Caméléon Spectral",
			},
			BaseMaxHP:      45,
			BaseMaxMana:    35,
			BaseAttack:     5,
			HPPerStage:     14,
			ManaPerStage:   15,
			AttackPerStage: 4,
			Initiative:     18,
			DodgeRate:      15,
			Resistances: map[Element]float64{
				ElementFeu:     1.2,
				ElementEau:     1.0,
				ElementSpectre: 0.5,
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
