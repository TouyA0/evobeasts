package main

import (
	"fmt"
	"strconv"
	"time"
)

func hasSpace(character *Character) bool {
	return len(character.Inventory) < character.InventorySize
}

func addInventory(character *Character, item Item) bool {
	if !hasSpace(character) {
		return false
	}
	character.Inventory = append(character.Inventory, item)
	return true
}

func removeInventory(character *Character, name string) bool {
	for i, item := range character.Inventory {
		if item.Name == name {
			character.Inventory = append(character.Inventory[:i], character.Inventory[i+1:]...)
			return true
		}
	}
	return false
}

func countItem(character *Character, name string) int {
	count := 0
	for _, item := range character.Inventory {
		if item.Name == name {
			count++
		}
	}
	return count
}

func accessInventory(character *Character) {
	if len(character.Inventory) == 0 {
		fmt.Println("Inventaire vide.")
		return
	}

	fmt.Println("Inventaire :")
	for i, item := range character.Inventory {
		fmt.Printf("%d. %s\n", i+1, item.Name)
	}
	fmt.Println("0. Retour")
	fmt.Print("Votre choix : ")

	choice := readInput()
	if choice == "0" {
		return
	}
	idx, err := strconv.Atoi(choice)
	if err != nil || idx < 1 || idx > len(character.Inventory) {
		fmt.Println("Choix invalide.")
		return
	}

	useItem(character, character.Inventory[idx-1])
}

func useItem(p *Character, item Item) {
	switch item.Type {
	case ItemPotion:
		usePotion(p, item)
	case ItemSpellbook:
		if learnSpell(p, item.SpellName) {
			fmt.Println("Vous apprenez : " + item.SpellName)
			removeInventory(p, item.Name)
		} else {
			fmt.Println("Vous connaissez déjà ce sort.")
		}
	case ItemArtifact:
		if equipItem(p, item) {
			fmt.Println(item.Name + " est équipé.")
		}
	default:
		fmt.Println("Cet objet ne peut pas être utilisé ainsi.")
	}
}

func usePotion(p *Character, item Item) {
	switch {
	case item.HealValue > 0:
		p.CurrentHP += item.HealValue
		if p.CurrentHP > p.MaxHP {
			p.CurrentHP = p.MaxHP
		}
		fmt.Printf("PV : %d/%d\n", p.CurrentHP, p.MaxHP)
		removeInventory(p, item.Name)

	case item.ManaValue > 0:
		p.CurrentMana += item.ManaValue
		if p.CurrentMana > p.MaxMana {
			p.CurrentMana = p.MaxMana
		}
		fmt.Printf("Mana : %d/%d\n", p.CurrentMana, p.MaxMana)
		removeInventory(p, item.Name)

	case item.PoisonDamage > 0:
		if poisonPot(p, item.PoisonDamage, item.PoisonTicks) {
			removeInventory(p, item.Name)
		}
	}
}

func takePot(character *Character) {
	if !removeInventory(character, PotionVie) {
		fmt.Println("Vous n'avez pas d'Essence vitale.")
		return
	}
	character.CurrentHP += 50
	if character.CurrentHP > character.MaxHP {
		character.CurrentHP = character.MaxHP
	}
	fmt.Printf("PV : %d/%d\n", character.CurrentHP, character.MaxHP)
}

func poisonPot(p *Character, damage, ticks int) bool {
	fmt.Println("⚠️  Cette potion est du poison : elle va vous infliger des dégâts.")
	fmt.Print("Voulez-vous vraiment la boire ? (oui/non) : ")
	if readInput() != "oui" {
		fmt.Println("Vous rangez la potion.")
		return false
	}

	for i := 0; i < ticks; i++ {
		time.Sleep(1 * time.Second)
		p.CurrentHP -= damage
		if p.CurrentHP < 0 {
			p.CurrentHP = 0
		}
		fmt.Printf("PV : %d/%d\n", p.CurrentHP, p.MaxHP)
	}
	isDead(p)
	return true
}
