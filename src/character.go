package main

import (
	"fmt"
	"strings"
)

type Character struct {
	Name string
	Race string
}

func characterCreation() *Character {
	var name string
	var race string

	fmt.Println("Entrez le nom de votre personnage :")
	fmt.Scanln(&name)

	fmt.Println("Entrez la race de votre personnage (Phénix/Kraken/CaméléonSpectral) :")
	fmt.Scanln(&race)
	for !strings.EqualFold(race, "Phénix") && !strings.EqualFold(race, "Phenix") && !strings.EqualFold(race, "Kraken") && !strings.EqualFold(race, "CaméléonSpectral") && !strings.EqualFold(race, "CameleonSpectral") {
		fmt.Println("Veuillez entrer Phénix, Kraken ou CaméléonSpectral uniquement :")
		fmt.Scanln(&race)
	}

	return initCharacter(name, race)
}

func initCharacter(name string, race string) *Character {
	character := &Character{
		Name: name,
		Race: race,
	}
	return character
}

func main() {
	characterCreation()
}
