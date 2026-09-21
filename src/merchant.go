package main

import (
	"fmt"
	"strconv"
	"time"
)

var prices = map[string]int{
	"Potion de vie":                3,
	"Potion de poison":             6,
	"Livre de Sort : Boule de Feu": 25,
	"Fourrure de Loup":             4,
	"Peau de Troll":                7,
	"Cuir de Sanglier":             3,
	"Plume de Corbeau":             1,
	"Augmentation d'inventaire":    30,
}

// on ne débite que si l'achat a réussi
func buy(p *Character, item Item) bool {
	price := prices[item.Name]
	if p.Money < price {
		fmt.Println("Pas assez d'or.")
		return false
	}
	if item.Name == "Augmentation d'inventaire" {
		if !upgradeInventorySlot(p) {
			return false
		}
	} else if !addInventory(p, item) {
		fmt.Println("Inventaire plein.")
		return false
	}
	p.Money -= price
	fmt.Println("Tu as obtenu :", item.Name)
	return true
}

func upgradeInventorySlot(p *Character) bool {
	if p.UpgradeCount >= 3 {
		fmt.Println("Augmentation max atteinte (3/3).")
		return false
	}
	p.MaxInventory += 10
	p.UpgradeCount++
	return true
}

func poisonPot(p *Character) {
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		p.HP -= 10
		if p.HP < 0 {
			p.HP = 0
		}
		fmt.Printf("PV : %d/%d\n", p.HP, p.MaxHP)
	}
	isDead(p)
}

func merchant(p *Character) {
	items := []string{
		"Potion de vie", "Potion de poison", "Livre de Sort : Boule de Feu",
		"Fourrure de Loup", "Peau de Troll", "Cuir de Sanglier",
		"Plume de Corbeau", "Augmentation d'inventaire",
	}
	for {
		fmt.Printf("\n=== Marchand === (Or : %d)\n", p.Money)
		for i, name := range items {
			fmt.Printf("%d. %s (%d or)\n", i+1, name, prices[name])
		}
		fmt.Println("0. Retour")

		choice, err := strconv.Atoi(readInput())
		if err != nil || choice < 0 || choice > len(items) {
			fmt.Println("Choix invalide.")
			continue
		}
		if choice == 0 {
			return
		}
		buy(p, Item{Name: items[choice-1]})
	}
}
