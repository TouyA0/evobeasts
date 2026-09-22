package main

import (
	"fmt"
	"strconv"
	"strings"
)

func merchantCatalog() []Item {
	return []Item{
		{Name: PotionVie, Type: ItemPotion, Price: 3, HealValue: 50,
			Description: "Soigne 50 PV instantanément"},
		{Name: PotionPoison, Type: ItemPotion, Price: 6, PoisonDamage: 10, PoisonTicks: 3,
			Description: "Inflige 10 dégâts de poison sur 3 tours"},
		{Name: PotionMana, Type: ItemPotion, Price: 5, ManaValue: 30,
			Description: "Restaure 30 points de mana"},
		{Name: GrimoireSoufle, Type: ItemSpellbook, Price: 25, SpellName: SortSoufflePrimordial,
			Description: "Enseigne le sort Souffle primordial"},
		{Name: ResEcaille, Type: ItemResource, Price: 4},
		{Name: ResAbyssal, Type: ItemResource, Price: 7},
		{Name: ResSpectre, Type: ItemResource, Price: 3},
		{Name: ResPlume, Type: ItemResource, Price: 1},
		{Name: ExpansionSac, Type: ItemUpgrade, Price: 30,
			Description: "+10 emplacements d'inventaire"},
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
		fmt.Println(strings.Repeat("─", 30))
		merchantDisplay(items, p.Gold)

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
	p.InventorySize += 10
	p.UpgradesUsed++
	return true
}

func merchantDisplay(items []Item, gold int) {
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
			if it.Price > gold {
				line = "\x1b[90m" + line + "\x1b[0m"
			}
			fmt.Println(line)
		}
	}
	fmt.Println("\n0. Retour")
}
