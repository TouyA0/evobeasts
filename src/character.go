package main

import (
	"fmt"
	"strings"
)

var phoenix = []Character{
	{Species: Species{Name: "Œuf ardent"}, Stage: 1, MaxHP: 50, MaxMana: 40, Attack: 5},
	{Species: Species{Name: "Poussin"}, Stage: 2, MaxHP: 60, MaxMana: 45, Attack: 8},
	{Species: Species{Name: "Poulet"}, Stage: 3, MaxHP: 80, MaxMana: 50, Attack: 12},
	{Species: Species{Name: "Phénix"}, Stage: 4, MaxHP: 100, MaxMana: 70, Attack: 20},
}

var kraken = []Character{
	{Species: Species{Name: "Œuf immergé"}, Stage: 1, MaxHP: 75, MaxMana: 30, Attack: 5},
	{Species: Species{Name: "Alevin"}, Stage: 2, MaxHP: 100, MaxMana: 35, Attack: 7},
	{Species: Species{Name: "Poisson"}, Stage: 3, MaxHP: 140, MaxMana: 40, Attack: 10},
	{Species: Species{Name: "Kraken"}, Stage: 4, MaxHP: 200, MaxMana: 60, Attack: 15},
}

var chameleon = []Character{
	{Species: Species{Name: "Œuf voilé"}, Stage: 1, MaxHP: 45, MaxMana: 35, Attack: 5},
	{Species: Species{Name: "Caméléonneau"}, Stage: 2, MaxHP: 55, MaxMana: 40, Attack: 9},
	{Species: Species{Name: "Caméléon"}, Stage: 3, MaxHP: 70, MaxMana: 50, Attack: 14},
	{Species: Species{Name: "Caméléon Spectral"}, Stage: 4, MaxHP: 90, MaxMana: 65, Attack: 25},
}

func readInput() string {
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

func characterCreation() *Character {
	var name string
	var choice string

	fmt.Println("Entrez le nom de votre personnage :")
	name = readInput()

	fmt.Println("Choisissez votre œuf :")
	fmt.Println("1. Œuf ardent - DPS / Feu")
	fmt.Println("2. Œuf immergé - Tank / Eau")
	fmt.Println("3. Œuf voilé - Assassin / Spectral")
	choice = readInput()

	for choice != "1" && choice != "2" && choice != "3" {
		fmt.Println("Veuillez entrer un nombre valide (1, 2 ou 3) :")
		choice = readInput()
	}

	var species Species

	switch choice {
	case "1":
		species = Species{Name: "Œuf ardent"}
	case "2":
		species = Species{Name: "Œuf immergé"}
	case "3":
		species = Species{Name: "Œuf voilé"}
	}

	return initCharacter(name, species)
}

func initCharacter(name string, species Species) *Character {
	character := &Character{
		Name:    name,
		Species: species,
		Stage:   1,
	}

	switch {
	case strings.EqualFold(species.Name, "Œuf ardent"):
		character.MaxHP = phoenix[0].MaxHP
		character.MaxMana = phoenix[0].MaxMana
		character.Attack = phoenix[0].Attack

	case strings.EqualFold(species.Name, "Œuf immergé"):
		character.MaxHP = kraken[0].MaxHP
		character.MaxMana = kraken[0].MaxMana
		character.Attack = kraken[0].Attack

	case strings.EqualFold(species.Name, "Œuf voilé"):
		character.MaxHP = chameleon[0].MaxHP
		character.MaxMana = chameleon[0].MaxMana
		character.Attack = chameleon[0].Attack
	}

	character.CurrentHP = character.MaxHP

	return character
}

func displayInfo(p *Character) {
	fmt.Printf("Nom : %s\n", p.Name)
	fmt.Printf("Œuf : %s\n", p.Species.Name)
	fmt.Printf("Palier : %d\n", p.Stage)
	fmt.Printf("HP : %d/%d\n", p.CurrentHP, p.MaxHP)
	fmt.Printf("Mana : %d\n", p.MaxMana)
	fmt.Printf("Attack : %d\n", p.Attack)

	displayCreature(p.Species.Name, p.Stage)
}
