# EvoBeasts

les branchs git :
1. celui qui taff sur les personnages : feature/character-base
2. celui qui taff sur les marchands : feature/merchant
3. celui qui taff sur les crafts : feature/crafting
4. celui qui taff sur les combats : feature/combat

## Fonctionnement git :
### Étape 1 — Avant de toucher au code, vérifier où on est
 
**Toujours commencer une session par `git status`.** Ça affiche sur quelle
branche vous êtes. Si ce n'est pas votre branche (`feature/mon-module`),
basculez dessus avant de modifier quoi que ce soit :
 
```bash
git status
git checkout feature/mon-module
```
 
**Ne codez jamais directement sur `main`.** Si `git status` vous dit que
vous êtes sur `main`, ne touchez à aucun fichier tant que vous n'avez pas
fait le `checkout` vers votre branche.
 
### Étape 2 — Récupérer les derniers changements avant de coder
 
Une fois sur sa branche, on récupère ce que les autres ont pu envoyer sur
`main` depuis la dernière fois :
 
```bash
git checkout main
git pull origin main
git checkout feature/mon-module
git merge main
```
 
**Faites ça à chaque fois que vous commencez à coder**, même si vous avez
l'impression que rien n'a changé. Ça évite de travailler sur une base
périmée et de découvrir un gros conflit trois jours plus tard.
 
### Étape 3 — Pendant qu'on code
 
On commit souvent, par petits blocs qui fonctionnent, pas un seul énorme
commit en fin de journée :
 
```bash
git add .
git commit -m "combat: ajoute le pattern du sanglier"
git push origin feature/mon-module
```
 
Le message de commit décrit ce qu'on vient de faire, en commençant par le
domaine (`combat:`, `merchant:`, `craft:`, `character:`).
 
**Repassez par l'étape 1 (`git status`) avant chaque nouvelle session**,
pour être sûr d'être sur la bonne branche avant de modifier quoi que ce
soit.
 
### Étape 4 — Une fois son bout terminé : envoyer vers main
#### en vrai le mieux serait d'attendre qu'on soit ensemble pour merge sur main comme ça on le fait ensemble proprement.
```bash
git checkout main
git pull origin main
git checkout feature/mon-module
git merge main
```
 
Si Git signale un conflit ici, on ne force rien : on regarde ensemble ce
qui a changé des deux côtés avant de choisir quoi garder.
 
Une fois le merge propre, on envoie sur `main` :
 
```bash
git checkout main
git merge feature/mon-module
git push origin main
```
 
**Juste après le `git push origin main`, on repart sur sa propre
branche**, pour ne pas continuer à coder par erreur directement sur
`main` :
 
```bash
git checkout feature/mon-module
```

---

## 1. L'arborescence

```
src/
├── go.mod
├── main.go          COMMUN  → menu principal uniquement
├── types.go         COMMUN  → toutes les structs et constantes
│
├── character.go     Membre 2 → création, affichage, stats
├── inventory.go     Membre 2 → inventaire + fonctions partagées
│
├── merchant.go      Membre 3 → boutique
├── economy.go       Membre 3 → or, XP, niveaux
│
├── craft.go         Membre 4 → forgeron
├── equipment.go     Membre 4 → artefacts équipés
│
├── combat.go        Quentin → boucle de combat, tours
├── monsters.go      Quentin → les 4 ennemis + phases du boss
├── spells.go        Quentin → sorts et effets
├── evolution.go     Quentin → espèces, paliers, évolutions
└── display.go       Quentin → ASCII art des créatures
```

Un fichier par domaine, comme ça chacun bosse dans les siens et on évite
les conflits Git. Les deux seuls fichiers communs sont `types.go` (posé
une fois au départ) et `main.go`, où chacun n'ajoute que son `case` dans
le switch du menu.

Le détail des fonctions de chaque fichier est en fin de doc (section 12).

---

## 2. Qui fait quoi

