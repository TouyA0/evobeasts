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
	imgTaniere		 = "assets/story/taniere.png"
	imgValleeBrumes  = "assets/story/vallee_brumes.png"
	imgVallee        = "assets/story/vallee_ruines.png"
	imgTroisOeufs    = "assets/story/trois_oeufs.png"
	imgDragonVaincu  = "assets/story/dragon_vaincu.png"
	imgParentsPierre = "assets/story/parents_pierre.png"
	imgFinClemence   = "assets/story/fin_clemence.png"
	imgFinRancune    = "assets/story/fin_rancune.png"
)

var storyPhaseLines = map[string]string{
	"Vague Ancestrale": "L'eau des Krakens ruisselle de ses écailles.",
	"Voile Spectral":   "Il disparaît dans le voile des Caméléons.",
	"Convergence":      "Il appelle les trois à la fois. Ses écailles se fendent.",
}

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

func bySpecies(p *Character, phenix, kraken, spectre string) string {
	switch p.Species.Name {
	case EspeceKraken:
		return kraken
	case EspeceSpectre:
		return spectre
	default:
		return phenix
	}
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
		Image: imgValleeBrumes,
		Pages: [][]string{
			{
				"À l'écart du monde, il existe une vallée qu'aucune carte n'a jamais",
				"nommée.",
				"Ceux qui y vivaient l'appelaient la Vallée des Brumes.",
			},
			{
				"Trois familles y régnaient en maîtres.",
				"",
				"  Les Phénix, qui ne connaissaient pas la mort.",
				"  Les Krakens, qui gouvernaient les profondeurs.",
				"  Les Caméléons Spectraux, que personne ne voyait jamais venir.",
			},
			{
				"La brume, c'était eux.",
				"",
				"Le souffle des Phénix sur les eaux des Krakens, et le voile des",
				"Caméléons pour cacher la vallée au reste du monde.",
			},
			{
				"Sous la montagne dormait une créature plus vieille que la vallée :",
				"le Dragon Ancien.",
				"",
				"Il n'avait aucun élément à lui. Il volait celui des autres.",
				"",
				"Une nuit, il se réveilla.",
			},
		},
	})

	playBeat(storyBeat{
		Key:   "prologue-dragon",
		Title: "LA NUIT DU DRAGON",
		Image: imgVallee,
		Pages: [][]string{
			{
				"Il ravagea la vallée.",
				"",
				"Il arracha aux trois familles ce qui les faisait vivre : leur feu,",
				"leur eau, leur voile.",
				"",
				"Puis il redescendit sous la montagne, et les trois gardiens",
				"disparurent avec lui.",
				"Personne ne sait ce qu'il en a fait.",
			},
			{
				"Sans ceux qui la faisaient naître, la brume changea.",
				"Elle devint grise, épaisse, et rendit malade tout ce qui la",
				"respirait.",
			},
		},
	})

	playBeat(storyBeat{
		Key:   "prologue-oeufs",
		Title: "TROIS ŒUFS",
		Image: imgTroisOeufs,
		Pages: [][]string{
			{
				"Il ne laissa derrière lui que quelques bêtes… et trois œufs,",
				"derniers vestiges des grandes familles.",
				"",
				"Trois œufs dont l'élément dormait encore, demeurés invisibles",
				"à ses yeux.",
			},
			{
				"Aujourd'hui, dans le silence, l'un d'eux se met à trembler…",
				"",
				"C'est ainsi que commence l'histoire des Evobeasts.",
			},
		},
	})
}

// ═══════════════════════════════════════════════════════════════════════
//  BEAT 1 — CRÉATION DU PERSONNAGE
// ═══════════════════════════════════════════════════════════════════════

