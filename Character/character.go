package character

import (
	"fmt"
	"slices"
	"strings"

	"john-v3/dndCharacterGenerator/DataGenerator"
)

type Character struct {
	Metadata  Metadata
	Character CharacterInfo
	Biography Biography
	Physical  Physical

	Attributes Attributes

	DerivedStrengthStats     DerivedStrengthStats
	DerivedDexterityStats    DerivedDexterityStats
	DerivedConstitutionStats DerivedConstitutionStats
	DerivedIntelligenceStats DerivedIntelligenceStats
	DerivedWisdomStats       DerivedWisdomStats
	DerivedCharismaStats     DerivedCharismaStats

	SavingThrows SavingThrows
	Combat       Combat
	Economy      Economy

	Spells                   []Spell
	WeaponSpecializations    []WeaponSpecialization
	NonWeaponSpecializations []NonWeaponSpecialization
	Languages                []string
	Gear                     []GearItem

	ThiefData  *ThiefData
	WizardData *WizardData
}

type Metadata struct {
	Seed    int64
	Ruleset string
	Version string
}

type CharacterInfo struct {
	Level       int
	Name        string
	Class       string
	ClassConfig DataGenerator.Class
	ClassHitDie int
	Race        string
	Alignment   string
}

type Biography struct {
	Family           string
	RaceClan         string
	Homeland         string
	Patron           string
	Religion         string
	SocialClass      string
	Status           string
	BirthRank        int
	NumberOfSiblings int
	Honor            int
	BaseHonor        int
}

type Physical struct {
	Hair       string
	Eyes       string
	Age        int
	Height     int
	Weight     int
	Appearance string
}

type Attributes struct {
	STR Strength
	DEX int
	CON int
	INT int
	WIS int
	CHA int
}

type Strength struct {
	Score       int
	Exceptional int
}

type DerivedStrengthStats struct {
	HitProbability    int
	DamageAdjustment  int
	WeightAllowance   int
	MaxPress          int
	OpenDoors         int
	BendBarsLiftGates int
}

type DerivedDexterityStats struct {
	ReactionAdjustment      int
	MissileAttackAdjustment int
	DefenseAdjustment       int
}

type DerivedConstitutionStats struct {
	Hitpoints    int
	HPAdjustment int
	SystemShock  int
	Resurrection int
	PoisonSave   int
	Regeneration int
}

type DerivedIntelligenceStats struct {
	NumberOfLanguages int
	SpellLevel        int
	LearnSpell        int
	SpellsPerLevel    int
	SpellImmunities   []string
}

type DerivedWisdomStats struct {
	MagicDefenseAdjustment int
	BonusSpells            []DataGenerator.BonusSpellEntry
	SpellFailure           int
	SpellImmunities        []string
}

type DerivedCharismaStats struct {
	MaxNumberOfHenchmen int
	LoyaltyBase         int
	ReactionAdjustment  int
}

type SavingThrows struct {
	ParalyzePoisonDeath int
	RodStaffOrWand      int
	PetrifyPolymorph    int
	BreathWeapon        int
	Spell               int
}

type Combat struct {
	ArmorClass    int
	THACO         int
	Surprised     int
	Rear          int
	AttackPenalty int
}

type Economy struct {
	Gold int
}

type Spell struct {
	Name        string
	Level       int
	School      string
	Description string

	Properties map[string]string
}

type WeaponSpecialization struct {
	Name string
	Rank string
}

type NonWeaponSpecialization struct {
	Name       string
	BaseStat   string
	TargetRoll int
}

type GearItem struct {
	Name     string
	Weight   int
	Quantity int

	Properties map[string]string
}

type ThiefData struct {
	PickPocket     int
	OpenLock       int
	FindRemoveTrap int
	MoveSilently   int
	HideInShadows  int
	DetectNoise    int
	ClimbWalls     int
	ReadLanguages  int
}

type WizardData struct {
	Spellbook       []Spell
	MemorizedSpells []string

	SpellSlots SpellSlots
}

type SpellSlots struct {
	Level1 int
	Level2 int
	Level3 int
	Level4 int
	Level5 int
	Level6 int
	Level7 int
	Level8 int
	Level9 int
}

type BonusSpell struct {
	SpellLevel int
	Count      int
}

