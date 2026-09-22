package main

import "fmt"

func main() {
	displayTitleScreen()
	playMusic(musicAmbiance)
	fmt.Println()
	fmt.Println("Appuyez sur Entrée pour continuer...")
	readInput()
	clearScreen()
	storyPrologue()
	player := characterCreation()
	clearScreen()
	mainMenu(player)
}

func mainMenu(p *Character) {
	for {
		clearScreen()
		displayBanner(bannerMenu)
 
		lines := []string{
			"\x1b[38;2;100;220;120m" + p.Name + " (" + string(p.Species.Element) + ")\x1b[0m",
			fmt.Sprintf("Niveau %d   PV %s   Mana %s   Or %d",
				p.Level,
				healthBar(p.CurrentHP, p.MaxHP, 12, nil),
				healthBar(p.CurrentMana, p.MaxMana, 12, nil),
				p.Gold),
			"",
			"1. Informations du personnage",
			"2. Inventaire",
			"3. Marchand",
			"4. Forgeron",
			"5. Combat",
			"6. Qui sont-ils ?",
			"7. Mémoires de la vallée",
			"8. Quitter",
		}
		fmt.Print(boxify(lines))
		fmt.Print("Votre choix : ")
 
		choice := readInput()
		clearScreen()

		switch choice {
		case "1":
			displayInfo(p)
		case "2":
			accessInventory(p)
		case "3":
			merchantMenu(p)
		case "4":
			craftMenu(p)
		case "5":
			combatMenu(p)
			if isGameComplete(p) {
				stopMusic()
				return
			}
		case "6":
			showArtists()
		case "7":
			storyJournal()
			continue
		case "8":
			stopMusic()
			fmt.Println("À bientôt !")
			return
		default:
			fmt.Println("Choix invalide. Veuillez entrer un chiffre entre 1 et 8.")
		}

		pauseForInput()
	}
}
