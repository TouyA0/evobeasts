package main

import "fmt"

func recipes() []Recipe {
	return []Recipe{
		{
			Result:   ArtefactTete,
			GoldCost: 5,
			Materials: map[string]int{
				ResPlume:  1,
				ResEcaille: 1,
			},
		},
		{
			Result:   ArtefactTorse,
			GoldCost: 5,
			Materials: map[string]int{
				ResAbyssal: 2,
				ResEcaille: 1,
			},
		},
		{
			Result:   ArtefactPieds,
			GoldCost: 5,
			Materials: map[string]int{
				ResSpectre: 1,
				ResEcaille: 1,
			},
		},
	}
}

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

func canCraft(p *Character, r Recipe) bool {
	if p.Gold < r.GoldCost {
		return false
	}

	if len(p.Inventory) >= p.InventorySize {
		return false
	}

	for resource, needed := range r.Materials {
		count := 0

		for _, item := range p.Inventory {
			if item.Type == ItemResource && item.Name == resource {
				count++
			}
		}

		if count < needed {
			return false
		}
	}

	return true
}

func craftItem(p *Character, r Recipe) {
	p.Gold -= r.GoldCost

	for resource, needed := range r.Materials {
		removed := 0

		for i := 0; i < len(p.Inventory) && removed < needed; {
			if p.Inventory[i].Type == ItemResource && p.Inventory[i].Name == resource {
				p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
				removed++
			} else {
				i++
			}
		}
	}

	p.Inventory = append(p.Inventory, r.Result)
}

func craftDisplay(recipes []Recipe) {
	for i, recipe := range recipes {
		fmt.Println(i+1, "-", recipe.Result.Name)
		fmt.Println("Coût :", recipe.GoldCost, "or")
		fmt.Println("Ressources :")

		for resource, quantity := range recipe.Materials {
			fmt.Println("-", quantity, resource)
		}

		fmt.Println()
	}
}

