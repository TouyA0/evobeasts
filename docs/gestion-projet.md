# EvoBeasts — Gestion de projet

## Les branches

| Branche | Qui | Domaine |
|---|---|---|
| `feature/character-base` | Lucas B | personnage, inventaire |
| `feature/merchant` | Kelissa | marchand, économie |
| `feature/crafting` | Lucas H | forgeron, équipement |
| `feature/combat` | Quentin | combat, ennemis, sorts, évolution, affichage |

Règle : on ne code jamais directement sur `main`. Chacun travaille sur sa
branche et on merge sur `main` tous ensemble pour régler les conflits
proprement.

---

## 1. Le jeu en quelques phrases

La Vallée des Brumes était protégée par trois familles : les Phénix, les
Krakens et les Caméléons Spectraux. Le Dragon Ancien leur a volé leurs
éléments et changé leurs gardiens en pierre. Il ne reste que trois œufs.

On incarne l'un d'eux. On le fait évoluer en battant quatre ennemis de plus
en plus forts, jusqu'au Dragon. À la fin, on choisit de l'épargner ou de
l'achever, et chaque choix mène à une fin différente.

Combat au tour par tour, en 1 contre 1, dans le terminal.

---

## 2. Lancer le jeu

```bash
cd src
go build -o evobeasts.exe
ouvrir evobeasts.exe
```

Il faut aussi un terminal qui gère la couleur 24 bits (Windows Terminal,
le terminal de VS Code).

---

## 3. L'arborescence

```
src/
├── go.mod
├── main.go             COMMUN   → lancement + menu principal
├── types.go            COMMUN   → toutes les structs et constantes
├── menu.go             COMMUN   → menus au clavier (flèches, Entrée, Échap)
│
├── character.go        Lucas B  → création, fiche perso, stats
├── inventory.go        Lucas B  → inventaire + fonctions partagées
│
├── merchant.go         Kelissa  → boutique
├── economy.go          Kelissa  → or, XP, niveaux
│
├── craft.go            Lucas H  → forgeron
├── equipment.go        Lucas H  → artefacts équipés
│
├── combat.go           Quentin  → boucle de combat, dégâts, patterns
├── combat_screen.go    Quentin  → écran de combat
├── monsters.go         Quentin  → les 4 ennemis + phases du boss
├── spells.go           Quentin  → sorts
├── evolution.go        Quentin  → espèces, paliers, évolutions
├── imageart.go         Quentin  → affichage des images PNG dans le terminal
├── images_map.go       Quentin  → quelle image pour quelle créature
│
├── story.go            COMMUN   → histoire, choix final, Mémoires
├── audio.go            COMMUN   → musiques
├── banners.go          COMMUN   → bannières des menus
├── title.go            COMMUN   → écran titre
├── flush_windows.go    COMMUN   → vide le buffer clavier (Windows)
├── flush_other.go      COMMUN   → version vide pour les autres systèmes
│
└── assets/
    ├── creatures/      les 3 espèces × 4 paliers + les 4 ennemis
    ├── story/          images des scènes
    ├── ui/             écran titre + bannières
    └── audio/          ambiance, combat, boss, gloire
```

---

## 4. Qui fait quoi

**Lucas B — le personnage et la base**
: création du perso, fiche d'infos, inventaire,
menu, mort, limite d'inventaire.

**Kelissa — l'économie**
: marchand, prix, or, potion de poison,
agrandissement de l'inventaire, XP et montée de niveau.

**Lucas H — le craft**
: forgeron, recettes, équipement, bonus des
artefacts.

**Quentin — le combat**
: combat, sorts, les quatre ennemis,
le boss à quatre phases, l'évolution, l'écran de combat et l'affichage des
créatures.

---

## 5. Les règles à respecter

### Les stats dérivées passent par `recalculateStats()`

**Personne n'écrit `p.MaxHP += X`.** On modifie la source (`Stage`, `Level`,
`Equipment`), puis on appelle `recalculateStats()` qui recalcule tout depuis
zéro :