**Membre 2 — le personnage et la base**
Tâches 1 à 6, 8, 11, 12. Il pose la structure du jeu : création du perso,
affichage des infos, inventaire, menu principal, mort et résurrection,
limite d'inventaire. Il écrit aussi les fonctions que tout le monde
utilise.
**Il est sur le chemin critique** : tant qu'il n'a pas poussé ses
fonctions, les trois autres ne peuvent pas compiler. Même vides, il faut
qu'elles existent dès le premier jour.

**Membre 3 — l'économie**
Tâches 7, 9, 13, 14, 18 et Mission 2. Le marchand, les prix, l'or, la
potion de poison, l'expansion d'inventaire, et tout le système
d'expérience et de montée en niveau.

**Membre 4 — le craft**
Tâches 15, 16, 17 et Mission 6. Le forgeron, les recettes, la structure
Equipment, équiper et déséquiper les artefacts avec les bonus de stats.

**Quentin — le combat et l'univers**
Tâches 10, 19 à 22, Missions 1, 3 et 4. Toute la partie combat, les
sorts, les quatre ennemis avec leurs patterns, le boss à quatre phases,
le système d'évolution et l'affichage ASCII des créatures.

---

## 3. Les règles à respecter absolument

Tout est dans `package main`, donc **un seul espace de noms pour tout le
monde**. Si deux d'entre nous écrivent une fonction `displayMenu`, ça ne
compile plus. Chacun préfixe donc ses fonctions internes par son domaine :
`merchantDisplay`, `craftDisplay`, `combatDisplay`.

### Les statistiques dérivées

`MaxHP`, `MaxMana`, `Attack` et `Initiative` sont modifiés par **trois
sources différentes** : l'évolution (moi), la montée de niveau (membre 3)
et les artefacts équipés (membre 4).

**Personne n'écrit `player.MaxHP += X`.** On modifie sa source (`Stage`,
`Level`, `Equipment`), puis on appelle `recalculateStats()` qui recalcule
tout depuis zéro :

```
MaxHP = BaseMaxHP
      + HPPerStage × (Stage − 1)
      + GainPVNiveau × (Level − 1)
      + somme des BonusMaxHP des artefacts équipés
```

Pareil pour le mana et l'attaque. Sans ça, déséquiper un artefact devient
impossible à calculer correctement et les stats vont dériver.

### Les noms d'objets sont des constantes

Tous les noms sont définis en haut de `types.go` (`ResEcaille`,
`PotionVie`, `ArtefactTete`...). **Ne comparez jamais avec une chaîne
écrite en dur** : une apostrophe typographique au lieu d'une droite et le
craft échoue sans message d'erreur. Pour renommer un objet, on modifie la
constante, pas les usages.

---

## 4. Le jeu en deux phrases

On incarne un œuf qui évolue en créature légendaire en battant quatre
ennemis de plus en plus durs, jusqu'à l'ancêtre commun des trois espèces.
Combat tour par tour, 1 contre 1, en ligne de commande.

**Contrainte qu'on s'impose : la démo doit tenir en quelques minutes
devant la classe.** Tous les chiffres sont calibrés pour des combats de
4 à 6 tours et une partie complète en une dizaine de minutes. Si un
combat traîne pendant l'oral, on perd le jury.

---

## 5. Les deux progressions

On a volontairement séparé deux choses souvent mélangées dans les RPG.

### L'évolution — c'est la trame

Le palier monte à la **première victoire** contre chaque ennemi.

| Ennemi battu | Palier atteint |
|---|---|
| Ver | 2 — Poussin / Alevin / Caméléonneau |
| Sanglier | 3 — Poulet / Poisson / Caméléon |
| Loup-garou | 4 — Phénix / Kraken / Caméléon Spectral |
| Dragon Ancien | fin du jeu |

On démarre au palier 1, sous forme d'œuf, et le premier combat c'est
l'œuf contre le Ver — deux trucs minuscules qui se battent, ça se tient
narrativement.

Rejouer un ennemi déjà battu ne fait pas réévoluer.

