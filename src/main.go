package main

import "fmt"

func main() {
	clearScreen()
	displayTitleScreen()
	playMusic(musicAmbiance)
	pauseForInput()
	clearScreen()
	storyPrologue()
	player := characterCreation()
	clearScreen()
	mainMenu(player)
}

func mainMenu(p *Character) {
	options := []menuOption{
		{"Informations du personnage", themeVert},
		{"Inventaire", themeBleu},
		{"Marchand", themeOr},
		{"Forgeron", themeOrange},
		{"Combat", themeRouge},
		{"Qui sont-ils ?", themeViolet},
		{"Mémoires de la vallée", themeTurquoise},
		{"Quitter", themeGris},
	}

	// Calculée une seule fois : c'est ce qui rend la navigation fluide.
	banner, _ := renderImageString(bannerMenu, bannerWidth)

	for {
		choice := selectMenu(len(options), func(sel int) string {
			t := elementTheme(p.Species.Element)
			lines := []string{
				"\x1b[38;2;" + t.bg + "m" + fmt.Sprintf("%s — %s (%s)", p.Name, currentStageName(p), p.Species.Element) + "\x1b[0m",
				fmt.Sprintf("Niveau %d   PV %s   Mana %s   Or %d",
					p.Level,
					healthBar(p.CurrentHP, p.MaxHP, 12, nil),
					healthBar(p.CurrentMana, p.MaxMana, 12, nil),
					p.Gold),
				"",
			}
			lines = append(lines, menuLines(options, sel)...)
			return banner + "\n" + boxify(lines)
		})
		clearScreen()

		switch choice {
		case 0:
			displayInfo(p)
			backButton(themeVert)
			continue
		case 1:
			accessInventory(p)
			continue
		case 2:
			merchantMenu(p)
			continue
		case 3:
			craftMenu(p)
			continue
		case 4:
			combatMenu(p)
			if isGameComplete(p) {
				stopMusic()
				return
			}
			continue
		case 5:
			showArtists()
		case 6:
			storyJournal()
			continue
		case 7:
			stopMusic()
			fmt.Println("À bientôt !")
			return
		}

		pauseForInput()
	}
}