```
MaxHP = BaseMaxHP
      + HPPerStage   × (Stage − 1)
      + GainPVNiveau × (Level − 1)
      + bonus PV des artefacts équipés
```
Pareil pour mana / atk.

### Constantes / struct = types.go

---

## 6. Les deux progressions

### L'évolution : le fil de l'histoire
Le palier monte à la **première victoire** contre chaque ennemi :

| Ennemi battu | Palier atteint |
|---|---|
| Ver | 2 — Poussin / Alevin / Caméléonneau |
| Sanglier | 3 — Poulet / Poisson / Caméléon |
| Loup-garou | 4 — Phénix / Kraken / Caméléon Spectral |
| Dragon Ancien | fin du jeu |

### Le niveau :
On gagne de l'XP à chaque combat gagné, même contre un ennemi déjà battu. On
peut donc farmer sans casser l'histoire.

- 100 XP pour le niveau 2, puis +50 à chaque niveau (150, 200, 250...)
- Par niveau : +5 PV max, +3 mana max, +1 attaque, et PV/mana remis au max
- L'XP en trop passe au niveau suivant

---

## 7. Les trois espèces

| Espèce | Élément | Rôle | PV | PV/palier | Mana | Mana/palier | Atk | Atk/palier | Init | Esquive |
|---|---|---|---|---|---|---|---|---|---|---|
| Phénix | Feu | DPS | 50 | +20 | 44 | +18 | 6 | +3 | 12 | 5 % |
| Kraken | Eau | Tank | 70 | +25 | 36 | +17 | 6 | +2 | 8 | 5 % |
| Caméléon Spectral | Spectre | Assassin | 48 | +16 | 32 | +15 | 6 | +4 | 18 | 20 % |

Au palier 4 (sans niveaux/artefacts) :

| Espèce | PV max | Mana max | Attaque |
|---|---|---|---|
| Phénix | 110 | 98 | 15 |
| Kraken | 145 | 87 | 12 |
| Caméléon Spectral | 96 | 77 | 18 |

### Le triangle des éléments
**L'Eau bat le Feu, le Feu bat le Spectre, le Spectre bat l'Eau.** Chaque
espèce résiste à son propre élément et craint celui qui la bat.
Multiplicateur appliqué aux dégâts reçus :

| Espèce | Feu | Eau | Spectre | Neutre |
|---|---|---|---|---|
| Phénix | 0.6 | **1.35** | 1.0 | 1.0 |
| Kraken | 1.0 | 0.6 | **1.35** | 1.0 |
| Caméléon Spectral | **1.35** | 1.0 | 0.6 | 1.0 |

### Au départ
- PV à 50 % du max, mana plein
- 100 or et 2 Essences vitales dans le sac

---

## 8. Les sorts

Au palier 1 on n'a que l'attaque de base, qui prend l'élément de l'espèce.
Chaque évolution apprend automatiquement les sorts du nouveau palier.

**Les buffs d'attaque s'appliquent aussi aux sorts**, pas seulement à
l'attaque de base.

### Phénix — Feu : brûlures et soins

| Sort | Palier | Mana | Effet |
|---|---|---|---|
| Picorage enflammé | 2 | 8 | 10 dégâts + brûlure 4/tour pendant 2 tours |
| Plumage chaud | 2 | 8 | rend 24 PV |
| Coup de bec ardent | 3 | 11 | 24 dégâts |
| Cocorico enflammé | 3 | 9 | +7 dégâts pendant 4 tours |
| Griffes de braise | 4 | 15 | 34 dégâts |
| Nova solaire | 4 | 19 | 16 dégâts + brûlure 7/tour pendant 4 tours |
| **Renaissance** (ultime) | 4 | — | passif : à la mort, revient à 40 % PV avec +8 d'attaque pendant 4 tours |

### Kraken — Eau : encaisser et contrôler

