package DataGenerator

import (
	"math/rand"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type GeneratorConfig struct {
	Alignments   Alignments
	Classes      Classes
	Races        Races
	Biography    Biography
	Physical     Physical
	Attributes   Attributes
	SavingThrows SavingThrowTable
	THACO        THACOTable

	StrengthTable     StrengthTable
	DexterityTable    DexterityTable
	ConstitutionTable ConstitutionTable
	WisdomTable       WisdomTable
	IntelligenceTable IntelligenceTable
	CharismaTable     CharismaTable

	ThievingTable    ThievingTable
	SpellSlotTable   SpellSlotTable
	WizardSpellTable WizardSpellTable

	Languages Languages

	Proficiencies    Proficiencies
	ProficiencySlots ProficiencySlots

	GearAssignments GearAssignments
	GearPool        GearPool

	LevelDistribution LevelDistribution
}

type RangeConfig struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type DiceConfig struct {
	DiceCount  int `yaml:"diceCount"`
	DiceSides  int `yaml:"diceSides"`
	Multiplier int `yaml:"multiplier,omitempty"`
}

type Alignments struct {
	Alignments []Alignment `yaml:"alignments"`
}

type Alignment struct {
	DisplayName string `yaml:"displayName"`
	Weight      int    `yaml:"weight"`
}

type Classes struct {
	Classes []Class `yaml:"classes"`
}

type Class struct {
	Name   string `yaml:"name"`
	Weight int    `yaml:"weight"`

	StartingGold DiceConfig `yaml:"startingGold"`

	HitDie         int `yaml:"hitDie"`
	MaxHitDieLevel int `yaml:"maxHitDieLevel"`

	HasThiefSkills bool `yaml:"hasThiefSkills"`
	CanCastArcane  bool `yaml:"canCastArcane"`
	CanCastDivine  bool `yaml:"canCastDivine"`

	MinStats MinStats `yaml:"minStats"`

	AllowedAlignments []string `yaml:"allowedAlignments"`
}

type MinStats struct {
	Strength     int `yaml:"strength"`
	Dexterity    int `yaml:"dexterity"`
	Constitution int `yaml:"constitution"`
	Intelligence int `yaml:"intelligence"`
	Wisdom       int `yaml:"wisdom"`
	Charisma     int `yaml:"charisma"`
}
type Races struct {
	Races []Race `yaml:"races"`
}

type Race struct {
	Name   string `yaml:"name"`
	Weight int    `yaml:"weight"`

	Age             RangeConfig `yaml:"age"`
	Height          RangeConfig `yaml:"height"`
	CharacterWeight RangeConfig `yaml:"characterWeight"`
}

type StringProbabilityEntry struct {
	Name   string `yaml:"name"`
	Weight int    `yaml:"weight"`
}

type Biography struct {
	Families      []StringProbabilityEntry `yaml:"families"`
	RaceClans     []StringProbabilityEntry `yaml:"raceClans"`
	HomeLands     []StringProbabilityEntry `yaml:"homeLands"`
	Patrons       []StringProbabilityEntry `yaml:"patrons"`
	Religons      []StringProbabilityEntry `yaml:"religons"`
	SocialClasses []StringProbabilityEntry `yaml:"socialClasses"`
	Statuses      []StringProbabilityEntry `yaml:"statuses"`
}

type Physical struct {
	Hair        []StringProbabilityEntry `yaml:"hair"`
	Eyes        []StringProbabilityEntry `yaml:"eyes"`
	Appearances []StringProbabilityEntry `yaml:"appearances"`
}

type Attributes struct {
	Attributes Attribute `yaml:"attributes"`
}

type Attribute struct {
	Strength     DiceConfig `yaml:"strength"`
	Dexterity    DiceConfig `yaml:"dexterity"`
	Constitution DiceConfig `yaml:"constitution"`
	Intelligence DiceConfig `yaml:"intelligence"`
	Wisdom       DiceConfig `yaml:"wisdom"`
	Charisma     DiceConfig `yaml:"charisma"`
}

type SavingThrowTable struct {
	SavingThrows map[string][]SavingThrowRow `yaml:"savingThrows"`
}

type SavingThrowRow struct {
	MinLevel int `yaml:"minLevel"`
	MaxLevel int `yaml:"maxLevel"`

	ParalyzePoisonDeath int `yaml:"paralyzePoisonDeath"`
	PetrifyPolymorph    int `yaml:"petrifyPolymorph"`
	RodStaffOrWand      int `yaml:"rodStaffOrWand"`
	BreathWeapon        int `yaml:"breathWeapon"`
	Spell               int `yaml:"spell"`
}

type THACOTable struct {
	THACO map[string][]THACORow `yaml:"thaco"`
}

type THACORow struct {
	MinLevel int `yaml:"minLevel"`
	MaxLevel int `yaml:"maxLevel"`
	Value    int `yaml:"value"`
}

type StrengthTable struct {
	Entries []StrengthEntry `yaml:"strengthTable"`
}

type StrengthEntry struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`

	ExceptionalMin int `yaml:"exceptionalMin,omitempty"`
	ExceptionalMax int `yaml:"exceptionalMax,omitempty"`

	HitProbability    int `yaml:"hitProbability"`
	DamageAdjustment  int `yaml:"damageAdjustment"`
	WeightAllowance   int `yaml:"weightAllowance"`
	MaxPress          int `yaml:"maxPress"`
	OpenDoors         int `yaml:"openDoors"`
	BendBarsLiftGates int `yaml:"bendBarsLiftGates"`
}

func LoadComponent[T any](target *T, filepath string) error {
	payload, err := LoadYAMLFile[T](filepath)

	if err != nil {
		return err
	}

	*target = payload

	return nil
}

func LoadYAMLFile[T any](filepath string) (T, error) {
	var cfg T

	data, err := os.ReadFile(filepath)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

type DexterityTable struct {
	Entries []DexterityEntry `yaml:"dexterityTable"`
}

type DexterityEntry struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`

	ReactionAdjustment      int `yaml:"reactionAdjustment"`
	MissileAttackAdjustment int `yaml:"missileAttackAdjustment"`
	DefenseAdjustment       int `yaml:"defenseAdjustment"`
}