### Le niveau — c'est le filet de sécurité

L'XP tombe à chaque combat, y compris les combats répétés. Ça permet de
farmer si on bloque, sans casser la trame des évolutions. C'est aussi la
Mission 2, donc des points bonus.

- 100 XP pour le niveau 2, puis +50 par niveau (150, 200, 250...)
- Par niveau : +5 PV max, +3 mana max, +1 attaque
- L'excédent d'XP est reporté sur le niveau suivant

En battant les quatre ennemis une fois, on arrive vers le niveau 4.

---

## 6. Les trois espèces

### Stats de base (palier 1) et progression par évolution

| Espèce | Élément | Rôle | PV | PV/palier | Mana | Mana/palier | Atk | Atk/palier | Init | Esquive |
|---|---|---|---|---|---|---|---|---|---|---|
| Phénix | Feu | DPS | 50 | +18 | 40 | +15 | 5 | +3 | 12 | 5 % |
| Kraken | Eau | Tank | 75 | +25 | 30 | +12 | 5 | +2 | 6 | 0 % |
| Caméléon Spectral | Spectre | Assassin | 45 | +14 | 35 | +15 | 5 | +4 | 18 | 15 % |

Le rôle est juste indicatif, il n'intervient dans aucun calcul.

Ce que ça donne au palier 4, sans compter les niveaux ni les artefacts :

| Espèce | PV max | Mana max | Attaque |
|---|---|---|---|
| Phénix | 104 | 85 | 14 |
| Kraken | 150 | 66 | 11 |
| Caméléon Spectral | 87 | 80 | 17 |

### Résistances

Multiplicateur de dégâts reçus. Sous 1.0 = on encaisse moins.

| Espèce | Feu | Eau | Spectre | Neutre |
|---|---|---|---|---|
| Phénix | 0.6 | 1.4 | 1.0 | 1.0 |
| Kraken | 0.7 | 0.6 | 1.2 | 0.9 |
| Caméléon Spectral | 1.2 | 1.0 | 0.5 | 1.0 |

### Les noms des paliers

| Espèce | Palier 1 | Palier 2 | Palier 3 | Palier 4 |
|---|---|---|---|---|
| Phénix | Œuf ardent | Poussin | Poulet | Phénix |
| Kraken | Œuf immergé | Alevin | Poisson | Kraken |
| Caméléon Spectral | Œuf voilé | Caméléonneau | Caméléon | Caméléon Spectral |

Au démarrage : PV à 50 % du max (imposé par le sujet, tâche 11), mana
plein, 100 or, et deux Essences vitales dans le sac pour que le premier
combat soit jouable.

---

## 7. Les compétences

Chaque sort s'apprend automatiquement à l'évolution correspondante. Au
palier 1 on n'a que l'attaque de base, c'est normal.

**Les buffs d'attaque s'appliquent aussi aux dégâts des sorts**, pas
seulement à l'attaque de base.

### Phénix — Feu

| Sort | Palier | Mana | Effet |
|---|---|---|---|
| Coup de bec | 1 | 0 | attaque de base |
| Picorage enflammé | 2 | 8 | 10 dégâts + brûlure 3/tour pendant 2 tours |
| Plumage chaud | 2 | 6 | rend 12 PV |
| Coup de bec ardent | 3 | 12 | 22 dégâts |
| Cocorico enflammé | 3 | 10 | +8 dégâts pendant 3 tours |
| Griffes de braise | 4 | 15 | 32 dégâts |
| Nova solaire | 4 | 20 | 18 dégâts + brûlure 6/tour pendant 4 tours |
| Renaissance (ultime) | 4 | 0 | à la mort : revient à 50 % PV + 10 dégâts bonus, une fois par combat |

### Kraken — Eau

