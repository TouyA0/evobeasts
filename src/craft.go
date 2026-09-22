package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func recipes() []Recipe {
	return []Recipe{
		{
			Result: Item{
				Name:         ArtefactTete,
				Type:         ItemArtifact,
				Slot:         SlotHead,
				BonusMaxHP:   10,
				BonusMaxMana: 10,
			},
			GoldCost: CoutCraft,
			Materials: map[string]int{
				ResPlume:   1,
				ResEcaille: 1,
			},
		},
		{
			Result: Item{
				Name:       ArtefactTorse,
				Type:       ItemArtifact,
				Slot:       SlotChest,
				BonusMaxHP: 25,
			},
			GoldCost: CoutCraft,
			Materials: map[string]int{
				ResAbyssal: 2,
				ResEcaille: 1,
			},
		},
		{
			Result: Item{
				Name:         ArtefactPieds,
				Type:         ItemArtifact,
				Slot:         SlotFeet,
				BonusMaxHP:   15,
				BonusMaxMana: 5,
			},
			GoldCost: CoutCraft,
			Materials: map[string]int{
				ResSpectre: 1,
				ResEcaille: 1,
			},
		},
	}
}

func canCraft(p *Character, r Recipe) (bool, string) {
	if p.Gold < r.GoldCost {
		return false, "Pas assez d'or."
	}

	materialsCount := 0
	for _, needed := range r.Materials {
		materialsCount += needed
	}
	if len(p.Inventory)-materialsCount >= p.InventorySize {
		return false, "Pas assez de place dans l'inventaire."
	}

	for resource, needed := range r.Materials {
		if countItem(p, resource) < needed {
			return false, "Ressources manquantes : " + resource
		}
	}

	return true, ""
}

func craftItem(p *Character, r Recipe) bool {
	ok, _ := canCraft(p, r)
	if !ok {
		return false
	}

	p.Gold -= r.GoldCost

	for resource, needed := range r.Materials {
		for i := 0; i < needed; i++ {
			removeInventory(p, resource)
		}
	}

	p.Inventory = append(p.Inventory, r.Result)
	return true
}

func craftDisplay(recipeList []Recipe) {
	for i, recipe := range recipeList {
		fmt.Printf("%d. %s (%d or)\n", i+1, recipe.Result.Name, recipe.GoldCost)

		names := make([]string, 0, len(recipe.Materials))
		for resource := range recipe.Materials {
			names = append(names, resource)
		}
		sort.Strings(names)

		for _, resource := range names {
			fmt.Printf("   - %d %s\n", recipe.Materials[resource], resource)
		}
		fmt.Println()
	}
	fmt.Println("0. Retour")
}

func craftMenu(p *Character) {
	list := recipes()

	for {
		clearScreen()
		displayBanner(bannerForge)
		fmt.Printf("\n\x1b[38;2;255;215;0mVous avez %d pièces d'or en poche.\x1b[0m\n", p.Gold)
		fmt.Println(strings.Repeat("─", 30))
		craftDisplay(list)
		fmt.Print("Votre choix : ")

		choice := readInput()
		if choice == "0" {
			return
		}

		idx, err := strconv.Atoi(choice)
		if err != nil || idx < 1 || idx > len(list) {
			fmt.Println("Choix invalide.")
			pauseForInput()
			continue
		}

		recipe := list[idx-1]
		ok, reason := canCraft(p, recipe)
		if !ok {
			fmt.Println(reason)
			pauseForInput()
			continue
		}

		craftItem(p, recipe)
		fmt.Println("Vous avez fabriqué : " + recipe.Result.Name)
		pauseForInput()
	}
}

func showArtists() {
	fmt.Println("Partie 2 (l'économie) fait référence à des chansons d'ABBA.")
	fmt.Println("Partie 3 (le combat) fait référence à des films de Steven Spielberg.")
}
