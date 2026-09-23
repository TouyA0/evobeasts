package main

// IDENTIFIANTS PARTAGES

// Ressources de craft
const (
	ResEcaille = "Fragment d'écaille"
	ResAbyssal = "Éclat abyssal"
	ResSpectre = "Poussière spectrale"
	ResPlume   = "Plume calcinée"
)
 
// Consommables
const (
	PotionVie     = "Essence vitale"
	PotionPoison  = "Essence corrompue"
	PotionMana    = "Essence de mana"
	GrimoireSoufle = "Grimoire : Souffle primordial"
	ExpansionSac  = "Expansion de besace"
)
 
// Artefacts
const (
	ArtefactTete  = "Diadème des Origines"
	ArtefactTorse = "Carapace primordiale"
	ArtefactPieds = "Serres du Voile"
)
 
// Noms des espèces
const (
	EspecePhenix  = "Phénix"
	EspeceKraken  = "Kraken"
	EspeceSpectre = "Caméléon Spectral"
)

// Noms de sorts
const (
	SortSoufflePrimordial   = "Souffle primordial"
	SortRenaissance         = "Renaissance"
	SortExecutionSpectrale  = "Exécution spectrale"
)

// Noms des ennemis
const (
	MonstreVer          = "Ver"
	MonstreSanglier     = "Sanglier"
	MonstreLoupGarou    = "Loup-garou"
	MonstreDragonAncien = "Dragon Ancien"
)
 
// Constantes de règles
const (
	InventaireBase    = 10  // capacité de départ
	InventaireBonus   = 10  // gain par expansion
	MaxExpansions     = 3   // nombre maximum d'expansions
	OrDepart          = 100 // or au démarrage
	XPPremierNiveau   = 100 // XP requise pour le niveau 2
	XPIncrementNiveau = 50  // XP supplémentaire requise par niveau
	GainPVNiveau      = 5   // PV max gagnés par niveau
	GainManaNiveau    = 3   // mana max gagné par niveau
	GainAtkNiveau     = 1   // attaque gagnée par niveau
	CoutCraft         = 10   // or dépensé par fabrication
	PermanentEffectDuration = 9999 // pour un effet qui dure tout le combat
)

// Valeurs des consommables.
const (
	SoinPotionVie      = 30
	ManaPotionMana     = 25
	DegatsPotionPoison = 8
	TicksPotionPoison  = 4
)

// Prix marchand
const (
	PrixPotionVie    = 12
	PrixPotionMana   = 10
	PrixPotionPoison = 10
	PrixGrimoire     = 35
	PrixExpansion    = 30

	PrixPlume   = 2
	PrixEcaille = 5
	PrixSpectre = 4
	PrixAbyssal = 12
)

// Valeurs des ultimes
const (
	RenaissanceHPPct     = 40 // % des PV max rendus par Renaissance
	RenaissanceAtkBuff   = 8  // bonus d'attaque accordé par Renaissance
	RenaissanceBuffTurns = 4  // durée de ce bonus
	ExecutionFurtiveDmg  = 110 // dégâts de l'Exécution spectrale lancée depuis l'invisibilité
	InvisibleDodgeBonus  = 40  // % d'esquive en plus tant qu'on est invisible
)

// ÉLÉMENTS
type Element string
const (
	ElementFeu 		Element = "Feu"
	ElementEau 		Element = "Eau"
	ElementSpectre 	Element = "Spectre"
	ElementNeutre 	Element = "Neutre" //atk de base sans élément
)

// LIGNEES ET EVOLUTIONS
type Species struct {
	Name 		string
	Element 	Element
	StageNames 	[]string

	BaseMaxHP 	int
	BaseMaxMana int
	BaseAttack 	int

	HPPerStage 	int
	ManaPerStage   int
	AttackPerStage int

	Initiative int
	DodgeRate  int

	Resistances map[Element]float64
}

// SORTS
type Spell struct {
	Name        string
	Description string
	Element     Element
	ManaCost    int
	MinStage    int // palier requis (1 à 4)
 
	// Effets offensifs
	Damage       int  // dégâts directs
	DotDamage    int  // dégâts par tour
	DotDuration  int  // durée du DoT en tours
	IgnoreResist bool // ignore les résistances de la cible
 
	// Effets défensifs et utilitaires
	Heal              int  // soins immédiats
	ShieldValue       int  // points de bouclier accordés
	DamageReduction   int  // % de réduction des dégâts subis
	ReductionDuration int  // durée de la réduction en tours
	AttackBuff        int  // bonus de dégâts pendant BuffDuration tours
	BuffDuration      int  // durée du buff en tours
	Invisible         bool // rend invisible 1 tour
	InvisibleBonus    int  // % de dégâts bonus à la sortie d'invisibilité
	Intangible        bool // esquive tout mais ne peut pas agir 1 tour
	SkipEnemyTurn     bool // l'adversaire perd son prochain tour
	InitiativeDebuff  int  // réduction d'initiative infligée à la cible
	AttackDebuff      int  // réduction d'attaque infligée à la cible
	DebuffDuration    int  // durée du debuff en tours
	IsUltimate        bool // utilisable une seule fois par combat
}