| Sort | Palier | Mana | Effet |
|---|---|---|---|
| Éclaboussure | 2 | 6 | 12 dégâts + initiative adverse −6 pendant 3 tours |
| Écailles naissantes | 2 | 7 | −30 % de dégâts subis pendant 3 tours |
| Coup de nageoire lourd | 3 | 11 | 22 dégâts + attaque adverse −4 pendant 3 tours |
| Bulle protectrice | 3 | 11 | bouclier de 35 |
| Étreinte tentaculaire | 4 | 17 | 22 dégâts + l'adversaire perd son prochain tour |
| Raz-de-marée | 4 | 16 | 28 dégâts + initiative adverse −10 pendant 3 tours |
| **Abysse** (ultime) | 4 | 24 | 13 dégâts/tour pendant 5 tours |

### Caméléon Spectral — Spectre : frapper depuis l'invisibilité

| Sort | Palier | Mana | Effet |
|---|---|---|---|
| Morsure furtive | 2 | 7 | 15 dégâts |
| Camouflage | 2 | 11 | invisible : +40 % d'esquive, prochaine attaque +130 % |
| Frappe chromatique | 3 | 13 | 30 dégâts |
| Voile chromatique | 3 | 15 | invisible : +40 % d'esquive, prochaine attaque +170 % |
| Lame spectrale | 4 | 15 | 32 dégâts, ignore les résistances |
| **Exécution spectrale** (ultime) | 4 | 20 | 35 dégâts, **110** depuis l'invisibilité, ignore les résistances |

Invisibilité : la prochaine attaque, que ce
soit l'attaque de base ou un sort, prend le bonus puis on redevient visible.

Ignore les résistances : la cible ne peut pas réduire les
dégâts. Si elle est faible à l'élément, le bonus s'applique quand même.

### Le sort acheté
**Souffle primordial** : grimoire à 35 or, 13 mana, 24 dégâts, élément
neutre. Pas de résistance à l'élément neutre. On peut l'apprendre qu'une fois.

### Ultimes
Un seul ultime par combat.

---

## 9. Les quatre ennemis

| Ennemi | Tier | PV | Atk | Init | XP | Or | Loot |
|---|---|---|---|---|---|---|---|
| Ver | 1 — Novice | 15 | 4 | 3 | 40 | 6 | Plume calcinée 70 % |
| Sanglier | 2 — Intermédiaire | 55 | 7 | 10 | 90 | 15 | Fragment d'écaille 75 % |
| Loup-garou | 3 — Avancé | 90 | 11 | 14 | 150 | 30 | Poussière spectrale 70 %, Fragment d'écaille 50 % |
| Dragon Ancien | 4 — Boss | 230 | 17 | 16 | 300 | 80 | Éclat abyssal 100 % |

Trois premiers ennemis : neutres et sans résistances. Seul le dragon joue avec les éléments.
Si le sac est plein alors le loot est perdu.

### Leurs patterns

**Ver** : morsure à 100 %, tous les tours. C'est le tutoriel.

**Sanglier**
- Charge tous les 3 tours : 140 % + 35 % de chances d'étourdir
- Enragé : dégâts ×1,3 sous 50 % PV, ×1,6 sous 25 %

**Loup-garou**
- Griffure enragée les tours pairs : 2 coups à 65 %
- Hurlement au tour 3, une seule fois : +4 d'attaque jusqu'à la fin du combat
- Régénération lunaire : +5 PV par tour

**Dragon Ancien, le boss en quatre phases**

Il change d'élément à chaque phase. **Chaque
phase a une faiblesse qui correspond à une des trois espèces**, pour que
tout le monde ait un moment où il domine et un moment où il galère.