| Sort | Palier | Mana | Effet |
|---|---|---|---|
| Coup de nageoire | 1 | 0 | attaque de base |
| Éclaboussure | 2 | 6 | 8 dégâts + initiative adverse −5 pendant 2 tours |
| Écailles naissantes | 2 | 8 | −20 % de dégâts subis pendant 3 tours |
| Coup de nageoire lourd | 3 | 10 | 16 dégâts + attaque adverse −3 pendant 1 tour |
| Bulle protectrice | 3 | 12 | bouclier de 25 points |
| Étreinte tentaculaire | 4 | 18 | 20 dégâts + l'adversaire perd son prochain tour |
| Raz-de-marée | 4 | 16 | 26 dégâts + initiative adverse −8 pendant 3 tours |
| Abysse (ultime) | 4 | 25 | 10 dégâts/tour pendant 5 tours |

### Caméléon Spectral — Spectre

| Sort | Palier | Mana | Effet |
|---|---|---|---|
| Morsure | 1 | 0 | attaque de base |
| Morsure furtive | 2 | 7 | 14 dégâts, doublés si invisible |
| Camouflage | 2 | 8 | invisible 1 tour, +50 % dégâts à la prochaine attaque |
| Frappe chromatique | 3 | 12 | 24 dégâts |
| Voile chromatique | 3 | 10 | invisible 1 tour, +100 % dégâts à la prochaine attaque |
| Lame spectrale | 4 | 15 | 28 dégâts, ignore les résistances |
| Dissolution | 4 | 12 | intangible 1 tour : on esquive tout mais on ne peut pas agir |
| Exécution spectrale (ultime) | 4 | 30 | 60 dégâts depuis l'invisibilité, 25 sinon |

### Le sort acheté

La tâche 10 impose un livre de sort acheté au marchand, donc on en garde
un en plus des sorts d'évolution :

**Souffle primordial** — grimoire à 25 or, coûte 14 mana, 20 dégâts,
élément neutre. Achetable dès le début, apprenable une seule fois.

---

## 8. Les quatre ennemis

| Ennemi | Palier | PV | Atk | Init | Élément | XP | Or | Tours estimés |
|---|---|---|---|---|---|---|---|---|
| Ver | 1 | 20 | 3 | 3 | Neutre | 40 | 5 | 4 |
| Sanglier | 2 | 45 | 7 | 9 | Neutre | 90 | 12 | 4-5 |
| Loup-garou | 3 | 80 | 12 | 13 | Neutre | 150 | 25 | 5-6 |
| Dragon Ancien | 4 | 150 | 16 | 15 | les trois | 300 | 80 | 8-10 |

Les trois premiers n'ont aucune résistance particulière (1.0 partout),
seul le Dragon joue là-dessus.

### Ce qu'ils droppent

En plus d'être achetables au marchand, certaines ressources tombent en
fin de combat :

| Ennemi | Ressource | Chance |
|---|---|---|
| Ver | Plume calcinée | 60 % |
| Sanglier | Fragment d'écaille | 60 % |
| Loup-garou | Poussière spectrale | 50 % |
| Loup-garou | Fragment d'écaille | 40 % |
| Dragon Ancien | Éclat abyssal | 100 % |

Si le sac est plein, le loot est perdu avec un message clair.

### Leurs patterns

**Ver** — Morsure, 100 % de son attaque, chaque tour. C'est tout, c'est
le tutoriel.

**Sanglier**
- Charge un tour sur 3 : 130 % de dégâts, 30 % de chance d'étourdir 1 tour
- Enragé (passif) : +30 % de dégâts sous 50 % PV, +60 % sous 25 %
- Les autres tours : attaque simple

**Loup-garou**
- Griffure enragée un tour sur 2 : deux attaques à 70 % chacune
- Hurlement au tour 3, une seule fois : +5 d'attaque pour le reste du combat
- Régénération lunaire : +6 PV en fin de chaque tour

**Dragon Ancien — le boss, quatre phases**

C'est l'ancêtre commun des trois espèces, né de la fusion corrompue du
feu, de l'eau et du spectre. Chaque phase exprime un élément, et ses
résistances changent avec.

