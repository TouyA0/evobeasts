package main

import (
	"fmt"
	"strconv"
	"strings"
)

func merchantCatalog() []Item {
	return []Item{
		{Name: PotionVie, Type: ItemPotion, Price: 3, HealValue: 50},
		{Name: PotionPoison, Type: ItemPotion, Price: 6, PoisonDamage: 10, PoisonTicks: 3},
		{Name: PotionMana, Type: ItemPotion, Price: 5, ManaValue: 30},
		{Name: GrimoireSoufle, Type: ItemSpellbook, Price: 25, SpellName: SortSoufflePrimordial},
		{Name: ResEcaille, Type: ItemResource, Price: 4},
		{Name: ResAbyssal, Type: ItemResource, Price: 7},
		{Name: ResSpectre, Type: ItemResource, Price: 3},
		{Name: ResPlume, Type: ItemResource, Price: 1},
		{Name: ExpansionSac, Type: ItemUpgrade, Price: 30},
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
		merchantDisplay(items)
		merchantDisplay(items)

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

func merchantDisplay(items []Item) {
	for i, it := range items {
		fmt.Printf("%d. %-30s (%d or)\n", i+1, it.Name, it.Price)
	}
	fmt.Println("0. Retour")
}
