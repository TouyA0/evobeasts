package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	storyCharDelay = 14 * time.Millisecond  // vitesse machine à écrire
	storyLineDelay = 180 * time.Millisecond // pause en fin de ligne
	storyPageWidth = 78                     // largeur des séparateurs
	storyColWidth  = 40                     // largeur d'une colonne du choix final
	storyArtWidth  = 78                     // largeur des images narratives
)

const (
	imgVallee        = "assets/story/vallee_ruines.png"
	imgTroisOeufs    = "assets/story/trois_oeufs.png"
	imgDragonVaincu  = "assets/story/dragon_vaincu.png"
	imgParentsPierre = "assets/story/parents_pierre.png"
	imgFinClemence   = "assets/story/fin_clemence.png"
	imgFinRancune    = "assets/story/fin_rancune.png"
)

var storySeen = map[string]bool{}
var storyLog []storyBeat
type storyBeat struct {
	Key   string
	Title string
	Image string // chemin d'image affiché avant le texte, "" si aucune
	Pages [][]string
}

func storyPrint(lines ...string) {
	for _, line := range lines {
		for _, r := range line {
			fmt.Print(string(r))
			time.Sleep(storyCharDelay)
		}
		fmt.Println()
		time.Sleep(storyLineDelay)
	}
}

func storyHeader(title string) {
	rule := strings.Repeat("─", storyPageWidth)
	fmt.Println("\x1b[38;2;170;150;90m" + rule + "\x1b[0m")
	fmt.Println("\x1b[38;2;235;215;150m" + centerText(title, storyPageWidth) + "\x1b[0m")
	fmt.Println("\x1b[38;2;170;150;90m" + rule + "\x1b[0m")
	fmt.Println()
}

func storyOnce(key string) bool {
	if storySeen[key] {
		return false
	}
	storySeen[key] = true
	return true
}

func playBeat(b storyBeat) {
	if len(b.Pages) == 0 {
		return
	}
	storyLog = append(storyLog, b)

	for i, page := range b.Pages {
		clearScreen()
		if b.Title != "" {
			storyHeader(b.Title)
		}
		if b.Image != "" {
			displayImage(b.Image, storyArtWidth)
			fmt.Println()
		}

		storyPrint(page...)

		if i < len(b.Pages)-1 {
			pauseForInput()
		}
	}
	pauseForInput()
}

func replayBeat(b storyBeat) {
	for i, page := range b.Pages {
		clearScreen()
		if b.Title != "" {
			storyHeader(b.Title)
		}
		if b.Image != "" {
			displayImage(b.Image, storyArtWidth)
			fmt.Println()
		}
		for _, line := range page {
			fmt.Println(line)
		}
		if i < len(b.Pages)-1 {
			pauseForInput()
		}
	}
	pauseForInput()
}

// ═══════════════════════════════════════════════════════════════════════
//  BEAT 0 — PROLOGUE
// ═══════════════════════════════════════════════════════════════════════

func storyPrologue() {
	playBeat(storyBeat{
		Key:   "prologue",
		Title: "LA VALLÉE DES BRUMES",
		Image: imgVallee,
		Pages: [][]string{
			{
				"À l'écart du monde, il existe une vallée que les cartes n'ont",
				"jamais nommée.",
				"Ceux qui y vivaient l'appelaient la Vallée des Brumes.",
			},
			{
				"Trois familles s'y partageaient le silence.",
				"",
				"  Les Phénix, qui ne connaissaient pas la mort.",
				"  Les Krakens, qui régnaient sur les profondeurs.",
				"  Les Caméléons Spectraux, que personne n'a jamais vus venir.",
			},
			{
				"Puis quelque chose de très ancien s'est retourné dans son sommeil.",
				"",
				"Le Dragon Ancien est remonté. Il n'a pas rugi, il n'a pas menacé.",
				"Il a simplement pris.",
				"",
				"Quand la brume est retombée, la vallée n'était plus qu'un champ",
				"de ruines.",
			},
		},
	})

	playBeat(storyBeat{
		Key:   "prologue-oeufs",
		Title: "TROIS ŒUFS",
		Image: imgTroisOeufs,
		Pages: [][]string{
			{
				"Il n'a rien laissé derrière lui. Rien, sauf trois œufs.",
				"",
				"Trois œufs qui ne portaient encore aucun élément, et qu'il n'a",
				"donc pas vus.",
				"",
				"L'un d'eux vient de bouger.",
			},
		},
	})
}

// ═══════════════════════════════════════════════════════════════════════
//  BEAT 1 — CRÉATION DU PERSONNAGE
// ═══════════════════════════════════════════════════════════════════════

func storyNamePrompt() string {
	return "Un nom remonte des cendres de la vallée. Lequel ?"
}

