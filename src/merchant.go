package main

import (
	"fmt"
	"strconv"
	"strings"
)

func merchantCatalog() []Item {
	return []Item{
		{Name: PotionVie, Type: ItemPotion, Price: PrixPotionVie, HealValue: SoinPotionVie,
			Description: fmt.Sprintf("Soigne %d PV instantanément", SoinPotionVie)},
		{Name: PotionPoison, Type: ItemPotion, Price: PrixPotionPoison,
			PoisonDamage: DegatsPotionPoison, PoisonTicks: TicksPotionPoison,
			Description: fmt.Sprintf("Inflige %d dégâts de poison sur %d tours", DegatsPotionPoison, TicksPotionPoison)},
		{Name: PotionMana, Type: ItemPotion, Price: PrixPotionMana, ManaValue: ManaPotionMana,
			Description: fmt.Sprintf("Restaure %d points de mana", ManaPotionMana)},
		{Name: GrimoireSoufle, Type: ItemSpellbook, Price: PrixGrimoire, SpellName: SortSoufflePrimordial,
			Description: "Enseigne le sort Souffle primordial"},
		{Name: ResEcaille, Type: ItemResource, Price: PrixEcaille},
		{Name: ResAbyssal, Type: ItemResource, Price: PrixAbyssal},
		{Name: ResSpectre, Type: ItemResource, Price: PrixSpectre},
		{Name: ResPlume, Type: ItemResource, Price: PrixPlume},
		{Name: ExpansionSac, Type: ItemUpgrade, Price: PrixExpansion,
			Description: fmt.Sprintf("+%d emplacements d'inventaire", InventaireBonus)},
	}
}

// on ne débite que si l'achat a réussi
func buyItem(p *Character, item Item) bool {
	if p.Gold < item.Price {
		fmt.Println("Pas assez d'or.")
		return false
	}
	if item.Type == ItemUpgrade {
		if !upgradeInventorySlot(p) {
			return false
		}
	} else if !addInventory(p, item) {
		fmt.Println("Inventaire plein.")
		return false
	}
	p.Gold -= item.Price
	fmt.Println("Tu as obtenu :", item.Name)
	return true
}

func merchantMenu(p *Character) {
	items := merchantCatalog()

	for {
		clearScreen()
		displayBanner(bannerMerchant)
		fmt.Printf("\n\x1b[38;2;255;215;0mVous avez %d pièces d'or en poche.\x1b[0m\n", p.Gold)
		fmt.Printf("Inventaire : %d/%d\n", len(p.Inventory), p.InventorySize)
		fmt.Println(strings.Repeat("─", 30))
		merchantDisplay(items, p.Gold, p)

		fmt.Print("Votre choix : ")
		choice, err := strconv.Atoi(readInput())
		if err != nil || choice < 0 || choice > len(items) {
			fmt.Println("Choix invalide.")
			pauseForInput()
			continue
		}
		if choice == 0 {
			return
		}
		buyItem(p, items[choice-1])
		pauseForInput()
	}
}

func upgradeInventorySlot(p *Character) bool {
	if p.UpgradesUsed >= MaxExpansions {
		fmt.Println("Augmentation max atteinte.")
		return false
	}
	p.InventorySize += InventaireBonus
	p.UpgradesUsed++
	return true
}

func merchantDisplay(items []Item, gold int, p *Character) {
	categories := []struct {
		label string
		typ   ItemType
	}{
		{"Consommables", ItemPotion},
		{"Matériaux", ItemResource},
		{"Objets rares", ItemSpellbook},
		{"Améliorations", ItemUpgrade},
	}

	for _, cat := range categories {
		printed := false
		for i, it := range items {
			if it.Type != cat.typ {
				continue
			}
			if !printed {
				fmt.Printf("\n-- %s --\n", cat.label)
				printed = true
			}
			line := fmt.Sprintf("%d. %-30s (%d or)", i+1, it.Name, it.Price)
			if it.Type == ItemUpgrade {
				line += fmt.Sprintf("  \x1b[90m[%d/%d utilisées]\x1b[0m", p.UpgradesUsed, MaxExpansions)
			} else if owned := countItem(p, it.Name); owned > 0 {
				line += fmt.Sprintf("  \x1b[90m[possédé: x%d]\x1b[0m", owned)
			}
			if it.Price > gold {
				line = "\x1b[90m" + line + "\x1b[0m"
			}
			fmt.Println(line)
			if it.Description != "" {
				fmt.Printf("   \x1b[90m%s\x1b[0m\n", it.Description)
			}
		}
	}
	fmt.Println("\n0. Retour")
}