| Phase | Déclenche à | Élément | Résiste (Feu/Eau/Spectre) | Ce qu'il fait |
|---|---|---|---|---|
| 1 — Éveil Igné | 100 → 75 % PV | Feu | 0.5 / 1.5 / 1.0 | Souffle ardent (120 % + brûlure), Griffe (100 %) |
| 2 — Vague Ancestrale | 75 → 50 % PV | Eau | 1.5 / 0.5 / 1.0 | Raz-de-marée (120 %), Étreinte liquide (80 % + perte de tour, 1 tour sur 3), régénère 8 PV/tour |
| 3 — Voile Spectral | 50 → 25 % PV | Spectre | 1.0 / 1.5 / 0.5 | Frappe hors du plan (130 %, ignore le bouclier), Dissolution (esquive totale 1 tour sur 4) |
| 4 — Convergence | 25 → 0 % PV | les trois | 1.0 / 1.0 / 1.0 | rotation des trois attaques, +50 % de dégâts |

L'idée : chaque espèce a une phase où elle cartonne et une où elle
galère. Ça force à s'adapter au lieu de spammer le même sort.

### Le menu de combat

Le sujet demande juste un combat d'entraînement, mais on fait mieux : un
menu qui liste les adversaires, débloqués au fur et à mesure. Au départ
seul le Ver est dispo, chaque victoire ouvre le suivant.

Un ennemi déjà battu reste rejouable pour farmer.

---

## 9. Le marchand et le forgeron

### Ce que vend le marchand

| Objet | Type | Prix |
|---|---|---|
| Essence vitale | potion, +50 PV | 3 or |
| Essence corrompue | potion de poison, 10 dégâts × 3 ticks | 6 or |
| Essence de mana | potion, +30 mana | 5 or |
| Grimoire : Souffle primordial | livre de sort | 25 or |
| Fragment d'écaille | ressource | 4 or |
| Éclat abyssal | ressource | 7 or |
| Poussière spectrale | ressource | 3 or |
| Plume calcinée | ressource | 1 or |
| Expansion de besace | +10 places, 3 fois max | 30 or |

### Le cas de l'Essence corrompue

La tâche 9 impose une potion qui fait des dégâts avec une seconde entre
chaque tick (`time.Sleep`). On la rend utile en la faisant dépendre du
contexte :

- **En combat** : on la balance sur l'ennemi, 10 dégâts × 3 ticks
- **Hors combat** : on se l'inflige à soi-même, avec un message
  d'avertissement et une confirmation avant

Comme ça on respecte la consigne à la lettre tout en ayant un objet qui
sert vraiment.

### Les artefacts du forgeron

Trois emplacements : tête, torse, pieds. 5 or par fabrication, en plus
des ressources.

| Artefact | Slot | Ressources | Bonus |
|---|---|---|---|
| Diadème des Origines | tête | 1 Plume calcinée + 1 Fragment d'écaille | +10 PV max, +10 mana max |
| Carapace primordiale | torse | 2 Éclats abyssaux + 1 Fragment d'écaille | +25 PV max |
| Serres du Voile | pieds | 1 Poussière spectrale + 1 Fragment d'écaille | +15 PV max, +5 mana max |

Équiper retire l'objet de l'inventaire. Si le slot est déjà occupé,
l'ancien revient dans le sac — donc il faut vérifier qu'il reste une
place libre avant de faire l'échange.

Les noms sont provisoires, membres 3 et 4 peuvent les changer, mais en
modifiant les constantes dans `types.go`.

---

## 10. Les règles de combat et d'inventaire

### Combat

- Celui qui a la plus haute initiative commence le tour
- Dégâts = (base + buffs − malus) × résistance de la cible à l'élément
- L'esquive se teste avant chaque attaque reçue
- Le bouclier absorbe avant les PV
- Intangible = on esquive tout mais on ne peut pas agir
- Les effets se résolvent en fin de tour et perdent un tour de durée
- **PV et mana remis à fond après chaque combat gagné.** On ne veut pas
  d'une gestion de ressources entre les combats, ça rallongerait la démo
  pour rien