func storyEggLines() []string {
	return []string{
		"1. Œuf ardent   — DPS / Feu",
		"   Il est brûlant, comme une braise qui refuse de s'éteindre.",
		"2. Œuf immergé  — Tank / Eau",
		"   On entend la marée à l'intérieur.",
		"3. Œuf voilé    — Assassin / Spectre",
		"   Vous n'êtes pas certain de l'avoir vraiment vu.",
	}
}

func storyHatch(p *Character) {
	var hatch []string

	switch p.Species.Name {
	case EspecePhenix:
		hatch = []string{
			"La coquille ne se brise pas : elle se consume.",
			p.Name + " respire pour la première fois, et l'air a un goût de fumée.",
		}
	case EspeceKraken:
		hatch = []string{
			"La coquille se fend, et l'eau qu'elle contenait se répand sans bruit.",
			p.Name + " ouvre les yeux dans le noir. Le noir ne le gêne pas.",
		}
	case EspeceSpectre:
		hatch = []string{
			"Personne n'a vu la coquille s'ouvrir.",
			p.Name + " était simplement là, comme s'il y avait toujours été.",
		}
	default:
		hatch = []string{"La coquille se fend."}
	}

	playBeat(storyBeat{
		Key:   "eclosion",
		Title: "ÉCLOSION",
		Image: speciesImagePath(p.Species.Name, 1),
		Pages: [][]string{
			hatch,
			{
				"Vous êtes faible. Vous êtes seul.",
				"",
				"Et quelque part, très loin sous la vallée, quelque chose digère",
				"ce qu'il vous a pris.",
			},
		},
	})
}

// ═══════════════════════════════════════════════════════════════════════
//  BEATS 2 à 5 — AVANT CHAQUE PREMIER COMBAT
// ═══════════════════════════════════════════════════════════════════════

func storyPreFight(p *Character, tier int) {
	if !storyOnce(fmt.Sprintf("pre-%d", tier)) {
		return
	}

	switch tier {
	case 1:
		playBeat(storyBeat{
			Key:   "pre-1",
			Title: "PREMIÈRE SORTIE",
			Image: monsterImagePath(MonstreVer),
			Pages: [][]string{{
				"Un ver de terre.",
				"",
				"Vous avez survécu au Dragon Ancien, et votre premier adversaire",
				"est un ver de terre.",
			}},
		})

	case 2:
		playBeat(storyBeat{
			Key:   "pre-2",
			Title: "CE QUI RESTE DES VIVANTS",
			Image: monsterImagePath(MonstreSanglier),
			Pages: [][]string{{
				"La bête a les yeux voilés de blanc.",
				"",
				"Elle n'est pas hostile : elle est malade. Ce que le Dragon a",
				"réveillé s'est répandu dans tout ce qui vivait ici.",
				"",
				"Elle ne vous laissera pas passer pour autant.",
			}},
		})

	case 3:
		playBeat(storyBeat{
			Key:   "pre-3",
			Title: "LE TÉMOIN",
			Image: monsterImagePath(MonstreLoupGarou),
			Pages: [][]string{{
				"Il parle.",
				"",
				"Ce qui sort de sa gueule n'a plus grand-chose d'un langage,",
				"mais il parle.",
				"",
				"Il vivait ici, avant.",
			}},
		})

	case 4:
		storyPreDragon(p)
	}
}

func storyPreDragon(p *Character) {
	stopMusic()

	playBeat(storyBeat{
		Key:   "pre-4",
		Title: "LA TANIÈRE ANCESTRALE",
		Image: monsterImagePath(MonstreDragonAncien),
		Pages: [][]string{
			{
				"Deux silhouettes vous attendent à l'entrée de la tanière.",
				"",
				"Les deux autres œufs. Ils ont grandi eux aussi.",
				"Aucun des trois ne dit un mot. Ce n'est pas nécessaire.",
				"",
				"En bas, il vous attend.",
			},
			{
				"Le Dragon Ancien ne possède pas d'élément.",
				"Il copie ceux des autres, et il ne les rend jamais.",
				"",
				"Mais il a une limite : il ne peut en contenir que trois.",
				"",
				"        Feu.        Eau.        Spectre.",
				"",
				"Trois éléments arrachés à trois créatures qui, quelque part sous",
				"vos pattes, respirent encore sans savoir pourquoi.",
			},
			{
				"Quand il n'aura plus le choix, il essaiera de porter les trois",
				"à la fois.",
				"",
				"C'est à ce moment-là qu'il se brisera.",
			},
		},
	})
}

// ═══════════════════════════════════════════════════════════════════════
//  BEATS 2 à 4 — APRÈS CHAQUE PREMIÈRE VICTOIRE
// ═══════════════════════════════════════════════════════════════════════

