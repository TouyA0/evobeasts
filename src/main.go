package main

import "fmt"

func main() {
	displayTitleScreen()
	fmt.Println()
	fmt.Println("Appuyez sur Entrée pour continuer...")
	readInput()
	clearScreen()
 
	player := characterCreation()
	clearScreen()
 
	playMusic(musicAmbiance)
 
	mainMenu(player)
}

func mainMenu(p *Character) {
	for {
		clearScreen()
 
		lines := []string{
			"\x1b[38;2;100;220;120mEvoBeasts — " + p.Name + " (" + string(p.Species.Element) + ")\x1b[0m",
			fmt.Sprintf("PV %s   Mana %s   Or %d",
				healthBar(p.CurrentHP, p.MaxHP, 12),
				healthBar(p.CurrentMana, p.MaxMana, 12),
				p.Gold),
			"",
			"1. Informations du personnage",
			"2. Inventaire",
			"3. Marchand",
			"4. Forgeron",
			"5. Combat",
			"6. Qui sont-ils ?",
			"7. Quitter",
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
		case "6":
			showArtists()
		case "7":
			stopMusic()
			fmt.Println("À bientôt !")
			return
		default:
			fmt.Println("Choix invalide.")
		}

		pauseForInput()
	}
}
