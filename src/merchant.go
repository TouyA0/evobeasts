package main

import (
	"fmt"
	"strconv"
	"time"
)

var merchantItems = []Item{
	{Name: "Potion de vie", Type: ItemPotion, Price: 3, HealValue: 50},
	{Name: "Potion de poison", Type: ItemPotion, Price: 6, PoisonDamage: 10, PoisonTicks: 3},
	{Name: "Livre de Sort : Boule de Feu", Type: ItemSpellbook, Price: 25, SpellName: "Boule de Feu"},
	{Name: "Fourrure de Loup", Type: ItemResource, Price: 4},
	{Name: "Peau de Troll", Type: ItemResource, Price: 7},
	{Name: "Cuir de Sanglier", Type: ItemResource, Price: 3},
	{Name: "Plume de Corbeau", Type: ItemResource, Price: 1},
	{Name: "Augmentation d'inventaire", Type: ItemUpgrade, Price: 30},
}

// on ne débite que si l'achat a réussi
func buy(p *Character, item Item) bool {
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

func merchant(p *Character) {
	for {
		fmt.Printf("\n=== Marchand === (Or : %d)\n", p.Gold)
		for i, it := range merchantItems {
			fmt.Printf("%d. %s (%d or)\n", i+1, it.Name, it.Price)
		}
		fmt.Println("0. Retour")

		choice, err := strconv.Atoi(readInput())
		if err != nil || choice < 0 || choice > len(merchantItems) {
			fmt.Println("Choix invalide.")
			continue
		}
		if choice == 0 {
			return
		}
		buy(p, merchantItems[choice-1])
	}
}
func poisonPot(p *Character, damage, ticks int) {
	for i := 0; i < ticks; i++ {
		time.Sleep(1 * time.Second)
		p.CurrentHP -= damage
		if p.CurrentHP < 0 {
			p.CurrentHP = 0
		}
		fmt.Printf("PV : %d/%d\n", p.CurrentHP, p.MaxHP)
	}
	isDead(p)
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
