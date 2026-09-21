package main

import "fmt"

func main() {
	fmt.Println("=== EvoBeasts ===")
	player := characterCreation()
	mainMenu(player)
}

func mainMenu(p *Character) {
	for {
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
			fmt.Println("À bientôt !")
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}
