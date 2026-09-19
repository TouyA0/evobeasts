package main

import (
	"fmt"
	"strings"
)

type CreatureArt struct {
	Art    string
	Colors [][3]int
}

func printArt(art string, colors [][3]int) {
	lines := strings.Split(art, "\n")
	for i, line := range lines {
		c := colors[i%len(colors)]
		fmt.Printf("\x1b[38;2;%d;%d;%dm%s\x1b[0m\n", c[0], c[1], c[2], line)
	}
}

var speciesArts = map[string]CreatureArt{}
var monsterArts = map[string]CreatureArt{}

func displayCreature(speciesName string, stage int) {
	key := fmt.Sprintf("%s-%d", speciesName, stage)
	ca, ok := speciesArts[key]
	if !ok {
		fmt.Println("(design manquant pour cette forme)")
		return
	}
	printArt(ca.Art, ca.Colors)
}

func displayMonster(name string) {
	ca, ok := monsterArts[name]
	if !ok {
		fmt.Println("(design manquant)")
		return
	}
	printArt(ca.Art, ca.Colors)
}
