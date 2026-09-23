package main

import (
	"fmt"
	"sort"
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

func craftScreen(recipeList []Recipe, sel int, p *Character) string {
	labels := []string{"Retour"}
	for _, r := range recipeList {
		labels = append(labels, fmt.Sprintf("%s (%d or)", r.Result.Name, r.GoldCost))
	}
	width := maxWidth(labels...)

	lines := []string{
		fmt.Sprintf("\x1b[38;2;255;215;0mVous avez %d pièces d'or en poche.\x1b[0m", p.Gold),
		strings.Repeat("─", 30),
		"",
	}

	for i, recipe := range recipeList {
		ok, _ := canCraft(p, recipe)
		theme := themeOrange
		if !ok {
			theme = themeGris
		}
		line := menuLine(labels[i+1], width, i == sel, theme)
		if !ok && i != sel {
			line = "\x1b[90m" + line + "\x1b[0m"
		}
		lines = append(lines, line)

		bonuses := []string{}
		if recipe.Result.BonusMaxHP > 0 {
			bonuses = append(bonuses, fmt.Sprintf("+%d PV", recipe.Result.BonusMaxHP))
		}
		if recipe.Result.BonusMaxMana > 0 {
			bonuses = append(bonuses, fmt.Sprintf("+%d Mana", recipe.Result.BonusMaxMana))
		}
		if len(bonuses) > 0 {
			lines = append(lines, "      \x1b[36m"+strings.Join(bonuses, ", ")+"\x1b[0m")
		}

		names := make([]string, 0, len(recipe.Materials))
		for resource := range recipe.Materials {
			names = append(names, resource)
		}
		sort.Strings(names)

		for _, resource := range names {
			needed := recipe.Materials[resource]
			mat := fmt.Sprintf("      - %d %s", needed, resource)
			if countItem(p, resource) < needed {
				mat = "\x1b[91m" + mat + "\x1b[0m"
			}
			lines = append(lines, mat)
		}
		lines = append(lines, "")
	}

	lines = append(lines, menuLine("Retour", width, sel == len(recipeList), themeGris))
	return strings.Join(lines, "\n")
}

func craftMenu(p *Character) {
	list := recipes()
	banner, _ := renderImageString(bannerForge, bannerWidth)
	selected := 0

	for {
		choice := selectMenuFrom(len(list)+1, selected, func(sel int) string {
			return banner + "\n" + craftScreen(list, sel, p)
		})
		if choice == len(list) {
			return
		}
		selected = choice

		recipe := list[choice]
		if ok, reason := canCraft(p, recipe); !ok {
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