- Le mana ne remonte pas pendant un combat, seulement avec les Essences
  de mana
- À 0 PV : le monstre meurt, ou nous on meurt (tâche 8 : on revient à
  50 % des PV max, hors combat)

### Inventaire

- Les doublons prennent des places distinctes : trois Essences vitales =
  trois places sur dix
- 10 places au départ, extensible trois fois de +10, donc 40 max
- Toute fonction qui ajoute un objet vérifie la place et affiche une
  erreur explicite si c'est plein

---

## 11. Mission 6 — les deux artistes cachés

Les titres des tâches sont des références.

- **Partie 2** : *Money, money, money*, *Two for the Price of One*,
  *Gimme! Gimme! Gimme!*, *I saw it in the mirror*, *Mamma Mia*,
  *On and on and on* → des chansons d'**ABBA**
- **Partie 3** : *Duel*, *A.I. Intelligence artificielle*,
  *Ready Player One*, *Fighter Squad* → des films de **Steven Spielberg**

À revérifier avant de le mettre dans le jeu, mais ça colle bien.

---

## 12. Les fonctions par fichier

Aperçu de ce que chacun a à écrire. Les fonctions marquées **partagée**
sont appelées par d'autres modules : leur nom et leur signature ne
changent pas sans prévenir tout le monde.

### main.go — commun

| Fonction | Ce qu'elle fait |
|---|---|
| `main()` | lance la création du perso puis le menu |
| `mainMenu(p *Character)` | boucle du menu principal, un `case` par module |

### types.go — commun

Aucune fonction, uniquement les structs (`Character`, `Monster`, `Item`,
`Spell`, `Species`, `Equipment`, `Recipe`, `StatusEffect`, `Phase`,
`LootEntry`) et les constantes de noms et de règles.

### character.go — Membre 2

| Fonction | Ce qu'elle fait |
|---|---|
| `characterCreation() *Character` | demande le nom et l'espèce, renvoie le perso prêt |
| `initCharacter(...) *Character` | construit un Character à partir des valeurs de base |
| `formatName(s string) string` | majuscule au début, reste en minuscules |
| `displayInfo(p *Character)` | fiche du perso, appelle `displayCreature` pour l'ASCII |
| `recalculateStats(p *Character)` | **partagée** — recalcule MaxHP, MaxMana, Attack depuis zéro |
| `isDead(p *Character) bool` | **partagée** — vérifie 0 PV, ressuscite à 50 % |
| `readInput() string` | **partagée** — lecture clavier nettoyée |

### inventory.go — Membre 2

| Fonction | Ce qu'elle fait |
|---|---|
| `accessInventory(p *Character)` | affiche l'inventaire, permet d'utiliser un objet |
| `addInventory(p *Character, item Item) bool` | **partagée** — ajoute si place, sinon `false` |
| `removeInventory(p *Character, name string) bool` | **partagée** — retire une occurrence |
| `countItem(p *Character, name string) int` | **partagée** — compte les occurrences |
| `hasSpace(p *Character) bool` | **partagée** — reste-t-il une place ? |
| `useItem(p *Character, item Item)` | applique l'effet d'un objet hors combat |
| `takePot(p *Character)` | potion de vie, +50 PV plafonnés au max |

### merchant.go — Membre 3

| Fonction | Ce qu'elle fait |
|---|---|
| `merchantMenu(p *Character)` | boucle du menu boutique |
| `merchantCatalog() []Item` | la liste de tout ce qui est vendu |
| `buyItem(p *Character, item Item)` | vérifie l'or et la place, débite, ajoute |
| `upgradeInventorySlot(p *Character)` | +10 places, 3 fois maximum |
| `merchantDisplay(items []Item)` | affiche le catalogue avec les prix |

### economy.go — Membre 3