func GenerateBaseCharacter(cfg *DataGenerator.GeneratorConfig) Character {

	// fmt.Println("Races")
	// fmt.Println(cfg.Races.Races)
	// pick core identity
	race := PickWeighted(cfg.Races.Races, func(r DataGenerator.Race) int {
		return r.Weight
	})

	class := PickWeighted(cfg.Classes.Classes, func(c DataGenerator.Class) int {
		return c.Weight
	})

	var validAlignments []DataGenerator.Alignment

	for _, alignment := range cfg.Alignments.Alignments {
		containsAlignmnet := slices.Contains(class.AllowedAlignments, alignment.DisplayName)

		if containsAlignmnet {
			validAlignments = append(validAlignments, alignment)
		}
	}

	alignment := PickWeighted(validAlignments, func(a DataGenerator.Alignment) int {
		return a.Weight
	})

	// attributes
	attrs := DataGenerator.Attribute{
		Strength:     cfg.Attributes.Attributes.Strength,
		Dexterity:    cfg.Attributes.Attributes.Dexterity,
		Constitution: cfg.Attributes.Attributes.Constitution,
		Intelligence: cfg.Attributes.Attributes.Intelligence,
		Wisdom:       cfg.Attributes.Attributes.Wisdom,
		Charisma:     cfg.Attributes.Attributes.Charisma,
	}

	strength := Strength{
		Score: DataGenerator.RollDice(attrs.Strength),
	}

	str := strength
	dex := DataGenerator.RollDice(attrs.Dexterity)
	con := DataGenerator.RollDice(attrs.Constitution)
	intel := DataGenerator.RollDice(attrs.Intelligence)
	wis := DataGenerator.RollDice(attrs.Wisdom)
	cha := DataGenerator.RollDice(attrs.Charisma)

	// biography (simple weighted picks)
	family := PickWeighted(cfg.Biography.Families, func(b DataGenerator.StringProbabilityEntry) int {
		return b.Weight
	})

	raceClan := PickWeighted(cfg.Biography.RaceClans, func(b DataGenerator.StringProbabilityEntry) int {
		return b.Weight
	})

	home := PickWeighted(cfg.Biography.HomeLands, func(b DataGenerator.StringProbabilityEntry) int {
		return b.Weight
	})

	patron := PickWeighted(cfg.Biography.Patrons, func(b DataGenerator.StringProbabilityEntry) int {
		return b.Weight
	})

	religion := PickWeighted(cfg.Biography.Religons, func(b DataGenerator.StringProbabilityEntry) int {
		return b.Weight
	})

	social := PickWeighted(cfg.Biography.SocialClasses, func(b DataGenerator.StringProbabilityEntry) int {
		return b.Weight
	})

	status := PickWeighted(cfg.Biography.Statuses, func(b DataGenerator.StringProbabilityEntry) int {
		return b.Weight
	})

	// physical
	hair := PickWeighted(cfg.Physical.Hair, func(p DataGenerator.StringProbabilityEntry) int {
		return p.Weight
	})

	eyes := PickWeighted(cfg.Physical.Eyes, func(p DataGenerator.StringProbabilityEntry) int {
		return p.Weight
	})

	appearance := PickWeighted(cfg.Physical.Appearances, func(p DataGenerator.StringProbabilityEntry) int {
		return p.Weight
	})

	characterLevel := PickWeighted(cfg.LevelDistribution.IntegerProbabilityEntry, func(b DataGenerator.IntegerProbabilityEntry) int {
		return b.Weight
	})

	return Character{
		Character: CharacterInfo{
			Level:       characterLevel.Value,
			Class:       class.Name,
			ClassConfig: class,
			ClassHitDie: class.HitDie,
			Race:        race.Name,
			Alignment:   alignment.DisplayName,
		},

		Biography: Biography{
			Family:           family.Name,
			RaceClan:         raceClan.Name,
			Homeland:         home.Name,
			Patron:           patron.Name,
			Religion:         religion.Name,
			SocialClass:      social.Name,
			Status:           status.Name,
			BirthRank:        1,
			NumberOfSiblings: 0,
			Honor:            0,
			BaseHonor:        0,
		},

		Physical: Physical{
			Hair:       hair.Name,
			Eyes:       eyes.Name,
			Age:        race.Age.Min, // placeholder (you can randomize later)
			Height:     race.Height.Min,
			Weight:     race.CharacterWeight.Min,
			Appearance: appearance.Name,
		},

		Attributes: Attributes{
			STR: str,
			DEX: dex,
			CON: con,
			INT: intel,
			WIS: wis,
			CHA: cha,
		},

		Economy: Economy{
			Gold: DataGenerator.RollDice(class.StartingGold),
		},

		Spells: []Spell{},
	}
}