// OBJETS ET INVENTAIRE
type ItemType string
const (
	ItemPotion    ItemType = "potion"    // consommable
	ItemResource  ItemType = "resource"  // matériau de craft
	ItemArtifact  ItemType = "artifact"  // équipable
	ItemSpellbook ItemType = "spellbook" // grimoire, enseigne un sort
	ItemUpgrade   ItemType = "upgrade"   // expansion d'inventaire
)

// emplacement equipement
type EquipSlot string
const (
	SlotNone  EquipSlot = ""
	SlotHead  EquipSlot = "tête"
	SlotChest EquipSlot = "torse"
	SlotFeet  EquipSlot = "pieds"
)

type Item struct {
	Name  string
	Type  ItemType
	Price int // prix chez le marchand, 0 si non vendu
	Description string
 
	// Potions
	HealValue    int // PV rendus
	ManaValue    int // mana rendu
	PoisonDamage int // dégâts par tick
	PoisonTicks  int // nombre de ticks
 
	// Artefacts
	Slot         EquipSlot
	BonusMaxHP   int
	BonusMaxMana int
 
	// Grimoire
	SpellName string // nom du sort enseigné
}

// ressource pouvant être droppée par un ennemi.
type LootEntry struct {
	Item   Item
	Chance int // probabilité de drop en pourcentage (0 à 100)
}

// Equipment regroupe les emplacements équipés du personnage.
type Equipment struct {
	Head  *Item
	Chest *Item
	Feet  *Item
}
 
// Recipe décrit une recette du forgeron.
type Recipe struct {
	Result    Item
	GoldCost  int
	Materials map[string]int
}

// EFFETS DE COMBAT
type StatusEffect struct {
	Name             string
	TurnsLeft        int
	DamagePerTurn    int  // brûlure, poison, abysse
	HealPerTurn      int  // régénération
	AttackBuff       int  // bonus de dégâts, s'applique aux sorts aussi
	AttackDebuff     int  // malus de dégâts
	DamageReduction  int  // % de réduction des dégâts subis
	InitiativeDebuff int  // réduction d'initiative
	SkipTurn         bool // fait perdre le tour
}

// PERSONNAGE
type Character struct {
	Name  string
	Level int
	Species Species
	Stage   int // palier actuel : 1 (œuf) à 4 (forme finale)
 
	MaxTierDefeated int // débloque adv suivants et déclenche évolution si victoire
 
	// Statistiques dérivées, recalculées par recalculateStats()
	MaxHP      int
	MaxMana    int
	Attack     int
	Initiative int
	DodgeRate  int
 
	// Statistiques courantes, modifiables librement en combat
	CurrentHP   int
	CurrentMana int
 
	// Progression par expérience
	Experience    int
	MaxExperience int
 
	// Possessions
	Gold          int
	Inventory     []Item
	InventorySize int // InventaireBase au départ
	UpgradesUsed  int // maximum MaxExpansions
	Spells        []Spell
	Equipment     Equipment
 
	// États de combat, réinitialisés au début de chaque combat
	Shield         int // absorbe les dégâts avant les PV
	Effects        []StatusEffect
	Invisible      bool // prochaine attaque boostée
	InvisibleBonus int  // % de bonus en attente
	Intangible     bool // esquive tout ce tour
	UltimateUsed   bool // ultime déjà utilisé ce combat
	SkipNextTurn   bool // perd son prochain tour
}

// COMPÉTENCES ENNEMIES
type SkillKind string

const (
	SkillAttack SkillKind = "attaque"
	SkillBuff   SkillKind = "buff"
	SkillDebuff SkillKind = "debuff"
	SkillHeal   SkillKind = "soin"
)

type MonsterSkill struct {
	Name        string
	Kind        SkillKind
	Element     Element
	DamagePct   int    // % de l'ATK de base, 0 si non-offensif
	Description string // condition de déclenchement + effet
}

// MONSTRES
type Phase struct {
	Name        string              // "Éveil Igné", "Vague Ancestrale"...
	HPThreshold float64             // % de PV sous lequel la phase s'active
	Element     Element
	Resistances map[Element]float64 // multiplicateur de dégâts reçus
	DamageBonus int                 // % de dégâts bonus (phase finale)
	Attacks     []MonsterSkill
}

type Monster struct {
	Name       string
	Tier       int // 1 (Ver) à 4 (Dragon Ancien)
	MaxHP      int
	CurrentHP  int
	Attack     int
	Initiative int
	Element    Element
	Skills []MonsterSkill
 
	// Résistances hors phases
	Resistances map[Element]float64
 
	// Récompenses
	XPReward   int
	GoldReward int
	LootTable  []LootEntry
 
	// Boss multi-phases (nil pour les ennemis simples)
	Phases       []Phase
	CurrentPhase int
 
	// États de combat
	Effects      []StatusEffect
	Shield       int
	Intangible   bool
	SkipNextTurn bool
	StunImmune   bool // vient d'être étourdi, ne peut pas l'être au tour suivant
 
	// Suivi du pattern d'attaque
	TurnCounter int  // tour courant, pour les patterns cycliques
	HowlUsed    bool // Hurlement du Loup-garou, une seule fois par combat
}
