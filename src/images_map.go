package main

import "fmt"

const creatureImageDir = "assets/creatures/"

var speciesImageFile = map[string]string{
	EspecePhenix:  "phenix",
	EspeceKraken:  "kraken",
	EspeceSpectre: "spectre",
}

var monsterImageFile = map[string]string{
	MonstreVer:          "ver",
	MonstreSanglier:     "sanglier",
	MonstreLoupGarou:    "loup_garou",
	MonstreDragonAncien: "dragon_ancien",
}

func speciesImagePath(speciesName string, stage int) string {
	prefix, ok := speciesImageFile[speciesName]
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s%s_%d.png", creatureImageDir, prefix, stage)
}

func monsterImagePath(monsterName string) string {
	file, ok := monsterImageFile[monsterName]
	if !ok {
		return ""
	}
	return creatureImageDir + file + ".png"
}