package main

import "fmt"

// 100 XP au niveau 1, puis +50 par niveau
func xpForNextLevel(level int) int {
	return 100 + (level-1)*50
}

// partagée : appelée à la fin d'un combat
func gainXP(p *Character, amount int) {
	if p.MaxExperience <= 0 {
		p.MaxExperience = xpForNextLevel(p.Level)
	}
	fmt.Printf("%s gagne %d XP\n", p.Name, amount)
	p.Experience += amount
	for p.Experience >= p.MaxExperience {
		levelUp(p)
	}
	fmt.Printf("XP : %d / %d\n", p.Experience, p.MaxExperience)
}

// partagée : loot et récompenses
func gainGold(p *Character, amount int) {
	p.Gold += amount
	fmt.Printf("+%d or (total : %d)\n", amount, p.Gold)
}

func levelUp(p *Character) {
	p.Experience -= p.MaxExperience
	p.Level++
	p.MaxExperience = xpForNextLevel(p.Level)
	recalculateStats(p) // recalcule MaxHP/MaxMana/Attack depuis zéro, pas de +=
	p.CurrentHP = p.MaxHP
	p.CurrentMana = p.MaxMana
	fmt.Printf("NIVEAU %d ! +%d PV max, +%d mana max, +%d attaque\n",
		p.Level, GainPVNiveau, GainManaNiveau, GainAtkNiveau)
}
