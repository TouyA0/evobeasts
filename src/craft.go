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
				BonusMaxHP:   6,
				BonusMaxMana: 16,
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
				BonusMaxHP: 22,
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
				BonusMaxHP:   8,
				BonusMaxMana: 8,
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

func craftDisplay(recipeList []Recipe, p *Character) {
	for i, recipe := range recipeList {
		ok, _ := canCraft(p, recipe)
		titleLine := fmt.Sprintf("%d. %s (%d or)", i+1, recipe.Result.Name, recipe.GoldCost)
		if !ok {
			titleLine = "\x1b[90m" + titleLine + "\x1b[0m"
		}
		fmt.Println(titleLine)

		bonuses := []string{}
		if recipe.Result.BonusMaxHP > 0 {
			bonuses = append(bonuses, fmt.Sprintf("+%d PV", recipe.Result.BonusMaxHP))
		}
		if recipe.Result.BonusMaxMana > 0 {
			bonuses = append(bonuses, fmt.Sprintf("+%d Mana", recipe.Result.BonusMaxMana))
		}
		if len(bonuses) > 0 {
			fmt.Printf("   \x1b[36m%s\x1b[0m\n", strings.Join(bonuses, ", "))
		}

		names := make([]string, 0, len(recipe.Materials))
		for resource := range recipe.Materials {
			names = append(names, resource)
		}
		sort.Strings(names)

		for _, resource := range names {
			needed := recipe.Materials[resource]
			line := fmt.Sprintf("   - %d %s", needed, resource)
			if countItem(p, resource) < needed {
				line = "\x1b[91m" + line + "\x1b[0m"
			}
			fmt.Println(line)
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
		craftDisplay(list, p)
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