| Phase | PV | Élément | Feu / Eau / Spectre | Ce qu'il fait |
|---|---|---|---|---|
| 1 — Éveil Igné | 100 → 75 % | Feu | 0.5 / **1.5** / 1.0 | Souffle ardent (125 % + brûlure 6×2) un tour sur deux, sinon Griffe (105 %) |
| 2 — Vague Ancestrale | 75 → 50 % | Eau | 1.0 / 0.5 / **1.5** | Étreinte liquide (90 % + perte de tour) tous les 3 tours, sinon Raz-de-marée (115 %). Régénère 6 PV/tour |
| 3 — Voile Spectral | 50 → 25 % | Spectre | **1.5** / 1.0 / 0.5 | Dissolution (esquive la prochaine attaque) tous les 4 tours, sinon Frappe hors du plan (135 %, détruit le bouclier) |
| 4 — Convergence | 25 → 0 % | les trois | 1.0 / 1.0 / 1.0 | 135 % +30 %, l'élément change à chaque tour (Feu → Eau → Spectre) |

### Le menu de combat
Chaque victoire déverouille l'ennemi suivant. Les ennemis déjà
battus restent rejouables pour farmer.

---

## 10. Les règles de combat

- Celui qui a la plus haute initiative joue en premier. On compte les malus
  d'initiative, et l'écran indique qui commence.
- Dégâts = (base + buffs − malus) × résistance de la cible
- Côté joueur : l'esquive est testée en premier, puis la réduction de
  dégâts, puis le bouclier, et enfin les PV.
- Les effets (brûlure, poison, régénération) s'appliquent en fin de tour.
- Relancer un buff remet sa durée à zéro. Relancer une brûlure la cumule.
- **Un ennemi qui vient de perdre son tour ne peut pas le reperdre au tour
  suivant**.
- Le mana ne remonte pas pendant le combat, sauf avec les Essences de mana.
- **Victoire** : or, XP, loot, évolution si c'est la première victoire, puis
  PV et mana remis au max.
- **Défaite** : le combat s'arrête tout de suite. Un Phénix au
  palier 4 revient d'abord avec Renaissance. Sinon, `isDead` le remet à
  **1 PV** et on retourne au menu. Il faut se soigner avant de retenter,
  ce qui rend les potions utiles.

### L'écran de combat
- l'**intention** de l'ennemi pour le tour (quelle attaque arrive)
- ses faiblesses et ses résistances actuelles
- pour le Dragon : les seuils de phase sur la barre de vie
- la liste de ses compétences, avec leurs pourcentages
- nos effets actifs, notre bouclier, notre invisibilité et notre ultime

Dans le menu des sorts, un sort qu'on ne peut pas lancer s'affiche en rouge
avec la raison.

---

## 11. Le marchand et le forgeron

### Le marchand

| Catégorie | Objet | Effet | Prix |
|---|---|---|---|
| Consommables | Essence vitale | +30 PV | 12 or |
| | Essence de mana | +25 mana | 10 or |
| | Essence corrompue | poison, 8 dégâts × 4 | 10 or |
| Matériaux | Plume calcinée | ressource | 2 or |
| | Poussière spectrale | ressource | 4 or |
| | Fragment d'écaille | ressource | 5 or |
| | Éclat abyssal | ressource | 12 or |
| Objets rares | Grimoire : Souffle primordial | apprend le sort | 35 or |
| Améliorations | Expansion de besace | +10 places, 3 fois max | 30 or |

Ce qu'on ne peut pas payer est grisé. On voit combien on en possède déjà, et
combien d'expansions ont été utilisées. L'or n'est retiré que si l'achat a
bien marché.

### L'Essence corrompue
Le sujet demande une potion qui fait des dégâts avec une seconde entre chaque
tick. On lui a donné un usage différent selon le moment :
- **en combat** : on la lance sur l'ennemi, qui prend 8 dégâts par tour
  pendant 4 tours
- **hors combat** : on la boit soi-même après une confirmation (« Non » est
  sélectionné par défaut). On prend 4 ticks de 8 dégâts, une seconde
  d'écart entre chaque

### Le forgeron
Chaque fabrication coûte 10 or en plus des ressources.

