package main

import (
	"fmt"
	"strings"
)

func accessInventory(character *Character) {
	fmt.Println("Inventaire:")

	for _, item := range character.Inventory {
		fmt.Printf("- %s\n", item.Name)
	}
}

func hasSpace(character *Character) bool {
	return len(character.Inventory) < 10
}

func addInventory(character *Character, item Item) {
	if len(character.Inventory) >= 10 {
		fmt.Println("Inventaire plein !")
		return
	}

	character.Inventory = append(character.Inventory, item)
}

func formatName(s string) string {
	if s == "" {
		return ""
	}

	s = strings.ToLower(s)
	runes := []rune(s)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]

	return string(runes)
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

func takePot(character *Character) {
	character.CurrentHP += 50

	if character.CurrentHP > character.MaxHP {
		character.CurrentHP = character.MaxHP
	}
}
