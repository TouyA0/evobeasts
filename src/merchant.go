package main

import (
	"fmt"
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


func upgradeInventorySlot(p *Character) bool {
	if p.UpgradesUsed >= MaxExpansions {
		fmt.Println("Augmentation max atteinte.")
		return false
	}
	p.InventorySize += InventaireBonus
	p.UpgradesUsed++
	return true
}

var merchantCategories = []struct {
	label string
	typ   ItemType
	theme menuTheme
}{
	{"Consommables", ItemPotion, themeRouge},
	{"Matériaux", ItemResource, themeTurquoise},
	{"Objets rares", ItemSpellbook, themeViolet},
	{"Améliorations", ItemUpgrade, themeBleu},
}

// merchantOrder donne les index du catalogue dans l'ordre d'affichage.
func merchantOrder(items []Item) []int {
	var order []int
	for _, cat := range merchantCategories {
		for i, it := range items {
			if it.Type == cat.typ {
				order = append(order, i)
			}
		}
	}
	return order
}

func merchantScreen(items []Item, sel int, p *Character) string {
	const width = 40
	lines := []string{
		fmt.Sprintf("\x1b[38;2;255;215;0mVous avez %d pièces d'or en poche.\x1b[0m", p.Gold),
		fmt.Sprintf("Inventaire : %d/%d", len(p.Inventory), p.InventorySize),
		strings.Repeat("─", 30),
	}

	pos := 0
	for _, cat := range merchantCategories {
		first := true
		for _, it := range items {
			if it.Type != cat.typ {
				continue
			}
			if first {
				lines = append(lines, "", "-- "+cat.label+" --")
				first = false
			}

			tooExpensive := it.Price > p.Gold
			theme := cat.theme
			if tooExpensive {
				theme = themeGris
			}
			line := menuLine(fmt.Sprintf("%-30s (%d or)", it.Name, it.Price), width, pos == sel, theme)
			if tooExpensive && pos != sel {
				line = "\x1b[90m" + line + "\x1b[0m"
			}

			if it.Type == ItemUpgrade {
				line += fmt.Sprintf("  \x1b[90m[%d/%d utilisées]\x1b[0m", p.UpgradesUsed, MaxExpansions)
			} else if owned := countItem(p, it.Name); owned > 0 {
				line += fmt.Sprintf("  \x1b[90m[possédé: x%d]\x1b[0m", owned)
			}
			lines = append(lines, line)

			if it.Description != "" {
				lines = append(lines, "      \x1b[90m"+it.Description+"\x1b[0m")
			}
			pos++
		}
	}

	lines = append(lines, "", menuLine("Retour", width, sel == pos, themeGris))
	return strings.Join(lines, "\n")
}

func merchantMenu(p *Character) {
	items := merchantCatalog()
	order := merchantOrder(items)
	banner, _ := renderImageString(bannerMerchant, bannerWidth)
	selected := 0

	for {
		choice := selectMenuFrom(len(order)+1, selected, func(sel int) string {
			return banner + "\n" + merchantScreen(items, sel, p)
		})
		if choice == len(order) {
			return
		}
		selected = choice

		buyItem(p, items[order[choice]])
		pauseForInput()
	}
}