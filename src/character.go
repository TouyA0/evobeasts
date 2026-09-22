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
	fmt.Println("Entrez le nom de votre personnage :")
	name := readInput()
	for !isLettersOnly(name) {
		fmt.Println("Le nom ne doit contenir que des lettres. Réessayez :")
		name = readInput()
	}
	name = formatName(name)

	fmt.Println("Choisissez votre œuf :")
	fmt.Println("1. Œuf ardent - DPS / Feu")
	fmt.Println("2. Œuf immergé - Tank / Eau")
	fmt.Println("3. Œuf voilé - Assassin / Spectre")
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

	return initCharacter(name, species)
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

	addInventory(p, Item{Name: PotionVie, Type: ItemPotion, Price: 3, HealValue: 50})
	addInventory(p, Item{Name: PotionVie, Type: ItemPotion, Price: 3, HealValue: 50})

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
		p.CurrentHP = p.MaxHP / 2
		fmt.Println(p.Name + " renaît de ses cendres avec la moitié de ses PV.")
		return true
	}
 
	p.CurrentHP = 1
	fmt.Println(p.Name + " est vaincu et ne se relève qu'à grand-peine (1 PV).")
	return true
}

func displayInfo(p *Character) {
	fmt.Printf("Nom : %s\n", p.Name)
	fmt.Printf("Espèce : %s (%s)\n", p.Species.Name, p.Species.Element)
	fmt.Printf("Palier : %d — %s\n", p.Stage, currentStageName(p))
	fmt.Printf("Niveau : %d (XP : %d/%d)\n", p.Level, p.Experience, p.MaxExperience)
	fmt.Printf("PV   : %d/%d\n", p.CurrentHP, p.MaxHP)
	fmt.Printf("Mana : %d/%d\n", p.CurrentMana, p.MaxMana)
	fmt.Printf("Attaque : %d\n", p.Attack)
	fmt.Printf("Or : %d\n", p.Gold)

	fmt.Println("Équipement :")
	equipmentDisplay(p.Equipment)

	if len(p.Spells) > 0 {
		fmt.Println("Sorts connus :")
		for _, s := range p.Spells {
			fmt.Printf("  - %s\n", s.Name)
		}
	}

	displayPlayerArt(p)
}

const infoArtWidth = 70
func displayPlayerArt(p *Character) {
	if imgPath := speciesImagePath(p.Species.Name, p.Stage); imgPath != "" {
		if rendered, ok := renderImageString(imgPath, infoArtWidth); ok {
			fmt.Println(rendered)
		}
	}
}