| Artefact | Emplacement | Ressources | Bonus |
|---|---|---|---|
| Diadème des Origines | tête | 1 Plume calcinée + 1 Fragment d'écaille | +6 PV, +16 mana |
| Carapace primordiale | torse | 2 Éclats abyssaux + 1 Fragment d'écaille | +22 PV |
| Serres du Voile | pieds | 1 Poussière spectrale + 1 Fragment d'écaille | +8 PV, +8 mana |

On équipe depuis l'inventaire. Si l'emplacement est déjà pris, l'ancien
artefact revient dans le sac. Il y a forcément une place libre, puisque le
nouvel artefact vient d'en sortir.

---

## 12. L'inventaire

- Chaque exemplaire prend sa propre place : trois Essences vitales occupent
  trois places. À l'affichage, elles sont regroupées en « x3 ».
- 10 places au départ, +10 trois fois, soit 40 maximum
- Les objets sont rangés par catégorie : Potions, Grimoires, Artefacts,
  Ressources
- Toute fonction qui ajoute un objet vérifie d'abord s'il reste de la place

---

## 13. L'expérience de jeu

- **Menus au clavier** : ↑/↓ (ou ←/→) pour se déplacer, Entrée pour valider,
  Échap pour revenir, chiffres comme raccourcis, Ctrl+C pour quitter.
- **Buffer clavier vidé** avant chaque menu (sous Windows)
- **Images** : les créatures, les ennemis, les scènes et les bannières sont
  des PNG affichés directement dans le terminal en blocs colorés. Chaque
  palier de chaque espèce a son image.
- **Musique** : ambiance, combat, boss et gloire (générique).
- **Histoire** : prologue, une scène avant et après chaque premier combat,
  une phrase à chaque évolution, puis le choix final :
  - **Épargner** → *Fin I — La Clémence* : les gardiens se réveillent,
    mais le Dragon dort toujours sous la montagne
  - **Achever** → *Fin II — La Rancune* : la vallée est sauvée, mais les
    gardiens meurent avec le Dragon
- **Mémoires de la vallée** : un menu pour relire les scènes déjà vues.

---

## 14. Mission 6 — les deux artistes cachés

Menu « Qui sont-ils ? » :
- **Partie 2** : *Money, money, money*, *Gimme! Gimme! Gimme!*,
  *Mamma Mia*... → **ABBA**
- **Partie 3** : *Duel*, *A.I. Intelligence artificielle*,
  *Ready Player One*... → **Steven Spielberg**

---

## 15. Les fonctions par fichier

Les fonctions marquées **partagée** sont utilisées par d'autres fichiers :
on ne change pas leur nom ni leur signature sans prévenir tout le monde.

### main.go — commun
| Fonction | Rôle |
|---|---|
| `main()` | écran titre, prologue, création du perso, menu |
| `mainMenu(p)` | menu principal, un `case` par module |

### types.go — commun
Aucune fonction. Uniquement les structs (`Character`, `Species`, `Spell`,
`Item`, `Equipment`, `Recipe`, `LootEntry`, `StatusEffect`, `Monster`,
`MonsterSkill`, `Phase`) et toutes les constantes.

### menu.go — commun
| Fonction | Rôle |
|---|---|
| `selectMenu` / `selectMenuFrom` / `selectMenuHint` | **partagée** — menu plein écran au clavier |
| `inlineSelect(...)` | **partagée** — menu sous l'écran actuel (utilisé en combat) |
| `confirm(lines...) bool` | **partagée** — question Oui / Non |
| `button` / `backButton` | **partagée** — bouton qui attend Entrée |
| `menuLine`, `menuLines`, `maxWidth` | mise en forme des lignes de menu |

### character.go — Lucas B
| Fonction | Rôle |
|---|---|
| `characterCreation() *Character` | demande le nom et l'œuf |
| `isLettersOnly(s) bool` | vérifie que le nom ne contient que des lettres |
| `formatName(s) string` | première lettre en majuscule, le reste en minuscules |
| `initCharacter(name, species)` | construit le perso de départ |
| `recalculateStats(p)` | **partagée** — recalcule MaxHP, MaxMana, Attack |
| `isDead(p) bool` | **partagée** — gère le 0 PV |
| `readInput() string` | **partagée** — lecture clavier |
| `displayInfo(p)` | fiche perso : image, stats, équipement, sorts |

