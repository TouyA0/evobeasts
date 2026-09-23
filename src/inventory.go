package main

import (
	"fmt"
	"strings"
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

func itemDescription(item Item) string {
	switch {
	case item.HealValue > 0:
		return fmt.Sprintf("Rend %d PV", item.HealValue)
	case item.ManaValue > 0:
		return fmt.Sprintf("Rend %d Mana", item.ManaValue)
	case item.PoisonDamage > 0:
		return fmt.Sprintf("Poison : %d dégâts sur %d tours", item.PoisonDamage, item.PoisonTicks)
	case item.Type == ItemArtifact:
		var parts []string
		if item.Slot != SlotNone {
			parts = append(parts, "Emplacement : "+string(item.Slot))
		}
		if item.BonusMaxHP > 0 {
			parts = append(parts, fmt.Sprintf("+%d PV max", item.BonusMaxHP))
		}
		if item.BonusMaxMana > 0 {
			parts = append(parts, fmt.Sprintf("+%d Mana max", item.BonusMaxMana))
		}
		return strings.Join(parts, ", ")
	case item.Type == ItemSpellbook:
		return "Enseigne : " + item.SpellName
	case item.Type == ItemResource:
		return "Matériau de craft"
	case item.Type == ItemUpgrade:
		return "Agrandit l'inventaire"
	default:
		return ""
	}
}

var inventoryCategories = []struct {
	Type  ItemType
	Label string
	Theme menuTheme
}{
	{ItemPotion, "Potions", themeRouge},
	{ItemSpellbook, "Grimoires", themeViolet},
	{ItemArtifact, "Artefacts", themeOr},
	{ItemResource, "Ressources", themeTurquoise},
	{ItemUpgrade, "Améliorations", themeBleu},
}

type inventoryEntry struct {
	category    string
	theme       menuTheme
	label, desc string
	index       int // position dans l'inventaire
}

func inventoryEntries(items []Item) []inventoryEntry {
	var entries []inventoryEntry
	for _, cat := range inventoryCategories {
		seen := make(map[string]bool)
		for i, item := range items {
			if item.Type != cat.Type || seen[item.Name] {
				continue
			}
			seen[item.Name] = true

			count := 0
			for _, other := range items {
				if other.Name == item.Name {
					count++
				}
			}
			label := item.Name
			if count > 1 {
				label = fmt.Sprintf("%s x%d", item.Name, count)
			}

			entries = append(entries, inventoryEntry{
				category: cat.Label,
				theme:    cat.Theme,
				label:    label,
				desc:     itemDescription(item),
				index:    i,
			})
		}
	}
	return entries
}

func inventoryScreen(p *Character, entries []inventoryEntry, sel int) string {
	labels := []string{"Retour"}
	for _, e := range entries {
		labels = append(labels, e.label)
	}
	width := maxWidth(labels...)

	lines := []string{fmt.Sprintf("(Emplacements : %d/%d)", len(p.Inventory), p.InventorySize)}
	category := ""
	for i, e := range entries {
		if e.category != category {
			lines = append(lines, "", "-- "+e.category+" --")
			category = e.category
		}
		line := menuLine(e.label, width, i == sel, e.theme)
		if e.desc != "" {
			line += "  \x1b[90m" + e.desc + "\x1b[0m"
		}
		lines = append(lines, line)
	}
	lines = append(lines, "", menuLine("Retour", width, sel == len(entries), themeGris))
	return strings.Join(lines, "\n")
}

func accessInventory(character *Character) {
	banner, _ := renderImageString(bannerInventaire, bannerWidth)
	selected := 0

	for {
		if len(character.Inventory) == 0 {
			fmt.Println("Inventaire vide.")
			pauseForInput()
			return
		}

		entries := inventoryEntries(character.Inventory)
		choice := selectMenuFrom(len(entries)+1, selected, func(sel int) string {
			return banner + "\n" + inventoryScreen(character, entries, sel)
		})
		if choice == len(entries) {
			return
		}
		selected = choice

		useItem(character, character.Inventory[entries[choice].index])
		pauseForInput()
	}
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
	character.CurrentHP += SoinPotionVie
	if character.CurrentHP > character.MaxHP {
		character.CurrentHP = character.MaxHP
	}
	fmt.Printf("PV : %d/%d\n", character.CurrentHP, character.MaxHP)
}

func poisonPot(p *Character, damage, ticks int) bool {
	if !confirm("⚠️  Cette potion est du poison : elle va vous infliger des dégâts.",
		"Voulez-vous vraiment la boire ?") {
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
