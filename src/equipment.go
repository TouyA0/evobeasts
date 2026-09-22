package main

import "fmt"

func equipItem(p *Character, item Item) bool {
	if item.Type != ItemArtifact {
		return false
	}

	if !removeInventory(p, item.Name) {
		return false
	}

	switch item.Slot {
	case SlotHead:
		if p.Equipment.Head != nil {
			p.Inventory = append(p.Inventory, *p.Equipment.Head)
		}
		p.Equipment.Head = &item

	case SlotChest:
		if p.Equipment.Chest != nil {
			p.Inventory = append(p.Inventory, *p.Equipment.Chest)
		}
		p.Equipment.Chest = &item

	case SlotFeet:
		if p.Equipment.Feet != nil {
			p.Inventory = append(p.Inventory, *p.Equipment.Feet)
		}
		p.Equipment.Feet = &item

	default:
		p.Inventory = append(p.Inventory, item)
		return false
	}

	recalculateStats(p)
	return true
}

func unequipItem(p *Character, slot EquipSlot) bool {
	if len(p.Inventory) >= p.InventorySize {
		return false
	}

	switch slot {
	case SlotHead:
		if p.Equipment.Head == nil {
			return false
		}
		p.Inventory = append(p.Inventory, *p.Equipment.Head)
		p.Equipment.Head = nil

	case SlotChest:
		if p.Equipment.Chest == nil {
			return false
		}
		p.Inventory = append(p.Inventory, *p.Equipment.Chest)
		p.Equipment.Chest = nil

	case SlotFeet:
		if p.Equipment.Feet == nil {
			return false
		}
		p.Inventory = append(p.Inventory, *p.Equipment.Feet)
		p.Equipment.Feet = nil

	default:
		return false
	}

	recalculateStats(p)
	return true
}

func equipmentBonus(e Equipment) (hp int, mana int) {
	if e.Head != nil {
		hp += e.Head.BonusMaxHP
		mana += e.Head.BonusMaxMana
	}
	if e.Chest != nil {
		hp += e.Chest.BonusMaxHP
		mana += e.Chest.BonusMaxMana
	}
	if e.Feet != nil {
		hp += e.Feet.BonusMaxHP
		mana += e.Feet.BonusMaxMana
	}
	return hp, mana
}

func equipmentDisplay(e Equipment) {
	if e.Head != nil {
		fmt.Println("Tête :", e.Head.Name)
	} else {
		fmt.Println("Tête : vide")
	}

	if e.Chest != nil {
		fmt.Println("Torse :", e.Chest.Name)
	} else {
		fmt.Println("Torse : vide")
	}

	if e.Feet != nil {
		fmt.Println("Pieds :", e.Feet.Name)
	} else {
		fmt.Println("Pieds : vide")
	}
}