| Fonction | Ce qu'elle fait |
|---|---|
| `gainXP(p *Character, amount int)` | **partagée** — ajoute l'XP, gère les montées de niveau |
| `gainGold(p *Character, amount int)` | **partagée** — ajoute l'or |
| `levelUp(p *Character)` | +5 PV, +3 mana, +1 atk, reporte l'excédent d'XP |
| `xpForNextLevel(level int) int` | 100 puis +50 par niveau |

### craft.go — Membre 4

| Fonction | Ce qu'elle fait |
|---|---|
| `craftMenu(p *Character)` | boucle du menu forgeron |
| `recipes() []Recipe` | les trois recettes d'artefacts |
| `canCraft(p *Character, r Recipe) bool` | vérifie l'or, les ressources et la place |
| `craftItem(p *Character, r Recipe)` | consomme les matériaux, ajoute l'artefact |
| `craftDisplay(recipes []Recipe)` | affiche les recettes et ce qu'il manque |
| `showArtists()` | Mission 6, affiche ABBA et Spielberg |

### equipment.go — Membre 4

| Fonction | Ce qu'elle fait |
|---|---|
| `equipItem(p *Character, item Item) bool` | équipe, remet l'ancien au sac, appelle `recalculateStats` |
| `unequipItem(p *Character, slot EquipSlot) bool` | déséquipe si place dispo |
| `equipmentBonus(e Equipment) (hp int, mana int)` | somme des bonus, utilisée par `recalculateStats` |
| `equipmentDisplay(e Equipment)` | affiche les trois emplacements |

### combat.go — Quentin

| Fonction | Ce qu'elle fait |
|---|---|
| `combatMenu(p *Character)` | liste les ennemis débloqués, lance le combat choisi |
| `trainingFight(p *Character, m *Monster)` | boucle du combat, alterne les tours |
| `characterTurn(p *Character, m *Monster)` | menu du joueur : attaquer, sort, inventaire |
| `monsterTurn(p *Character, m *Monster)` | applique le pattern de l'ennemi |
| `dealDamage(base int, element Element, target ...)` | calcule les dégâts avec résistances et bouclier |
| `resolveEffects(...)` | applique les DoT, régén, buffs en fin de tour |
| `endFight(p *Character, m *Monster, won bool)` | XP, or, loot, évolution, restaure PV et mana |

### monsters.go — Quentin

| Fonction | Ce qu'elle fait |
|---|---|
| `initVer() *Monster` | et les trois autres `initXxx` |
| `monsterByTier(tier int) *Monster` | renvoie l'ennemi correspondant au palier |
| `updatePhase(m *Monster)` | fait basculer le Dragon de phase selon ses PV |
| `rollLoot(m *Monster) []Item` | tire les ressources droppées |

### spells.go — Quentin

| Fonction | Ce qu'elle fait |
|---|---|
| `allSpells() []Spell` | catalogue de tous les sorts du jeu |
| `spellsForStage(s Species, stage int) []Spell` | les sorts débloqués à ce palier |
| `learnSpell(p *Character, name string) bool` | **partagée** — ajoute un sort, `false` si déjà connu |
| `castSpell(p *Character, m *Monster, spell Spell)` | applique le sort, débite le mana |

### evolution.go — Quentin

| Fonction | Ce qu'elle fait |
|---|---|
| `allSpecies() []Species` | définition des trois espèces avec leurs stats |
| `speciesByName(name string) Species` | récupère une espèce par son nom |
| `evolve(p *Character)` | monte le palier, débloque les sorts, restaure PV et mana |
| `canEvolve(p *Character, defeatedTier int) bool` | première victoire contre cet ennemi ? |

### display.go — Quentin

| Fonction | Ce qu'elle fait |
|---|---|
| `displayCreature(s Species, stage int)` | **partagée** — ASCII art coloré de la créature |
| `displayMonster(m *Monster)` | ASCII art de l'ennemi |
| `healthBar(current, max int) string` | barre de vie visuelle |
| `combatDisplay(p *Character, m *Monster)` | écran de combat complet |