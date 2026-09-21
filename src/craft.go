package main

func craftMenu(p *Character)

const (
	ArtefactTete  = "Diadème des Origines"
	ArtefactTorse = "Carapace primordiale"
	ArtefactPieds = "Serres du Voile"

const (
	ResEcaille = "Fragment d'écaille"
	ResAbyssal = "Éclat abyssal"
	ResSpectre = "Poussière spectrale"
	ResPlume   = "Plume calcinée"

	type EquipSlot string
const (
	SlotNone  EquipSlot = ""
	SlotHead  EquipSlot = "tête"
	SlotChest EquipSlot = "torse"
	SlotFeet  EquipSlot = "pieds"


	func AllArtefact(artefact Artefact) {
    if artefact == ArtefactTete {
        fmt.Println(ArtefactTete.Nom)
        fmt.Println(ArtefactTete.Description)
        fmt.Println("PV max :", ArtefactTete.PVMax)
        fmt.Println("Mana max :", ArtefactTete.ManaMax)
    } else if artefact == ArtefactTorse {
        fmt.Println(ArtefactTorse.Nom)
        fmt.Println(ArtefactTorse.Description)
        fmt.Println("PV max :", ArtefactTorse.PVMax)
        fmt.Println("Mana max :", ArtefactTorse.ManaMax)
    } else if artefact == ArtefactPieds {
        fmt.Println(ArtefactPieds.Nom)
        fmt.Println(ArtefactPieds.Description)
        fmt.Println("PV max :", ArtefactPieds.PVMax)
        fmt.Println("Mana max :", ArtefactPieds.ManaMax)
    }
}

func AllArtefact(artefact Artefact) {
    fmt.Println(artefact.Nom)
    fmt.Println(artefact.Description)
    fmt.Println("PV max :", artefact.PVMax)
    fmt.Println("Mana max :", artefact.ManaMax)
}