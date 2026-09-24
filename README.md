<p align="center">
  <img src="src/assets/ui/title.png" alt="EvoBeasts" width="720">
</p>

<p align="center">
  <b>Un RPG au tour par tour dans le terminal, écrit en Go.</b><br>
  Faites éclore un œuf, faites-le évoluer, et affrontez le Dragon Ancien.
</p>

---

## L'histoire

<img src="src/assets/story/vallee_brumes.png" alt="La Vallée des Brumes" width="360" align="right">

La **Vallée des Brumes** était protégée par trois familles : les **Phénix**,
les **Krakens** et les **Caméléons Spectraux**.

Une nuit, le **Dragon Ancien** s'est réveillé. Il leur a volé leurs éléments
et a changé leurs gardiens en pierre. Il ne reste que trois œufs.

Vous incarnez l'un d'eux. Battez les bêtes de la vallée, évoluez, descendez
dans la tanière du Dragon… et décidez de son sort.

<br clear="right">

---

## Lancer le jeu

**Prérequis :** Go 1.27 ou plus, et un terminal qui gère la couleur 24 bits
(Windows Terminal, le terminal de VS Code). Au premier build, il faut
Internet pour télécharger les dépendances.

**Windows**
```bash
cd src
go run .
```
Vous pouvez aussi double-cliquer sur `evobeasts.exe` dans le dossier `src`.

**Linux / macOS**
```bash
cd src
go build -o evobeasts
./evobeasts
```

> ⚠️ L'exécutable doit être lancé **depuis le dossier `src/`**. Les images et
> les musiques sont chargées depuis `src/assets/`. Depuis un autre dossier, le
> jeu fonctionne, mais sans images ni son.

### Commandes

| Touche | Action |
|---|---|
| ↑ ↓ / ← → | se déplacer dans les menus |
| Entrée | valider |
| Échap | revenir en arrière |
| 1 – 9 | raccourci vers une option |
| Ctrl + C | quitter |

---

## Les trois espèces

Chaque œuf évolue trois fois, à la première victoire contre chaque ennemi.

| | Palier 1 | Palier 2 | Palier 3 | Palier 4 |
|---|:---:|:---:|:---:|:---:|
| **Phénix**<br>Feu · DPS | <img src="src/assets/creatures/phenix_1.png" width="110"><br>Œuf ardent | <img src="src/assets/creatures/phenix_2.png" width="110"><br>Poussin | <img src="src/assets/creatures/phenix_3.png" width="110"><br>Poulet | <img src="src/assets/creatures/phenix_4.png" width="110"><br>Phénix |
| **Kraken**<br>Eau · Tank | <img src="src/assets/creatures/kraken_1.png" width="110"><br>Œuf immergé | <img src="src/assets/creatures/kraken_2.png" width="110"><br>Alevin | <img src="src/assets/creatures/kraken_3.png" width="110"><br>Poisson | <img src="src/assets/creatures/kraken_4.png" width="110"><br>Kraken |
| **Caméléon Spectral**<br>Spectre · Assassin | <img src="src/assets/creatures/spectre_1.png" width="110"><br>Œuf voilé | <img src="src/assets/creatures/spectre_2.png" width="110"><br>Caméléonneau | <img src="src/assets/creatures/spectre_3.png" width="110"><br>Caméléon | <img src="src/assets/creatures/spectre_4.png" width="110"><br>Caméléon Spectral |

- **Phénix** : brûlures, soins, et **Renaissance**, qui le ramène à la vie
  une fois par combat.
- **Kraken** : bouclier, réduction de dégâts, ralentit l'ennemi et lui fait
  perdre des tours.
- **Caméléon Spectral** : devient invisible pour frapper beaucoup plus fort
  à l'attaque suivante.

**Le triangle des éléments :** l'Eau bat le Feu, le Feu bat le Spectre et le
Spectre bat l'Eau.

---

## Les ennemis

| <img src="src/assets/creatures/ver.png" width="140"> | <img src="src/assets/creatures/sanglier.png" width="140"> | <img src="src/assets/creatures/loup_garou.png" width="140"> | <img src="src/assets/creatures/dragon_ancien.png" width="140"> |
|:---:|:---:|:---:|:---:|
| **Ver**<br>Novice | **Sanglier**<br>Intermédiaire | **Loup-garou**<br>Avancé | **Dragon Ancien**<br>Boss |
| Le tutoriel | Charge, étourdit, s'enrage quand il est blessé | Double griffure, hurlement, se régénère | 4 phases : Feu, Eau, Spectre, puis les trois à la fois |

Chaque victoire débloque l'ennemi suivant. Les ennemis déjà battus restent
disponibles pour gagner de l'XP et des ressources.

---

## Fonctionnalités

- ⚔️ **Combat tactique** : initiative, résistances, effets sur la durée,
  bouclier, esquive. L'écran de combat montre **l'action que l'ennemi va
  faire** et ses faiblesses du moment.
- 🐉 **Boss en quatre phases** : le Dragon change d'élément et de
  résistances selon ses PV restants.
- 🧬 **Deux progressions** : l'**évolution** fait avancer l'histoire, et les
  **niveaux** permettent de farmer si on bloque.
- 🛒 **Marchand** : potions de soin, de mana et de poison, ressources,
  grimoire, et sac agrandissable jusqu'à 40 places.
- 🔨 **Forgeron** : trois artefacts à fabriquer et à équiper (tête, torse,
  pieds).
- 📖 **Histoire complète** : prologue, scènes avant et après chaque combat,
  et **deux fins** selon votre choix final. On peut relire les scènes dans
  les *Mémoires de la vallée*.
- 🖼️ **Images dans le terminal** : les créatures et les scènes sont des PNG
  affichés en blocs colorés.
- 🎵 **Musique** : ambiance, combat, boss et générique de fin.

---

## Structure du projet

```
src/
├── main.go, menu.go, types.go       base commune
├── character.go, inventory.go       personnage et inventaire
├── merchant.go, economy.go          marchand, or, XP
├── craft.go, equipment.go           forgeron et équipement
├── combat.go, combat_screen.go      combat et affichage
├── monsters.go, spells.go           ennemis et sorts
├── evolution.go                     espèces et évolutions
├── story.go                         histoire et fins
├── imageart.go, images_map.go       rendu des images
├── audio.go, banners.go, title.go   musique, bannières, écran titre
└── assets/                          images et musiques
```

L'organisation détaillée, les règles d'équipe et toutes les stats sont dans
[docs/gestion-projet.md](docs/gestion-projet.md).

---

## L'équipe

| Membre | Partie |
|---|---|
| **Lucas B** | personnage, inventaire |
| **Kelissa** | marchand, économie, niveaux |
| **Lucas H** | forgeron, équipement |
| **Quentin** | combat, ennemis, sorts, évolution, affichage |

<p align="center"><i>Projet RED, Ynov</i></p>