func storyPostFight(p *Character, tier int) {
	if tier >= 4 {
		return
	}
	if !storyOnce(fmt.Sprintf("post-%d", tier)) {
		return
	}

	switch tier {
	case 1:
		playBeat(storyBeat{
			Key:   "post-1",
			Title: "SOUS LA PEAU",
			Pages: [][]string{{
				"Ce n'était rien.",
				"",
				"Mais vous sentez quelque chose bouger sous votre peau.",
			}},
		})

	case 2:
		playBeat(storyBeat{
			Key:   "post-2",
			Title: "UNE ÉCAILLE ÉTRANGÈRE",
			Pages: [][]string{{
				"Sur son flanc, une écaille qui n'est pas la sienne.",
				"",
				"Vous la reconnaissez sans savoir pourquoi.",
			}},
		})

	case 3:
		playBeat(storyBeat{
			Key:   "post-3",
			Title: "CE QU'IL A VU",
			Image: monsterImagePath(MonstreLoupGarou),
			Pages: [][]string{
				{
					"Il ne cherche plus à se relever.",
				},
				{
					"« Je les ai vus. Tes vieux. Les trois.",
					"",
					"  Il ne les a pas tués — il les a vidés. Il leur a pris ce qui",
					"  les tient debout, et il est reparti sous la montagne avec.",
					"",
					"  Ils sont encore en bas. Encore en vie, si on peut appeler",
					"  ça vivre.",
					"",
					"  Le sort tient tant qu'il tient, lui. »",
				},
				{
					"Il ne dit plus rien.",
					"",
					"Pour la première fois, vous savez où aller.",
				},
			},
		})
	}
}

// ═══════════════════════════════════════════════════════════════════════
//  BEAT — ÉVOLUTIONS
// ═══════════════════════════════════════════════════════════════════════

func storyEvolve(p *Character) {
	var line string

	switch p.Stage {
	case 2:
		line = "Quelque chose se réveille dans " + p.Name + ". Ce n'est pas encore un nom."
	case 3:
		line = "La vallée commence à reconnaître ce que " + p.Name + " est en train de devenir."
	case 4:
		line = p.Name + " n'a plus besoin qu'on lui dise ce qu'il est. Il est " +
			currentStageName(p) + "."
	default:
		return
	}

	if !storyOnce(fmt.Sprintf("evo-%d", p.Stage)) {
		return
	}

	fmt.Println()
	storyPrint("\x1b[38;2;200;180;255m" + line + "\x1b[0m")
}

// ═══════════════════════════════════════════════════════════════════════
//  BEAT 6 — LE CHOIX FINAL ET LES DEUX FINS
// ═══════════════════════════════════════════════════════════════════════

func storyFinale(p *Character) {
	stopMusic()

	playBeat(storyBeat{
		Key:   "finale-scene",
		Title: "CE QU'IL RESTE DU DRAGON",
		Image: imgDragonVaincu,
		Pages: [][]string{{
			"Il ne peut plus se relever.",
			"",
			"Trois lueurs s'échappent de ses plaies et flottent au-dessus de",
			"lui, sans tomber.",
			"Une rouge. Une bleue. Une qui n'a pas de couleur.",
		}},
	})

	playBeat(storyBeat{
		Key:   "finale-parents",
		Title: "TROIS CORPS DE PIERRE",
		Image: imgParentsPierre,
		Pages: [][]string{
			{
				"Au fond de la salle, trois corps de pierre.",
				"Ils respirent, très lentement.",
				"",
				"Le Dragon vous regarde. Il ne supplie pas.",
			},
			{
				"« Si tu m'achèves, ce que je leur ai pris meurt avec moi.",
				"",
				"  Si tu me laisses, je le rends, et je me rendors.",
				"",
				"  Choisis. Je n'ai plus rien à voler. »",
			},
		},
	})

	if storyChoiceScreen() == "1" {
		storyEndingClemence(p)
	} else {
		storyEndingRancune(p)
	}

	storyCredits(p)
}

