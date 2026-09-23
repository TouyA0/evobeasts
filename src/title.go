package main
 
import "fmt"
 
const titleImagePath = "assets/ui/title.png"
const titleWidth = 120

func displayTitleScreen() {
	if !displayImage(titleImagePath, titleWidth) {
		fmt.Println("=== EvoBeasts ===")
	}
}