### inventory.go — Lucas B
| Fonction | Rôle |
|---|---|
| `addInventory(p, item) bool` | **partagée** — ajoute si place, sinon `false` |
| `removeInventory(p, name) bool` | **partagée** — retire un exemplaire |
| `countItem(p, name) int` | **partagée** — compte les exemplaires |
| `hasSpace(p) bool` | **partagée** — reste-t-il une place ? |
| `accessInventory(p)` | écran d'inventaire |
| `inventoryEntries(items)` | regroupe par catégorie et par nom |
| `itemDescription(item)` | texte d'effet d'un objet |
| `useItem(p, item)` | utilise un objet hors combat |
| `usePotion(p, item)` | effet d'une potion |
| `poisonPot(p, dmg, ticks) bool` | boire l'Essence corrompue, avec confirmation |
| `takePot(p)` | potion de vie (tâche du sujet, pas reliée au menu) |

### merchant.go — Kelissa
| Fonction | Rôle |
|---|---|
| `merchantMenu(p)` | boucle de la boutique |
| `merchantCatalog() []Item` | tout ce qui est vendu |
| `buyItem(p, item) bool` | vérifie l'or et la place, puis retire l'or |
| `upgradeInventorySlot(p) bool` | +10 places, 3 fois maximum |
| `merchantOrder`, `merchantScreen` | classement par catégorie et affichage |

### economy.go — Kelissa
| Fonction | Rôle |
|---|---|
| `gainXP(p, amount)` | **partagée** — ajoute l'XP et gère les niveaux |
| `gainGold(p, amount)` | **partagée** — ajoute l'or |
| `levelUp(p)` | monte d'un niveau, recalcule les stats, soigne |
| `xpForNextLevel(level) int` | 100, puis +50 par niveau |

### craft.go — Lucas H
| Fonction | Rôle |
|---|---|
| `craftMenu(p)` | boucle du forgeron |
| `recipes() []Recipe` | les trois recettes |
| `canCraft(p, r) (bool, string)` | vérifie or, ressources et place, et donne la raison si c'est non |
| `craftItem(p, r) bool` | consomme les ressources et ajoute l'artefact |
| `craftScreen(...)` | affiche les recettes et ce qui manque |
| `showArtists()` | Mission 6 |

### equipment.go — Lucas H
| Fonction | Rôle |
|---|---|
| `equipItem(p, item) bool` | équipe, remet l'ancien dans le sac, recalcule les stats |
| `unequipItem(p, slot) bool` | retire un artefact s'il y a de la place (pas relié au menu) |
| `equipmentBonus(e) (hp, mana)` | somme des bonus, utilisée par `recalculateStats` |
| `equipmentDisplay(e)` | affiche les trois emplacements |

### combat.go — Quentin
| Fonction | Rôle |
|---|---|
| `combatMenu(p)` | liste des ennemis débloqués |
| `trainingFight(p, m)` | boucle du combat : ordre de jeu, tours, effets, fin |
| `checkPlayerDefeat(p) bool` | Renaissance, sinon défaite + `isDead` |
| `runCharacterTurn(p, m)` | Attaquer / Sorts / Inventaire |
| `basicAttack(p, m)` | attaque de base, avec l'élément de l'espèce |
| `spellMenu`, `spellMenuText` | choix du sort, sorts impossibles en rouge |
| `combatInventoryMenu(p, m)` | potions en combat, poison lancé sur l'ennemi |
| `runMonsterTurn(p, m, turn)` | choisit le pattern du bon ennemi |
| `verPattern`, `sanglierPattern`, `loupGarouPattern`, `dragonPattern` | comportement de chaque ennemi |
| `monsterBasicAttack(p, m, pct)` | attaque à X % (format imposé, tâche 20) |
| `monsterDamage(m, pct)` | dégâts ennemis avec buffs et bonus de phase |
| `dealDamage(...)` | dégâts joueur → ennemi |
| `dealDamageToPlayerElement(...)` | dégâts ennemi → joueur |
| `applyEffect(...)` | ajoute un effet (un buff est rafraîchi, un DoT se cumule) |
| `effectiveInitiative`, `effectiveAttackBonus` | stats avec les effets actifs |
| `resolveEffects`, `tickEffects` | effets de fin de tour |
| `announcePhaseChange(m)` | message + phrase d'histoire au changement de phase |
| `endFight(p, m, won)` | or, XP, loot, évolution, soin complet |
| `printSequential`, `autoAdvance`, `pauseForInput` | rythme de l'affichage |