type ConstitutionTable struct {
	Entries []ConstitutionEntry `yaml:"constitutionTable"`
}

type ConstitutionEntry struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`

	HPAdjustment int `yaml:"hpAdjustment"`
	SystemShock  int `yaml:"systemShock"`
	Resurrection int `yaml:"resurrection"`
	PoisonSave   int `yaml:"poisonSave"`
	Regeneration int `yaml:"regeneration"`
}

type IntelligenceTable struct {
	Entries []IntelligenceEntry `yaml:"intelligenceTable"`
}

type IntelligenceEntry struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`

	NumberOfLanguages int `yaml:"numberOfLanguages"`
	SpellLevel        int `yaml:"spellLevel"`
	LearnSpell        int `yaml:"learnSpell"`
	SpellsPerLevel    int `yaml:"spellsPerLevel"`

	SpellImmunities []string `yaml:"spellImmunities"`
}

type WisdomTable struct {
	Entries []WisdomEntry `yaml:"wisdomTable"`
}

type WisdomEntry struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`

	MagicDefenseAdjustment int `yaml:"magicDefenseAdjustment"`
	SpellFailure           int `yaml:"spellFailure"`

	BonusSpells     []BonusSpellEntry `yaml:"bonusSpells"`
	SpellImmunities []string          `yaml:"spellImmunities"`
}

type BonusSpellEntry struct {
	SpellLevel int `yaml:"spellLevel"`
	Count      int `yaml:"count"`
}

type CharismaTable struct {
	Entries []CharismaEntry `yaml:"charismaTable"`
}

type CharismaEntry struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`

	MaxNumberOfHenchmen int `yaml:"maxNumberOfHenchmen"`
	LoyaltyBase         int `yaml:"loyaltyBase"`
	ReactionAdjustment  int `yaml:"reactionAdjustment"`
}