func storyChoiceScreen() string {
	totalWidth := storyColWidth*2 + 3

	leftBody := []string{
		"",
		"  Il perd son sort draconique.",
		"  Les trois éléments retournent",
		"  à leurs porteurs.",
		"",
		"  Tes parents se réveillent.",
		"",
		"  Il se rendort.",
		"  Et un jour, il se réveillera.",
		"",
		"  ──────────────────────────────",
		"  Tu leur rends leur vie.",
		"  Tu ne leur rends pas justice.",
	}
	rightBody := []string{
		"",
		"Il ne se relèvera plus jamais.",
		"Les trois éléments s'éteignent",
		"avec lui.",
		"",
		"Tes parents ne se réveillent pas.",
		"",
		"Plus jamais personne ne prendra",
		"rien à ta famille.",
		"",
		"──────────────────────────────",
		"Tu leur rends justice.",
		"Tu ne leur rends pas leur vie.",
	}

	// Titre de colonne : barre de couleur si sélectionné, texte coloré sinon.
	header := func(label string, selected bool, t menuTheme) string {
		if selected {
			return menuLine(label, 0, true, t)
		}
		return "\x1b[38;2;" + t.bg + "m" + menuLine(label, 0, false, t) + "\x1b[0m"
	}

	choice := selectMenuHint(2, 0, menuHintSide, func(sel int) string {
		left := strings.Join(append([]string{"  " + header("ÉPARGNER", sel == 0, themeVert)}, leftBody...), "\n")
		right := strings.Join(append([]string{header("ACHEVER", sel == 1, themeRouge)}, rightBody...), "\n")
		return "\n" + centerText("\x1b[1mQUE FAIS-TU ?\x1b[0m", totalWidth) + "\n\n" +
			dividedColumns(left, right, storyColWidth)
	})

	if choice == 0 {
		return "1"
	}
	return "2"
}

func storyEndingClemence(p *Character) {
	playMusic(musicAmbiance)

	playBeat(storyBeat{
		Key:   "fin-1",
		Title: "FIN I — LA CLÉMENCE",
		Image: imgFinClemence,
		Pages: [][]string{
			{
				"Le Dragon incline la tête.",
				"",
				"Ce n'est pas de la soumission, c'est autre chose — quelque chose",
				"qui ressemble à de la gratitude, chez une créature qui n'a jamais",
				"eu de mot pour ça.",
			},
			{
				"Les trois lueurs redescendent.",
				"La pierre se fend. Trois souffles.",
				"",
				"Il se traîne jusqu'au fond de sa tanière ancestrale, et se rendort.",
				"",
				"Personne ne saura jamais combien de temps.",
			},
			{
				"Trois familles se relèveront dans la Vallée des Brumes.",
				"",
				"Et trois enfants dormiront mal, chaque nuit, en écoutant la",
				"montagne.",
			},
		},
	})
}

func storyEndingRancune(p *Character) {
	// Pas de musique sur cette fin

	playBeat(storyBeat{
		Key:   "fin-2",
		Title: "FIN II — LA RANCUNE",
		Image: imgFinRancune,
		Pages: [][]string{
			{
				"Le Dragon ne dit rien de plus. Il n'en a pas le temps.",
				"",
				"Les trois lueurs s'éteignent au moment exact où il cesse de",
				"respirer.",
			},
			{
				"Dans la salle, la pierre se fend quand même.",
				"Trois souffles — courts.",
				"",
				"Ils ont juste le temps de vous voir.",
				"C'est probablement mieux ainsi.",
			},
			{
				p.Name + " et les deux autres remontent à la surface.",
				"Ils ne se retournent pas.",
				"",
				"Trois familles se relèveront dans la Vallée des Brumes, et à leur",
				"tête se tiendront trois créatures que plus personne ne prendra",
				"jamais pour des œufs.",
				"",
				"Elles n'ont rien pardonné. Elles n'ont rien oublié.",
				"Et c'est très bien comme ça.",
			},
		},
	})
}

func storyCredits(p *Character) {
	clearScreen()
	storyHeader("EVOBEASTS")

	if path := speciesImagePath(p.Species.Name, p.Stage); path != "" {
		displayImage(path, storyArtWidth)
		fmt.Println()
	}

	fmt.Println(centerText(p.Name+" — "+currentStageName(p), storyPageWidth))
	fmt.Println()
	storyPrint("Merci d'avoir joué à EvoBeasts.")
	pauseForInput()
}

// ═══════════════════════════════════════════════════════════════════════
//  MÉMOIRES DE LA VALLÉE — relire les scènes déjà vues
// ═══════════════════════════════════════════════════════════════════════

func storyJournal() {
	for {
		clearScreen()
		displayBanner(bannerMemoires)

		if len(storyLog) == 0 {
			fmt.Print(boxify([]string{
				"Mémoires de la vallée",
				"",
				"Vous n'avez encore rien vécu qui mérite d'être raconté.",
			}))
			return
		}

		lines := []string{"Mémoires de la vallée", ""}
		for i, b := range storyLog {
			title := b.Title
			if title == "" {
				title = "Fragment"
			}
			lines = append(lines, fmt.Sprintf("%d. %s", i+1, title))
		}
		lines = append(lines, "", "0. Retour")

		fmt.Print(boxify(lines))
		fmt.Print("Votre choix : ")

		choice := readInput()
		if choice == "0" {
			return
		}

		index := -1
		fmt.Sscanf(choice, "%d", &index)
		if index < 1 || index > len(storyLog) {
			continue
		}

		replayBeat(storyLog[index-1])
	}
}
