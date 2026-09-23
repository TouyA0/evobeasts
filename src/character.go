package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)
func readInput() string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func formatName(s string) string {
	s = strings.ToLower(s)
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func isLettersOnly(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return false
		}
	}
	return true
}

func characterCreation() *Character {
	fmt.Println(storyNamePrompt())
	name := readInput()
	for !isLettersOnly(name) {
		fmt.Println("Le nom ne doit contenir que des lettres. Réessayez :")
		name = readInput()
	}
	name = formatName(name)

	fmt.Println("Choisissez votre œuf :")
	for _, l := range storyEggLines() {
		fmt.Println(l)
	}
	choice := readInput()

	for choice != "1" && choice != "2" && choice != "3" {
		fmt.Println("Veuillez entrer un nombre valide (1, 2 ou 3) :")
		choice = readInput()
	}

	var speciesName string
	switch choice {
	case "1":
		speciesName = EspecePhenix
	case "2":
		speciesName = EspeceKraken
	case "3":
		speciesName = EspeceSpectre
	}

	species := speciesByName(speciesName)

	p := initCharacter(name, species)
	storyHatch(p)
	return p
}

func initCharacter(name string, species Species) *Character {
	p := &Character{
		Name:    name,
		Species: species,
		Stage:   1,
		Level:   1,
	}

	recalculateStats(p)

	p.CurrentHP = p.MaxHP / 2
	p.CurrentMana = p.MaxMana

	p.Initiative = species.Initiative
	p.DodgeRate = species.DodgeRate

	p.Gold = OrDepart
	p.InventorySize = InventaireBase
	p.MaxExperience = xpForNextLevel(p.Level)

	addInventory(p, Item{Name: PotionVie, Type: ItemPotion, Price: PrixPotionVie, HealValue: SoinPotionVie})
	addInventory(p, Item{Name: PotionVie, Type: ItemPotion, Price: PrixPotionVie, HealValue: SoinPotionVie})

	return p
}

func recalculateStats(p *Character) {
	bonusHP, bonusMana := equipmentBonus(p.Equipment)

	p.MaxHP = p.Species.BaseMaxHP +
		p.Species.HPPerStage*(p.Stage-1) +
		GainPVNiveau*(p.Level-1) +
		bonusHP

	p.MaxMana = p.Species.BaseMaxMana +
		p.Species.ManaPerStage*(p.Stage-1) +
		GainManaNiveau*(p.Level-1) +
		bonusMana

	p.Attack = p.Species.BaseAttack +
		p.Species.AttackPerStage*(p.Stage-1) +
		GainAtkNiveau*(p.Level-1)

	if p.CurrentHP > p.MaxHP {
		p.CurrentHP = p.MaxHP
	}
	if p.CurrentMana > p.MaxMana {
		p.CurrentMana = p.MaxMana
	}
}

func isDead(p *Character) bool {
	if p.CurrentHP > 0 {
		return false
	}
 
	if p.Species.Name == EspecePhenix && p.Stage == 4 {
		p.CurrentHP = p.MaxHP * RenaissanceHPPct / 100
		fmt.Println(p.Name + " renaît de ses cendres.")
		return true
	}
 
	p.CurrentHP = 1
	fmt.Println(p.Name + " est vaincu et ne se relève qu'à grand-peine (1 PV).")
	return true
}

func displayInfo(p *Character) {
	header := "\x1b[38;2;100;220;120m" + fmt.Sprintf("%s — %s (%s)", p.Name, currentStageName(p), p.Species.Element) + "\x1b[0m"
 
	infoLines := []string{
		header, "",
		fmt.Sprintf("Niveau %d  (XP : %d/%d)", p.Level, p.Experience, p.MaxExperience),
		"",
		"PV   " + healthBar(p.CurrentHP, p.MaxHP, 20, nil),
		"Mana " + healthBar(p.CurrentMana, p.MaxMana, 20, nil),
		fmt.Sprintf("Attaque : %d", p.Attack),
		fmt.Sprintf("Or : %d", p.Gold),
		"",
		"Équipement :",
		"  Tête  : " + slotName(p.Equipment.Head),
		"  Torse : " + slotName(p.Equipment.Chest),
		"  Pieds : " + slotName(p.Equipment.Feet),
	}
 
	if len(p.Spells) > 0 {
		infoLines = append(infoLines, "", "Sorts connus :")
		for _, s := range p.Spells {
			icon := spellIcon(s)
			infoLines = append(infoLines, fmt.Sprintf("  %s %s (%d mana)", icon, s.Name, s.ManaCost))
			if s.Description != "" {
				infoLines = append(infoLines, "     "+s.Description)
			}
		}
	}
 
	artStr := ""
	if imgPath := speciesImagePath(p.Species.Name, p.Stage); imgPath != "" {
		artStr, _ = renderImageString(imgPath, infoArtWidth)
	}
 
	merged := dividedColumns(artStr, strings.Join(infoLines, "\n"), infoArtWidth)
	lines := strings.Split(strings.TrimRight(merged, "\n"), "\n")
 
	fmt.Print(boxify(lines))
}

const infoArtWidth = 40