func storyNamePrompt() string {
	return "Dans l'un de ces œufs, quelque chose attend un nom. Lequel ?"
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

func storyCoquille(p *Character) {
	playBeat(storyBeat{
		Key:   "coquille",
		Title: "LA COQUILLE",
		Image: speciesImagePath(p.Species.Name, 1),
		Pages: [][]string{{
			bySpecies(p,
				"Une fissure rougeoie. Un filet de fumée s'en échappe.",
				"Une fissure s'ouvre. Un filet d'eau s'en échappe.",
				"Une fissure apparaît. Disparaît. Réapparaît."),
			"",
			"Pas assez pour sortir. Pas encore.",
			"",
			"Vous avez faim, et personne ne viendra.",
			"",
			"Et dehors, quelque chose gratte contre la coquille.",
		}},
	})
}

func storyEclosion(p *Character) {
	var hatch []string

	switch p.Species.Name {
	case EspeceKraken:
		hatch = []string{
			"La coquille se fend, et l'eau qu'elle contenait se répand sans bruit.",
			p.Name + " ouvre les yeux dans le noir. Ses yeux sont faits pour ça.",
		}
	case EspeceSpectre:
		hatch = []string{
			"Personne n'a vu la coquille s'ouvrir.",
			p.Name + " était simplement là, comme depuis toujours.",
		}
	default:
		hatch = []string{
			"La coquille ne se brise pas : elle se consume.",
			p.Name + " respire pour la première fois, et l'air a un goût de fumée.",
		}
	}

	playBeat(storyBeat{
		Key:   "eclosion",
		Title: "ÉCLOSION",
		Image: speciesImagePath(p.Species.Name, 2),
		Pages: [][]string{
			hatch,
			{
				"Un peu plus loin, les deux autres œufs sont vides.",
				"D'autres que vous ont survécu, quelque part dans la brume.",
				"",
				"La faim vous envahit.",
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
			Pages: [][]string{
				{
					"Un ver de terre.",
					"",
					"Vous avez survécu au Dragon Ancien, et votre premier adversaire",
					"est un ver de terre.",
					"",
					"Il faut bien commencer quelque part.",
				},
				{
					"Ce ver n'est que le premier.",
					"La vallée est pleine de bêtes que la brume a rendues folles.",
					"",
					"Pour survivre, vous allez devoir les affronter…",
				},
			},
		})

	case 2:
		playBeat(storyBeat{
			Key:   "pre-2",
			Title: "CE QUI RESTE DES VIVANTS",
			Image: monsterImagePath(MonstreSanglier),
			Pages: [][]string{{
				"Au loin, une bête aux yeux brûlants de fièvre.",
				"",
				"Elle n'est pas hostile, mais malade. La brume est entrée en elle,",
				"comme dans tout ce qui vit encore ici.",
				"",
				"Elle ne vous laissera pas passer pour autant.",
			}},
		})

	case 3:
		playBeat(storyBeat{
			Key:   "pre-3",
			Title: "LE TÉMOIN",
			Image: monsterImagePath(MonstreLoupGarou),
			Pages: [][]string{
				{
					"Autrefois, c'était un loup. Il gardait l'entrée de la vallée pour",
					"les trois familles.",
					"",
					"La brume l'a changé plus que les autres.",
				},
				{
					"Il parle.",
					"",
					"Ce qui sort de sa gueule n'a plus grand-chose d'un langage,",
					"mais il parle.",
				},
			},
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
		Image: imgTaniere,
		Pages: [][]string{{
			"À l'entrée de la tanière, deux silhouettes.",
			"",
			"Les deux autres œufs. Ils ont grandi, eux aussi.",
			"Ils sont arrivés avant vous. Ils ont essayé.",
			"Leurs blessures disent le reste.",
			"",
			"Sans un mot, ils s'écartent pour vous laisser passer.",
		}},
	})

	playBeat(storyBeat{
		Key:   "pre-4-dragon",
		Title: "LE DRAGON ANCIEN",
		Image: monsterImagePath(MonstreDragonAncien),
		Pages: [][]string{
			{
				"Ils ont appris une chose en perdant :",
				"le Dragon se bat avec ce qu'il a volé.",
				"",
				"        Feu.        Eau.        Spectre.",
				"",
				bySpecies(p,
					"Il va vous attaquer avec la flamme de votre famille.",
					"Il va vous attaquer avec la marée de votre famille.",
					"Il va vous attaquer avec le voile de votre famille."),
			},
			{
				"Mais trois éléments, c'est trop, même pour lui.",
				"",
				"Quand il sera acculé, il essaiera de les porter tous à la fois.",
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
		storyEclosion(p)

	case 2:
		playBeat(storyBeat{
			Key:   "post-2",
			Title: "UNE ÉCAILLE ÉTRANGÈRE",
			Pages: [][]string{{
				"Plantée dans son flanc, une écaille noire.",
				"Bien trop grande pour venir d'une bête de la vallée.",
				"",
				"Quelque chose d'énorme est passé par ici.",
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
					"« Je te reconnais. Tu as l'odeur de ta famille.",
					"",
					"  J'étais là, la nuit où il est remonté. Il n'a pas tué les",
					"  gardiens : il leur a arraché leur élément, et ils se sont",
					"  changés en pierre.",
					"",
					"  Il les a traînés sous la montagne, dans sa tanière.",
					"  Ils respirent encore.",
					"",
					"  Tant qu'il vit, ce qu'il leur a pris existe encore. »",
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
		line = "L'élément qui dormait dans l'œuf vient de s'éveiller."
	case 3:
		line = "La vallée commence à reconnaître ce que " + p.Name + " est en train de devenir."
	case 4:
		line = p.Name + " porte désormais le nom de sa famille : " + currentStageName(p) + "."
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
			"lui. Une rouge. Une bleue. Une qui n'a pas de couleur.",
		}},
	})

	playBeat(storyBeat{
		Key:   "finale-parents",
		Title: "TROIS STATUES",
		Image: imgParentsPierre,
		Pages: [][]string{
			{
				"Au fond de la tanière, trois statues : un Phénix, un Kraken,",
				"un Caméléon. Les trois gardiens.",
				"Parmi eux, celui de votre famille.",
				"",
				"Ils respirent, très lentement.",
			},
			{
				"Les deux autres vous ont rejoint. Ils se tournent vers vous :",
				"c'est vous qui l'avez fait tomber.",
				"",
				"Le Dragon vous regarde. Il ne supplie pas.",
			},
			{
				"« Je ne peux plus les porter.",
				"",
				"  Laisse-moi vivre, et je les rends. Puis je retourne dormir.",
				"",
				"  Achève-moi, et ce que je leur ai pris meurt avec moi.",
				"  Eux aussi.",
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
		"  Il rend ce qu'il a volé.",
		"  Les gardiens se réveillent.",
		"  Il retourne dormir.",
		"  Un jour, il se réveillera.",
		"",
		"  ──────────────────────────────",
		"  Vous sauvez votre famille.",
		"  Vous laissez vivre la menace.",
	}
	rightBody := []string{
		"",
		"Il ne se relèvera jamais.",
		"Ce qu'il a volé meurt avec lui.",
		"Les gardiens ne se réveilleront pas.",
		"Plus rien ne menacera la vallée.",
		"",
		"──────────────────────────────",
		"Vous protégez la vallée.",
		"Vous perdez votre famille.",
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
				return "\n" + centerText("\x1b[1mQUE FAITES-VOUS ?\x1b[0m", totalWidth) + "\n\n" +
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
				"Ce n'est pas de la soumission. C'est autre chose, quelque chose",
				"qui ressemble à de la gratitude, chez une créature qui n'a jamais",
				"eu de mot pour ça.",
			},
			{
				"Les trois lueurs redescendent, chacune vers sa statue.",
				"La pierre se fend. Trois souffles.",
				"",
				"Les gardiens ouvrent les yeux sur trois enfants qu'ils n'ont",
				"jamais vus éclore.",
			},
			{
				"Le Dragon se traîne au fond de sa tanière, et se rendort.",
				"Personne ne sait pour combien de temps.",
			},
			{
				"Dehors, la brume redevient blanche.",
				"",
				"Trois grandes familles se relèveront dans la Vallée des Brumes.",
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
				"Trois souffles, très courts.",
				"",
				"Les gardiens ouvrent les yeux juste assez longtemps pour voir",
				"leurs enfants, debout.",
				"",
				"Puis ils les referment. Pour de bon.",
			},
			{
				p.Name + " et les deux autres remontent à la surface.",
				"Ils ne se retournent pas.",
				"",
				"Au-dessus de la vallée, une brume nouvelle se lève.",
				"La leur, cette fois.",
				"",
				"Trois grandes familles se relèveront dans la Vallée des Brumes.",
				"Et plus rien, jamais, ne descendra de la montagne pour les leur",
				"prendre.",
			},
		},
	})
}

func storyCredits(p *Character) {
	playMusic(musicGloire)
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
	banner, _ := renderImageString(bannerMemoires, bannerWidth)

	if len(storyLog) == 0 {
		clearScreen()
		fmt.Print(banner + "\n")
		fmt.Print(boxify([]string{
			"Mémoires de la vallée",
			"",
			"Vous n'avez encore rien vécu qui mérite d'être raconté.",
		}))
		pauseForInput()
		return
	}

	selected := 0
	for {
		choice := selectMenuFrom(len(storyLog)+1, selected, func(sel int) string {
			labels := make([]string, 0, len(storyLog)+1)
			for _, b := range storyLog {
				title := b.Title
				if title == "" {
					title = "Fragment"
				}
				labels = append(labels, title)
			}
			labels = append(labels, "Retour")
			width := maxWidth(labels...)

			lines := []string{"Mémoires de la vallée", ""}
			for i, label := range labels {
				theme := themeTurquoise
				if i == len(storyLog) {
					theme = themeGris
				}
				lines = append(lines, menuLine(label, width, i == sel, theme))
			}
			return banner + "\n" + strings.Join(lines, "\n")
		})

		if choice == len(storyLog) {
			return
		}
		selected = choice

		replayBeat(storyLog[choice])
	}
}