type RangeEntry interface {
	GetMin() int
	GetMax() int
}

type ThievingTable struct {
	PickPocket               int `yaml:"pickpocket"`
	OpenLock                 int `yaml:"openLock"`
	FindRemoveTrap           int `yaml:"findRemoveTrap"`
	MoveSilently             int `yaml:"moveSilently"`
	HideInShadows            int `yaml:"hideInShadows"`
	DetectNoise              int `yaml:"detectNoise"`
	ClimbWalls               int `yaml:"climbWalls"`
	ReadLanguages            int `yaml:"readLanguages"`
	PointsPerLevel           int `yaml:"pointsPerLevel"`
	SkillIncreasePerLevelCap int `yaml:"skillIncreasePerLevelCap"`
	Level1Variance           int `yaml:"level1Variance"`
}

func (e DexterityEntry) GetMin() int { return e.Min }
func (e DexterityEntry) GetMax() int { return e.Max }

func (e ConstitutionEntry) GetMin() int { return e.Min }
func (e ConstitutionEntry) GetMax() int { return e.Max }

func (e IntelligenceEntry) GetMin() int { return e.Min }
func (e IntelligenceEntry) GetMax() int { return e.Max }

func (e WisdomEntry) GetMin() int { return e.Min }
func (e WisdomEntry) GetMax() int { return e.Max }

func (e CharismaEntry) GetMin() int { return e.Min }
func (e CharismaEntry) GetMax() int { return e.Max }

type Languages struct {
	Languages []LanguageEntry `yaml:"languages"`
}

type LanguageEntry struct {
	Name   string `yaml:"name"`
	Weight int    `yaml:"weight"`
}

type Proficiencies struct {
	WeaponProficiencies    []ProficiencyEntry `yaml:"weaponProficiencies"`
	NonWeaponProficiencies []ProficiencyEntry `yaml:"nonWeaponProficiencies"`
}

type ProficiencyEntry struct {
	Name           string   `yaml:"name"`
	Weight         int      `yaml:"weight"`
	AllowedClasses []string `yaml:"allowedClasses"`
	Stat           string   `yaml:"stat,omitempty"`     // e.g., "STR", "DEX", "WIS"
	Modifier       int      `yaml:"modifier,omitempty"` // e.g., 0, -1, -3
}

type ProficiencySlots struct {
	ProficiencySlots map[string]ProficiencySlotGroup `yaml:"proficiencySlots"`
}

type ProficiencySlotGroup struct {
	Weapon    ProficiencySlotRule `yaml:"weapon"`
	NonWeapon ProficiencySlotRule `yaml:"nonWeapon"`
}

type ProficiencySlotRule struct {
	Initial   int `yaml:"initial"`
	PerLevels int `yaml:"perLevels"`
	Penalty   int `yaml:"penalty,omitempty"`
}

type SpellSlotTable struct {
	WizardSpellSlots []WizardSpellSlotRow `yaml:"wizardSpellSlots"`
}

type WizardSpellSlotRow struct {
	MinLevel int `yaml:"minLevel"`
	MaxLevel int `yaml:"maxLevel"`

	Slots SpellSlotsConfig `yaml:"slots"`
}

type SpellSlotsConfig struct {
	Level1 int `yaml:"level1"`
	Level2 int `yaml:"level2"`
	Level3 int `yaml:"level3"`
	Level4 int `yaml:"level4"`
	Level5 int `yaml:"level5"`
	Level6 int `yaml:"level6"`
	Level7 int `yaml:"level7"`
	Level8 int `yaml:"level8"`
	Level9 int `yaml:"level9"`
}

type WizardSpellTable struct {
	Spells []WizardSpellEntry `yaml:"Spells"`
}

type WizardSpellEntry struct {
	Name           string   `yaml:"name"`
	Level          int      `yaml:"level"`
	Sphere         string   `yaml:"sphere"`
	Weight         int      `yaml:"weight"`
	AllowedClasses []string `yaml:"allowedClasses"`
}