### combat_screen.go — Quentin
| Fonction | Rôle |
|---|---|
| `combatScreen(p, m, turn, playerFirst)` | assemble tout l'écran de combat |
| `enemyIntentText`, `predictMonsterSkillName` | prochaine action de l'ennemi |
| `resistanceLines`, `nextPhaseHint`, `phaseThresholds` | infos tactiques |
| `healthBar(...)` | **partagée** — barre `[████░░░░] 52/80` |
| `boxify`, `dividedColumns`, `centerText` | **partagée** — cadre et colonnes |
| `visibleWidth`, `padRight`, `truncateVisible` | largeur réelle à l'écran (accents, emojis, couleurs) |
| `clearScreen()` | **partagée** — efface le terminal |
| `slotName(item)` | nom de l'artefact équipé, ou `—` |

### monsters.go — Quentin
| Fonction | Rôle |
|---|---|
| `initVer`, `initSanglier`, `initLoupGarou`, `initDragonAncien` | les 4 ennemis |
| `monsterByTier(tier)` | nouvel ennemi selon le tier |
| `updatePhase(m) string` | fait passer le Dragon à la phase suivante |
| `rollLoot(m) []Item` | tire chaque ressource de la table de loot |

### spells.go — Quentin
| Fonction | Rôle |
|---|---|
| `allSpells()` | tous les sorts |
| `spellsForStage(s, stage)` | sorts débloqués à ce palier précis |
| `spellByName`, `hasSpell` | recherche |
| `learnSpell(p, name) bool` | **partagée** — apprend un sort |
| `canCastSpell(p, spell)` | vérifie le mana et l'ultime, avec un message |
| `castSpell(p, m, spell)` | applique tous les effets du sort |

### evolution.go — Quentin
| Fonction | Rôle |
|---|---|
| `allSpecies()`, `speciesByName(name)` | les 3 espèces |
| `canEvolve(p, tier) bool` | première victoire et pas encore au palier 4 |
| `evolve(p)` | palier +1, stats, soin, nouveaux sorts |
| `currentStageName(p)`, `onEvolve(...)` | nom du palier et message « X devient Y ! » |
| `getResistance(s, e)` | résistance d'une espèce |
| `isGameComplete(p)` | le Dragon est-il vaincu ? |

### imageart.go / images_map.go — Quentin
| Fonction | Rôle |
|---|---|
| `displayImage(path, width) bool` | affiche un PNG dans le terminal |
| `renderImageString(path, width)` | même chose sous forme de string, pour la placer dans un écran |
| `speciesImagePath(name, stage)` | image d'une espèce à un palier |
| `monsterImagePath(name)` | image d'un ennemi |

### story.go, audio.go, banners.go, title.go — commun
| Fonction | Rôle |
|---|---|
| `storyPrologue`, `storyPreFight`, `storyPostFight`, `storyEvolve` | scènes d'histoire |
| `storyFinale`, `storyChoiceScreen` | choix Épargner / Achever |
| `storyEndingClemence`, `storyEndingRancune`, `storyCredits` | les deux fins + générique |
| `storyJournal()` | Mémoires de la vallée |
| `playMusic(path)`, `stopMusic()` | musique en boucle |
| `displayBanner(path)` | bannière d'un menu |
| `displayTitleScreen()` | écran titre |