func (c Character) String() string {
	var b strings.Builder

	// HEADER BLOCK
	fmt.Fprintf(&b, "========================================================================\n")
	fmt.Fprintf(&b, "  %-45s LEVEL %-3d\n", fmt.Sprintf("%s %s", c.Character.Alignment, c.Character.Race), c.Character.Level)
	fmt.Fprintf(&b, "  Class: %-43s \n", c.Character.Class)
	fmt.Fprintf(&b, "========================================================================\n\n")

	// CORE STATS & COMBAT OVERVIEW
	fmt.Fprintf(&b, "--- CORE VITALITY ------------------------------------------------------\n")
	fmt.Fprintf(&b, "  HP: %-12d AC: %-12d THAC0: %-12d\n", c.DerivedConstitutionStats.Hitpoints, c.Combat.ArmorClass, c.Combat.THACO)
	fmt.Fprintf(&b, "  Surprised: %-5d Rear: %-7d Attack Penalty: %-5d\n\n", c.Combat.Surprised, c.Combat.Rear, c.Combat.AttackPenalty)

	// ATTRIBUTES BLOCK
	fmt.Fprintf(&b, "--- CHARACTER ATTRIBUTES -----------------------------------------------\n")
	strDisplay := fmt.Sprintf("%d", c.Attributes.STR.Score)
	if c.Attributes.STR.Score == 18 && c.Attributes.STR.Exceptional > 0 {
		strDisplay = fmt.Sprintf("18/%02d", c.Attributes.STR.Exceptional)
	}
	fmt.Fprintf(&b, "  STR: %-10s DEX: %-10d CON: %-10d\n", strDisplay, c.Attributes.DEX, c.Attributes.CON)
	fmt.Fprintf(&b, "  INT: %-10d WIS: %-10d CHA: %-10d\n\n", c.Attributes.INT, c.Attributes.WIS, c.Attributes.CHA)

	// DERIVED SUB-STATS BREAKDOWN
	fmt.Fprintf(&b, "--- SUB-STATS BREAKDOWN ------------------------------------------------\n")
	fmt.Fprintf(&b, "  [Strength]     Hit Prob: %+d  |  Dmg Adj: %+d  | Open Doors: %d | Max Press: %d lbs\n", c.DerivedStrengthStats.HitProbability, c.DerivedStrengthStats.DamageAdjustment, c.DerivedStrengthStats.OpenDoors, c.DerivedStrengthStats.MaxPress)
	fmt.Fprintf(&b, "                 Allow: %d  |  BB/LG: %d%%\n", c.DerivedStrengthStats.WeightAllowance, c.DerivedStrengthStats.BendBarsLiftGates)
	fmt.Fprintf(&b, "  [Dexterity]    Reaction Adj: %+d  |  Missile Adj: %+d  |  Defense Adj: %+d\n", c.DerivedDexterityStats.ReactionAdjustment, c.DerivedDexterityStats.MissileAttackAdjustment, c.DerivedDexterityStats.DefenseAdjustment)
	fmt.Fprintf(&b, "  [Constitution] HP Adj: %+d  |  System Shock: %d%%  |  Resurrection: %d%%\n", c.DerivedConstitutionStats.HPAdjustment, c.DerivedConstitutionStats.SystemShock, c.DerivedConstitutionStats.Resurrection)
	fmt.Fprintf(&b, "                 Poison Save: %+d  |  Regen: %d\n", c.DerivedConstitutionStats.PoisonSave, c.DerivedConstitutionStats.Regeneration)
	fmt.Fprintf(&b, "  [Intelligence] Languages Max: %d  |  Max Spell Lvl: %d  |  Learn Spell: %d%%\n", c.DerivedIntelligenceStats.NumberOfLanguages, c.DerivedIntelligenceStats.SpellLevel, c.DerivedIntelligenceStats.LearnSpell)
	fmt.Fprintf(&b, "                 Spells Per Lvl: %d\n", c.DerivedIntelligenceStats.SpellsPerLevel)
	if len(c.DerivedIntelligenceStats.SpellImmunities) > 0 {
		fmt.Fprintf(&b, "                 Immunities: %s\n", strings.Join(c.DerivedIntelligenceStats.SpellImmunities, ", "))
	}
	fmt.Fprintf(&b, "  [Wisdom]       Magic Def Adj: %+d  |  Spell Failure: %d%%\n", c.DerivedWisdomStats.MagicDefenseAdjustment, c.DerivedWisdomStats.SpellFailure)
	if len(c.DerivedWisdomStats.BonusSpells) > 0 {
		var spells []string
		for _, s := range c.DerivedWisdomStats.BonusSpells {
			spells = append(spells, fmt.Sprintf("Lvl %d (+%d)", s.SpellLevel, s.Count))
		}
		fmt.Fprintf(&b, "                 Bonus Spells: %s\n", strings.Join(spells, ", "))
	}
	if len(c.DerivedWisdomStats.SpellImmunities) > 0 {
		fmt.Fprintf(&b, "                 Immunities: %s\n", strings.Join(c.DerivedWisdomStats.SpellImmunities, ", "))
	}
	fmt.Fprintf(&b, "  [Charisma]     Max Henchmen: %d  |  Loyalty Base: %+d  |  Reaction Adj: %+d\n\n", c.DerivedCharismaStats.MaxNumberOfHenchmen, c.DerivedCharismaStats.LoyaltyBase, c.DerivedCharismaStats.ReactionAdjustment)

	// SAVING THROWS
	fmt.Fprintf(&b, "--- SAVING THROWS ------------------------------------------------------\n")
	fmt.Fprintf(&b, "  Paralyze/Poison/Death: %-4d Rod/Staff/Wand: %-4d Petrify/Polymorph: %-4d\n", c.SavingThrows.ParalyzePoisonDeath, c.SavingThrows.RodStaffOrWand, c.SavingThrows.PetrifyPolymorph)
	fmt.Fprintf(&b, "  Breath Weapon:         %-4d Spell:          %-4d\n\n", c.SavingThrows.BreathWeapon, c.SavingThrows.Spell)

	// PROFICIENCIES & CUSTOM TRACKS
	fmt.Fprintf(&b, "--- WEAPON SPECIALIZATIONS ---------------------------------------------\n")
	if len(c.WeaponSpecializations) == 0 {
		fmt.Fprintf(&b, "  (None)\n")
	}
	for _, w := range c.WeaponSpecializations {
		fmt.Fprintf(&b, "  * %s [%s]\n", w.Name, w.Rank)
	}
	fmt.Fprintf(&b, "\n")

	fmt.Fprintf(&b, "--- NON-WEAPON SPECIALIZATIONS -----------------------------------------\n")
	if len(c.NonWeaponSpecializations) == 0 {
		fmt.Fprintf(&b, "  (None)\n")
	}
	for _, n := range c.NonWeaponSpecializations {
		fmt.Fprintf(&b, "  * %-25s | Base Stat: %-12s | Target Roll: <= %d\n", n.Name, n.BaseStat, n.TargetRoll)
	}
	fmt.Fprintf(&b, "\n")

	// UTILITY / FLAVOR
	fmt.Fprintf(&b, "--- LANGUAGES KNOWN ----------------------------------------------------\n  %s\n\n", strings.Join(c.Languages, ", "))

	fmt.Fprintf(&b, "--- EQUIPMENT & ECONOMY ------------------------------------------------\n")
	fmt.Fprintf(&b, "  Wallet: %d gp\n", c.Economy.Gold)
	if len(c.Gear) == 0 {
		fmt.Fprintf(&b, "  Gear: En Route / Empty\n")
	} else {
		for _, g := range c.Gear {
			fmt.Fprintf(&b, "  - %-25s x%-3d (Weight: %d lbs)\n", g.Name, g.Quantity, g.Weight)
		}
	}
	fmt.Fprintf(&b, "\n")

	// BACKGROUNDS
	fmt.Fprintf(&b, "--- BIOGRAPHY & BACKGROUND ---------------------------------------------\n")
	fmt.Fprintf(&b, "  Homeland:   %-20s Family: %s\n", c.Biography.Homeland, c.Biography.Family)
	fmt.Fprintf(&b, "  Clan:       %-20s Patron: %s\n", c.Biography.RaceClan, c.Biography.Patron)
	fmt.Fprintf(&b, "  Religion:   %-20s Social: %s (Status: %s)\n", c.Biography.Religion, c.Biography.SocialClass, c.Biography.Status)
	fmt.Fprintf(&b, "  Birth Rank: %-20d Siblings: %d\n", c.Biography.BirthRank, c.Biography.NumberOfSiblings)
	fmt.Fprintf(&b, "  Honor:      %-20d (Base: %d)\n\n", c.Biography.Honor, c.Biography.BaseHonor)

	fmt.Fprintf(&b, "--- PHYSICAL CHARACTERISTICS -------------------------------------------\n")
	fmt.Fprintf(&b, "  Age: %-10d Height: %-9d Weight: %-10d\n", c.Physical.Age, c.Physical.Height, c.Physical.Weight)
	fmt.Fprintf(&b, "  Hair: %-9s Eyes: %-11s Look: %s\n\n", c.Physical.Hair, c.Physical.Eyes, c.Physical.Appearance)

	// CONDITIONAL ARDUOUS SYSTEM MODULES (Thief / Wizard)
	if c.ThiefData != nil {
		fmt.Fprintf(&b, "--- THIEF SKILLS SKILLSET ----------------------------------------------\n")
		fmt.Fprintf(&b, "  Pick Pocket:     %-4d%%  Open Locks:      %-4d%%  Find/Remove Traps: %-4d%%\n", c.ThiefData.PickPocket, c.ThiefData.OpenLock, c.ThiefData.FindRemoveTrap)
		fmt.Fprintf(&b, "  Move Silently:   %-4d%%  Hide In Shadows: %-4d%%  Detect Noise:      %-4d%%\n", c.ThiefData.MoveSilently, c.ThiefData.HideInShadows, c.ThiefData.DetectNoise)
		fmt.Fprintf(&b, "  Climb Walls:     %-4d%%  Read Languages:  %-4d%%\n\n", c.ThiefData.ClimbWalls, c.ThiefData.ReadLanguages)
	}

	if c.WizardData != nil {
		fmt.Fprintf(&b, "--- WIZARD ARCANUM DATA ------------------------------------------------\n")
		fmt.Fprintf(&b, "  Spell Slots: [L1:%d | L2:%d | L3:%d | L4:%d | L5:%d | L6:%d | L7:%d | L8:%d | L9:%d]\n\n",
			c.WizardData.SpellSlots.Level1, c.WizardData.SpellSlots.Level2, c.WizardData.SpellSlots.Level3,
			c.WizardData.SpellSlots.Level4, c.WizardData.SpellSlots.Level5, c.WizardData.SpellSlots.Level6,
			c.WizardData.SpellSlots.Level7, c.WizardData.SpellSlots.Level8, c.WizardData.SpellSlots.Level9)

		fmt.Fprintf(&b, "  Spellbook Contents:\n")
		if len(c.WizardData.Spellbook) == 0 {
			fmt.Fprintf(&b, "    (Empty)\n")
		}
		for _, spell := range c.WizardData.Spellbook {
			fmt.Fprintf(&b, "    - %s (%s)\n", spell.Name, spell.School)
		}

		fmt.Fprintf(&b, "\n  Currently Memorized Spells:\n")
		if len(c.WizardData.MemorizedSpells) == 0 {
			fmt.Fprintf(&b, "    (None Selected)\n")
		}
		for _, spellName := range c.WizardData.MemorizedSpells {
			fmt.Fprintf(&b, "    - %s\n", spellName)
		}
		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "========================================================================\n")
	return b.String()
}

func PickLevel(cfg DataGenerator.GeneratorConfig) {
	PickWeighted(cfg.LevelDistribution.IntegerProbabilityEntry, func(b DataGenerator.IntegerProbabilityEntry) int {
		return b.Weight
	})
}