type GearAssignments struct {
	GearAssignments map[string]GearAssignmentRule `yaml:"gearAssignments"`
}

type GearAssignmentRule struct {
	Guaranteed  []string `yaml:"guaranteed"`
	RandomSlots int      `yaml:"randomSlots"`
}

type GearPool struct {
	Pool []GearPoolEntry `yaml:"gear"`
}

type GearPoolEntry struct {
	Name              string `yaml:"name"`
	Category          string `yaml:"category"`
	Cost              string `yaml:"cost"`
	Weight            string `yaml:"weight"`
	ProbabilityWeight int    `yaml:"probabilityWeight"`
}

type LevelDistribution struct {
	IntegerProbabilityEntry []IntegerProbabilityEntry `yaml:"levelDistribution"`
}

type IntegerProbabilityEntry struct {
	Value  int `yaml:"value"`
	Weight int `yaml:"weight"`
}

func LoadGeneratorConfig(rulePath string) (*GeneratorConfig, error) {
	var cfg GeneratorConfig
	attributesFilepath := path.Join(rulePath, "attributes.yml")
	alignmentsFilepath := path.Join(rulePath, "alignments.yml")
	classesFilepath := path.Join(rulePath, "classes.yml")
	racesFilepath := path.Join(rulePath, "races.yml")
	biographyFilepath := path.Join(rulePath, "biography.yml")
	physicalFilepath := path.Join(rulePath, "physical.yml")
	thacoFilepath := path.Join(rulePath, "thaco.yml")
	savingThrowsFilepath := path.Join(rulePath, "savingthrows.yml")
	strengthFilepath := path.Join(rulePath, "./tables/strength.yml")
	dexterityFilepath := path.Join(rulePath, "./tables/dexterity.yml")
	constitutionFilepath := path.Join(rulePath, "./tables/constitution.yml")
	intelligenceFilepath := path.Join(rulePath, "./tables/intelligence.yml")
	wisdomFilepath := path.Join(rulePath, "./tables/wisdom.yml")
	charismaFilepath := path.Join(rulePath, "./tables/charisma.yml")
	languagesFilepath := path.Join(rulePath, "./language.yml")

	thievingTableFilepath := path.Join(rulePath, "./tables/thieving.yml")

	proficiencyFilepath := path.Join(rulePath, "./proficiencies.yml")
	proficiencySlotsFilepath := path.Join(rulePath, "./proficiency_slots.yml")

	spellSlotsFilepath := path.Join(rulePath, "./tables/wizardSpellSlot.yml")

	spellsFilepath := path.Join(rulePath, "./tables/wizardSpell.yml")

	gearPoolFilepath := path.Join(rulePath, "./tables/gearAssignmentTable.yml")
	gearFilepath := path.Join(rulePath, "./gear.yml")
	levelDistributionsFilepath := path.Join(rulePath, "./levelDistribution.yml")

	if err := LoadComponent(&cfg.THACO, thacoFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.SavingThrows, savingThrowsFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.Attributes, attributesFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.Alignments, alignmentsFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.Classes, classesFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.Races, racesFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.Biography, biographyFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.Physical, physicalFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.StrengthTable, strengthFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.DexterityTable, dexterityFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.ConstitutionTable, constitutionFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.IntelligenceTable, intelligenceFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.WisdomTable, wisdomFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.CharismaTable, charismaFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.Languages, languagesFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.Proficiencies, proficiencyFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.ProficiencySlots, proficiencySlotsFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.ThievingTable, thievingTableFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.SpellSlotTable, spellSlotsFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.WizardSpellTable, spellsFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.GearAssignments, gearPoolFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.GearPool, gearFilepath); err != nil {
		return nil, err
	}

	if err := LoadComponent(&cfg.LevelDistribution, levelDistributionsFilepath); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func RollDice(cfg DiceConfig) int {
	total := 0

	for i := 0; i < cfg.DiceCount; i++ {
		total += rand.Intn(cfg.DiceSides) + 1
	}

	if cfg.Multiplier != 0 {
		total *= cfg.Multiplier
	}

	return total
}
