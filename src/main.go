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
		fmt.Println("\n--- Menu principal ---")
		fmt.Println("1. Afficher les informations du personnage")
		fmt.Println("2. Accéder à l'inventaire")
		fmt.Println("3. Marchand")
		fmt.Println("4. Forgeron")
		fmt.Println("5. Combat")
		fmt.Println("6. Qui sont-ils ?")
		fmt.Println("7. Quitter")
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
