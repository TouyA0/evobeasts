package main

import (
	"fmt"
	"time"
)

// Tâche 14 + 18 : prix
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

// Tâche 14 : achat (on ne débite que si l'achat a réussi)
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

// Tâche 9 : poison
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

// Tâche 18 : augmentation d'inventaire (3 fois max)
func upgradeInventorySlot(p *Character) bool {
	if p.UpgradeCount >= 3 {
		fmt.Println("Augmentation max atteinte (3/3).")
		return false
	}
	p.MaxInventory += 10
	p.UpgradeCount++
	return true
}

// Mission 2 : expérience
func xpForLevel(level int) int {
	return 100 + (level-1)*50
}

func gainExperience(p *Character, amount int) {
	fmt.Printf("%s gagne %d XP\n", p.Name, amount)
	p.Experience += amount
	for p.Experience >= p.ExperienceMax {
		p.Experience -= p.ExperienceMax // l'excès est gardé
		levelUp(p)
	}
	fmt.Printf("XP : %d / %d\n", p.Experience, p.ExperienceMax)
}

func levelUp(p *Character) {
	p.Level++
	p.ExperienceMax = xpForLevel(p.Level)
	recalculateStats(p)
	p.HP = p.MaxHP
	fmt.Printf("NIVEAU %d ! PV max : %d | Mana max : %d | Attaque : %d\n",
		p.Level, p.MaxHP, p.MaxMana, p.Attack)
}
