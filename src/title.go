package main
 
import "fmt"
 
const titleImagePath = "assets/title.png"
const titleWidth = 180

func displayTitleScreen() {
	if !displayImage(titleImagePath, titleWidth) {
		fmt.Println("=== EvoBeasts ===")
	}
}
