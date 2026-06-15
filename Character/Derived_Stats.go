package character

import (
	"fmt"
	"john-v3/dndCharacterGenerator/DataGenerator"
	"math"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
)

func ApplyDerivedStats(c *Character, ruleSet *DataGenerator.GeneratorConfig, level int) {
	validateStatsVsClass(c)

	applyStrength(c, ruleSet)
	applyDexterity(c, ruleSet)
	applyConstitution(c, ruleSet)

	applyIntelligence(c, ruleSet)
	applyWisdom(c, ruleSet)
	applyCharisma(c, ruleSet)
	ApplyTHACO(c, &ruleSet.THACO)
	ApplySavingThrows(c, &ruleSet.SavingThrows)

	c.DerivedConstitutionStats.Hitpoints = RollHitPoints(c.Character.Level, c.Character.ClassHitDie, c.DerivedConstitutionStats.HPAdjustment)

	c.Combat.ArmorClass = 10 + c.DerivedDexterityStats.DefenseAdjustment

	c.Languages = generateLanguages(ruleSet, c.DerivedIntelligenceStats.NumberOfLanguages)

	weaponSlots, nonWeaponSlots, penalty := getProficiencySlots(ruleSet, c.Character.Class, c.Character.Level)

	c.WeaponSpecializations = generateWeaponProficiencies(ruleSet, c.Character.Class, weaponSlots)
	c.NonWeaponSpecializations = generateNonWeaponProficiencies(ruleSet, c, nonWeaponSlots)
	c.Combat.AttackPenalty = penalty

	c.Gear = generateGear(ruleSet, c)

	thiefSkills(c, ruleSet)
	wizardSpells(c, ruleSet)
}
func applyStrength(c *Character, cfg *DataGenerator.GeneratorConfig) {

	strengthRow, err := LookupStrength(
		cfg.StrengthTable,
		c.Attributes.STR.Score,
		c.Attributes.STR.Exceptional,
	)

	if err != nil {
		panic(err)
	}

	c.DerivedStrengthStats = DerivedStrengthStats{
		HitProbability:    strengthRow.HitProbability,
		DamageAdjustment:  strengthRow.DamageAdjustment,
		WeightAllowance:   strengthRow.WeightAllowance,
		MaxPress:          strengthRow.MaxPress,
		OpenDoors:         strengthRow.OpenDoors,
		BendBarsLiftGates: strengthRow.BendBarsLiftGates,
	}
}

func applyDexterity(c *Character, cfg *DataGenerator.GeneratorConfig) {

	dexterityRow, err := LookupByScore(
		c.Attributes.DEX,
		cfg.DexterityTable.Entries,
	)

	if err != nil {
		panic(err)
	}

	c.DerivedDexterityStats = DerivedDexterityStats{
		ReactionAdjustment:      dexterityRow.ReactionAdjustment,
		MissileAttackAdjustment: dexterityRow.MissileAttackAdjustment,
		DefenseAdjustment:       dexterityRow.DefenseAdjustment,
	}
}

func applyConstitution(c *Character, cfg *DataGenerator.GeneratorConfig) {

	constitutionRow, err := LookupByScore(
		c.Attributes.CON,
		cfg.ConstitutionTable.Entries,
	)

	if err != nil {
		panic(err)
	}

	c.DerivedConstitutionStats = DerivedConstitutionStats{
		HPAdjustment: constitutionRow.HPAdjustment,
		SystemShock:  constitutionRow.SystemShock,
		Resurrection: constitutionRow.Resurrection,
		PoisonSave:   constitutionRow.PoisonSave,
		Regeneration: constitutionRow.Regeneration,
	}
}

func applyIntelligence(c *Character, cfg *DataGenerator.GeneratorConfig) {

	// fmt.Println("int: ", c.Attributes.INT)
	intelligenceRow, err := LookupByScore(
		c.Attributes.INT,
		cfg.IntelligenceTable.Entries,
	)

	if err != nil {
		panic(err)
	}

	c.DerivedIntelligenceStats = DerivedIntelligenceStats{
		NumberOfLanguages: intelligenceRow.NumberOfLanguages,
		SpellLevel:        intelligenceRow.SpellLevel,
		LearnSpell:        intelligenceRow.LearnSpell,
		SpellsPerLevel:    intelligenceRow.SpellsPerLevel,
		// spell immunities?
	}
}

func applyWisdom(c *Character, cfg *DataGenerator.GeneratorConfig) {

	wisdomRow, err := LookupByScore(
		c.Attributes.WIS,
		cfg.WisdomTable.Entries,
	)

	if err != nil {
		panic(err)
	}

	c.DerivedWisdomStats = DerivedWisdomStats{
		MagicDefenseAdjustment: wisdomRow.MagicDefenseAdjustment,
		BonusSpells:            wisdomRow.BonusSpells,
		SpellFailure:           wisdomRow.SpellFailure,
		SpellImmunities:        wisdomRow.SpellImmunities,
	}
}

func applyCharisma(c *Character, cfg *DataGenerator.GeneratorConfig) {

	charismaRow, err := LookupByScore(
		c.Attributes.CHA,
		cfg.CharismaTable.Entries,
	)

	if err != nil {
		panic(err)
	}

	c.DerivedCharismaStats = DerivedCharismaStats{
		MaxNumberOfHenchmen: charismaRow.MaxNumberOfHenchmen,
		LoyaltyBase:         charismaRow.LoyaltyBase,
		ReactionAdjustment:  charismaRow.ReactionAdjustment,
	}
}

func ApplySavingThrows(
	c *Character,
	table *DataGenerator.SavingThrowTable,
) error {

	row, err := FindSavingThrowRow(
		table,
		c.Character.Class,
		c.Character.Level,
	)

	if err != nil {
		return err
	}

	c.SavingThrows = SavingThrows{
		ParalyzePoisonDeath: row.ParalyzePoisonDeath,
		RodStaffOrWand:      row.RodStaffOrWand,
		PetrifyPolymorph:    row.PetrifyPolymorph,
		BreathWeapon:        row.BreathWeapon,
		Spell:               row.Spell,
	}

	return nil
}

func FindSavingThrowRow(
	table *DataGenerator.SavingThrowTable,
	class string,
	level int,
) (*DataGenerator.SavingThrowRow, error) {

	rows, ok := table.SavingThrows[class]
	if !ok {
		return nil, fmt.Errorf("unknown class %s", class)
	}

	for _, row := range rows {
		if level >= row.MinLevel &&
			level <= row.MaxLevel {

			return &row, nil
		}
	}

	return nil, fmt.Errorf(
		"no saving throw row found for class=%s level=%d",
		class,
		level,
	)
}

func FindTHACO(
	table *DataGenerator.THACOTable,
	class string,
	level int,
) (int, error) {

	rows, ok := table.THACO[class]
	if !ok {
		return 0, fmt.Errorf("unknown class %s", class)
	}

	for _, row := range rows {
		if level >= row.MinLevel &&
			level <= row.MaxLevel {
			return row.Value, nil
		}
	}

	return 0, fmt.Errorf(
		"no THACO found for class=%s level=%d",
		class,
		level,
	)
}

func ApplyTHACO(
	c *Character,
	table *DataGenerator.THACOTable,
) error {

	thaco, err := FindTHACO(
		table,
		c.Character.Class,
		c.Character.Level,
	)

	if err != nil {
		return err
	}

	c.Combat.THACO = thaco

	return nil
}

func LookupByScore[T DataGenerator.RangeEntry](score int, rows []T) (T, error) {

	for _, entry := range rows {

		if score < entry.GetMin() || score > entry.GetMax() {
			continue
		}

		return entry, nil
	}

	panic("no record found when looking up score...")
}

func LookupStrength(
	table DataGenerator.StrengthTable,
	score int,
	exceptional int,
) (*DataGenerator.StrengthEntry, error) {

	for _, entry := range table.Entries {

		if score < entry.Min || score > entry.Max {
			continue
		}

		if score != 18 {
			return &entry, nil
		}

		if entry.ExceptionalMin == 0 &&
			entry.ExceptionalMax == 0 &&
			exceptional == 0 {
			return &entry, nil
		}

		if exceptional >= entry.ExceptionalMin &&
			exceptional <= entry.ExceptionalMax {
			return &entry, nil
		}
	}

	return nil, fmt.Errorf(
		"strength lookup failed: score=%d exceptional=%d",
		score,
		exceptional,
	)
}

func RollHitPoints(level int, hitDie int, hpAdjustment int) int {
	total := 0

	for range level {
		total += rand.Intn(hitDie) + 1
		total += hpAdjustment
	}

	if total < 1 {
		total = 1
	}

	return total
}

func IsAlignmentAllowed(class DataGenerator.Class, alignment string) bool {
	if slices.Contains(class.AllowedAlignments, alignment) {
		return true
	}
	return false
}

func generateLanguages(cfg *DataGenerator.GeneratorConfig, count int) []string {
	if count <= 0 {
		return []string{"Common"}
	}

	// build pool from config
	pool := make([]DataGenerator.LanguageEntry, len(cfg.Languages.Languages))
	copy(pool, cfg.Languages.Languages)

	// ensure Common always exists
	langs := []string{"Common"}

	// remove Common from pool so we don’t duplicate it
	filtered := pool[:0]
	for _, l := range pool {
		if l.Name != "Common" {
			filtered = append(filtered, l)
		}
	}
	pool = filtered

	// weighted shuffle selection (simple repeated pick)
	for len(langs) < count && len(pool) > 0 {
		pick := PickWeighted(pool, func(l DataGenerator.LanguageEntry) int {
			return l.Weight
		})

		// avoid duplicates
		if !slices.Contains(langs, pick.Name) {
			langs = append(langs, pick.Name)
		}
	}

	return langs
}

func filterByClass(entries []DataGenerator.ProficiencyEntry, class string) []DataGenerator.ProficiencyEntry {
	out := make([]DataGenerator.ProficiencyEntry, 0, len(entries))

	for _, e := range entries {
		if len(e.AllowedClasses) == 0 {
			out = append(out, e)
			continue
		}

		if slices.Contains(e.AllowedClasses, class) {
			out = append(out, e)
		}
	}

	return out
}

func generateWeaponProficiencies(cfg *DataGenerator.GeneratorConfig, class string, slots int) []WeaponSpecialization {
	pool := filterByClass(cfg.Proficiencies.WeaponProficiencies, class)

	if len(pool) == 0 {
		return nil
	}

	result := make([]WeaponSpecialization, 0, slots)

	for len(result) < slots {
		pick := PickWeighted(pool, func(p DataGenerator.ProficiencyEntry) int {
			return p.Weight
		})

		// prevent duplicates
		already := false
		for _, r := range result {
			if r.Name == pick.Name {
				already = true
				break
			}
		}

		if !already {
			result = append(result, WeaponSpecialization{
				Name: pick.Name,
				Rank: "Basic", // placeholder for now
			})
		}
	}

	return result
}

func generateNonWeaponProficiencies(cfg *DataGenerator.GeneratorConfig, c *Character, slots int) []NonWeaponSpecialization {
	pool := filterByClass(cfg.Proficiencies.NonWeaponProficiencies, c.Character.Class)

	if len(pool) == 0 {
		return nil
	}

	result := make([]NonWeaponSpecialization, 0, slots)

	for len(result) < slots {
		pick := PickWeighted(pool, func(p DataGenerator.ProficiencyEntry) int {
			return p.Weight
		})

		already := false
		for _, r := range result {
			if r.Name == pick.Name {
				already = true
				break
			}
		}

		if !already {
			// Pass the entire parsed proficiency configuration object directly
			baseStat, targetRoll := calculateNWPTarget(pick, c)

			result = append(result, NonWeaponSpecialization{
				Name:       pick.Name,
				BaseStat:   baseStat,
				TargetRoll: targetRoll,
			})
		}
	}

	return result
}

func calculateNWPTarget(p DataGenerator.ProficiencyEntry, c *Character) (string, int) {
	// Fall back to a standard unmodified Intelligence check if someone forgets config values
	if p.Stat == "" {
		return "INT", c.Attributes.INT
	}

	var baseScore int

	switch p.Stat {
	case "STR":
		baseScore = c.Attributes.STR.Score
	case "DEX":
		baseScore = c.Attributes.DEX
	case "CON":
		baseScore = c.Attributes.CON
	case "INT":
		baseScore = c.Attributes.INT
	case "WIS":
		baseScore = c.Attributes.WIS
	case "CHA":
		baseScore = c.Attributes.CHA
	default:
		baseScore = c.Attributes.INT // Ultimate safety net fallback
	}

	// AD&D Rule calculation: (Core Attribute Score + Checked Penalty/Bonus Modifier)
	return p.Stat, baseScore + p.Modifier
}

func calculateSlots(rule DataGenerator.ProficiencySlotRule, level int) int {
	if level <= 0 {
		return 0
	}

	// base slots
	slots := rule.Initial

	// growth over time
	if rule.PerLevels > 0 {
		slots += (level - 1) / rule.PerLevels
	}

	return slots
}

func getProficiencySlots(cfg *DataGenerator.GeneratorConfig, class string, level int) (weapon int, nonWeapon int, penalty int) {

	group, ok := cfg.ProficiencySlots.ProficiencySlots[class]
	if !ok {
		// fallback safety
		panic("error getting proficiency slots")
	}

	weapon = calculateSlots(group.Weapon, level)
	nonWeapon = calculateSlots(group.NonWeapon, level)
	penalty = group.Weapon.Penalty

	return
}

func thiefSkills(c *Character, cfg *DataGenerator.GeneratorConfig) {

	if c.Character.Class != "thief" {
		return
	}

	if c.ThiefData == nil {
		c.ThiefData = &ThiefData{}
	}

	c.ThiefData.PickPocket = cfg.ThievingTable.PickPocket
	c.ThiefData.OpenLock = cfg.ThievingTable.OpenLock
	c.ThiefData.FindRemoveTrap = cfg.ThievingTable.FindRemoveTrap
	c.ThiefData.MoveSilently = cfg.ThievingTable.MoveSilently
	c.ThiefData.HideInShadows = cfg.ThievingTable.HideInShadows
	c.ThiefData.DetectNoise = cfg.ThievingTable.DetectNoise
	c.ThiefData.ClimbWalls = cfg.ThievingTable.ClimbWalls
	c.ThiefData.ReadLanguages = cfg.ThievingTable.ReadLanguages

	allocateThiefPoints(c, 60, 30)

	for i := 2; i <= c.Character.Level; i++ {
		allocateThiefPoints(c, cfg.ThievingTable.PointsPerLevel, cfg.ThievingTable.SkillIncreasePerLevelCap)
	}

	// fmt.Print("theifData", c.ThiefData)
}

func allocateThiefPoints(c *Character, totalPoints int, maxPerSkill int) {
	// Track points invested in this specific allocation window/level loop
	investedThisWindow := map[string]int{
		"PickPocket":     0,
		"OpenLock":       0,
		"FindRemoveTrap": 0,
		"MoveSilently":   0,
		"HideInShadows":  0,
		"DetectNoise":    0,
		"ClimbWalls":     0,
		"ReadLanguages":  0,
	}

	pointsRemaining := totalPoints

	// Keep looping until all discretionary points for this level/window are spent
	for pointsRemaining > 0 {
		// Identify which skills are still legally allowed to receive points right now
		validSkills := make([]string, 0, 8)

		if c.ThiefData.PickPocket < 95 && investedThisWindow["PickPocket"] < maxPerSkill {
			validSkills = append(validSkills, "PickPocket")
		}
		if c.ThiefData.OpenLock < 95 && investedThisWindow["OpenLock"] < maxPerSkill {
			validSkills = append(validSkills, "OpenLock")
		}
		if c.ThiefData.FindRemoveTrap < 95 && investedThisWindow["FindRemoveTrap"] < maxPerSkill {
			validSkills = append(validSkills, "FindRemoveTrap")
		}
		if c.ThiefData.MoveSilently < 95 && investedThisWindow["MoveSilently"] < maxPerSkill {
			validSkills = append(validSkills, "MoveSilently")
		}
		if c.ThiefData.HideInShadows < 95 && investedThisWindow["HideInShadows"] < maxPerSkill {
			validSkills = append(validSkills, "HideInShadows")
		}
		if c.ThiefData.DetectNoise < 95 && investedThisWindow["DetectNoise"] < maxPerSkill {
			validSkills = append(validSkills, "DetectNoise")
		}
		if c.ThiefData.ClimbWalls < 95 && investedThisWindow["ClimbWalls"] < maxPerSkill {
			validSkills = append(validSkills, "ClimbWalls")
		}
		// In 2nd Edition, Read Languages is unavailable at level 1 base (0%),
		// but can be invested in normally as long as it doesn't break caps.
		if c.ThiefData.ReadLanguages < 95 && investedThisWindow["ReadLanguages"] < maxPerSkill {
			validSkills = append(validSkills, "ReadLanguages")
		}

		// Safety check: if all skills are maxed out at 95%, discard remaining points
		if len(validSkills) == 0 {
			break
		}

		// Pick a random eligible skill to receive a point
		targetSkill := validSkills[rand.Intn(len(validSkills))]

		// Allocate 1 point to the character and track it against caps
		switch targetSkill {
		case "PickPocket":
			c.ThiefData.PickPocket += 5
		case "OpenLock":
			c.ThiefData.OpenLock += 5
		case "FindRemoveTrap":
			c.ThiefData.FindRemoveTrap += 5
		case "MoveSilently":
			c.ThiefData.MoveSilently += 5
		case "HideInShadows":
			c.ThiefData.HideInShadows += 5
		case "DetectNoise":
			c.ThiefData.DetectNoise += 5
		case "ClimbWalls":
			c.ThiefData.ClimbWalls += 5
		case "ReadLanguages":
			c.ThiefData.ReadLanguages += 5
		}

		investedThisWindow[targetSkill] += 5
		pointsRemaining -= 5
	}
}

func FindWizardSpellSlots(
	table *DataGenerator.SpellSlotTable,
	level int,
) (*DataGenerator.WizardSpellSlotRow, error) {

	for _, row := range table.WizardSpellSlots {

		if level >= row.MinLevel &&
			level <= row.MaxLevel {

			return &row, nil
		}
	}

	return nil, fmt.Errorf(
		"no spell slot row found for wizard level %d",
		level,
	)
}

func ApplyWizardSpellSlots(
	c *Character,
	table *DataGenerator.SpellSlotTable,
) error {

	if c.WizardData == nil {
		c.WizardData = &WizardData{}
	}

	row, err := FindWizardSpellSlots(
		table,
		c.Character.Level,
	)

	if err != nil {
		return err
	}

	c.WizardData.SpellSlots = SpellSlots{
		Level1: row.Slots.Level1,
		Level2: row.Slots.Level2,
		Level3: row.Slots.Level3,
		Level4: row.Slots.Level4,
		Level5: row.Slots.Level5,
		Level6: row.Slots.Level6,
		Level7: row.Slots.Level7,
		Level8: row.Slots.Level8,
		Level9: row.Slots.Level9,
	}

	return nil
}

func filterSpellsByClass(entries []DataGenerator.WizardSpellEntry, class string) []DataGenerator.WizardSpellEntry {
	out := make([]DataGenerator.WizardSpellEntry, 0, len(entries))

	for _, e := range entries {
		if slices.Contains(e.AllowedClasses, class) {
			out = append(out, e)
		}
	}

	return out
}

func spellFromEntry(e DataGenerator.WizardSpellEntry) Spell {
	return Spell{
		Name:   e.Name,
		Level:  e.Level,
		School: e.Sphere,
	}
}

func generateSpellbook(c *Character, cfg *DataGenerator.GeneratorConfig) []Spell {
	pool := filterSpellsByClass(cfg.WizardSpellTable.Spells, c.Character.Class)

	slotsByLevel := map[int]int{
		1: c.WizardData.SpellSlots.Level1,
		2: c.WizardData.SpellSlots.Level2,
		3: c.WizardData.SpellSlots.Level3,
		4: c.WizardData.SpellSlots.Level4,
		5: c.WizardData.SpellSlots.Level5,
		6: c.WizardData.SpellSlots.Level6,
		7: c.WizardData.SpellSlots.Level7,
		8: c.WizardData.SpellSlots.Level8,
		9: c.WizardData.SpellSlots.Level9,
	}

	spellbook := make([]Spell, 0)

	for level := 1; level <= 9; level++ {

		count := slotsByLevel[level]
		if count == 0 {
			continue
		}

		var levelPool []DataGenerator.WizardSpellEntry
		for _, s := range pool {
			if s.Level == level {
				levelPool = append(levelPool, s)
			}
		}

		known := 0
		attempts := 0

		for known < count && len(levelPool) > 0 && attempts < len(levelPool)*4 {
			attempts++

			pick := PickWeighted(levelPool, func(s DataGenerator.WizardSpellEntry) int {
				return s.Weight
			})

			alreadyKnown := false
			for _, existing := range spellbook {
				if existing.Name == pick.Name {
					alreadyKnown = true
					break
				}
			}

			if alreadyKnown {
				continue
			}

			// Learn Spell % chance, 2e style: roll d100, succeed if <= LearnSpell
			if rand.Intn(100)+1 > c.DerivedIntelligenceStats.LearnSpell {
				continue
			}

			spellbook = append(spellbook, spellFromEntry(pick))
			known++
		}
	}

	return spellbook
}

func memorizeSpells(c *Character) []string {
	slotsByLevel := map[int]int{
		1: c.WizardData.SpellSlots.Level1,
		2: c.WizardData.SpellSlots.Level2,
		3: c.WizardData.SpellSlots.Level3,
		4: c.WizardData.SpellSlots.Level4,
		5: c.WizardData.SpellSlots.Level5,
		6: c.WizardData.SpellSlots.Level6,
		7: c.WizardData.SpellSlots.Level7,
		8: c.WizardData.SpellSlots.Level8,
		9: c.WizardData.SpellSlots.Level9,
	}

	memorized := make([]string, 0)

	for level := 1; level <= 9; level++ {

		count := slotsByLevel[level]
		if count == 0 {
			continue
		}

		var pool []Spell
		for _, s := range c.WizardData.Spellbook {
			if s.Level == level {
				pool = append(pool, s)
			}
		}

		for i := 0; i < count && len(pool) > 0; i++ {
			idx := rand.Intn(len(pool))
			memorized = append(memorized, pool[idx].Name)
			pool = append(pool[:idx], pool[idx+1:]...)
		}
	}

	return memorized
}

func wizardSpells(c *Character, cfg *DataGenerator.GeneratorConfig) {
	if c.Character.Class != "wizard" {
		return
	}

	if c.WizardData == nil {
		c.WizardData = &WizardData{}
	}

	if err := ApplyWizardSpellSlots(c, &cfg.SpellSlotTable); err != nil {
		panic(err)
	}

	c.WizardData.Spellbook = generateSpellbook(c, cfg)

	c.WizardData.MemorizedSpells = memorizeSpells(c)
}

func validateStatsVsClass(c *Character) {

	if c.Character.ClassConfig.MinStats.Strength > c.Attributes.STR.Score {
		c.Attributes.STR.Score = c.Character.ClassConfig.MinStats.Strength
	}
	if c.Character.ClassConfig.MinStats.Dexterity > c.Attributes.DEX {
		c.Attributes.DEX = c.Character.ClassConfig.MinStats.Dexterity
	}
	if c.Character.ClassConfig.MinStats.Constitution > c.Attributes.CON {
		c.Attributes.CON = c.Character.ClassConfig.MinStats.Constitution
	}
	if c.Character.ClassConfig.MinStats.Intelligence > c.Attributes.INT {
		c.Attributes.INT = c.Character.ClassConfig.MinStats.Intelligence
	}
	if c.Character.ClassConfig.MinStats.Wisdom > c.Attributes.WIS {
		c.Attributes.WIS = c.Character.ClassConfig.MinStats.Wisdom
	}
	if c.Character.ClassConfig.MinStats.Charisma > c.Attributes.CHA {
		c.Attributes.CHA = c.Character.ClassConfig.MinStats.Charisma
	}

}

func PickWeighted[T any](items []T, weightFn func(T) int) T {
	total := 0
	for _, item := range items {
		total += weightFn(item)
	}

	r := rand.Intn(total)

	for _, item := range items {
		r -= weightFn(item)
		if r < 0 {
			return item
		}
	}

	panic("unreachable")
}

var weightPattern = regexp.MustCompile(`[\d.]+`)

// parseWeight extracts a numeric weight (in lbs) from strings like "5 lb" or "0.5".
func parseWeight(weight string) int {
	match := weightPattern.FindString(weight)
	if match == "" {
		return 0
	}

	val, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0
	}

	return int(math.Round(val))
}

// gearItemFromEntry converts a pool entry into a character-owned GearItem.
func gearItemFromEntry(entry DataGenerator.GearPoolEntry, quantity int) GearItem {
	return GearItem{
		Name:     entry.Name,
		Weight:   parseWeight(entry.Weight),
		Quantity: quantity,
		// Properties: map[string]string{
		// 	"category": entry.Category,
		// 	"cost":     entry.Cost,
		// },
	}
}

// findGearEntry looks up a gear pool entry by name.
func findGearEntry(pool []DataGenerator.GearPoolEntry, name string) (DataGenerator.GearPoolEntry, error) {
	for _, e := range pool {
		if e.Name == name {
			return e, nil
		}
	}
	return DataGenerator.GearPoolEntry{}, fmt.Errorf("gear item %q not found in gear pool", name)
}

// filterAvailableGear returns pool entries not already added, and with a usable weight.
func filterAvailableGear(pool []DataGenerator.GearPoolEntry, added map[string]bool) []DataGenerator.GearPoolEntry {
	out := make([]DataGenerator.GearPoolEntry, 0, len(pool))
	for _, e := range pool {
		if added[e.Name] {
			continue
		}
		if e.ProbabilityWeight <= 0 {
			continue
		}
		out = append(out, e)
	}
	return out
}

// generateGear builds a character's starting gear from the class's gear
// assignment rule: guaranteed items plus a number of randomly weighted slots.
func generateGear(cfg *DataGenerator.GeneratorConfig, c *Character) []GearItem {
	rule, ok := cfg.GearAssignments.GearAssignments[c.Character.Class]
	if !ok {
		return nil
	}

	gear := make([]GearItem, 0, len(rule.Guaranteed)+rule.RandomSlots)
	added := make(map[string]bool)

	// Guaranteed items first
	for _, name := range rule.Guaranteed {
		entry, err := findGearEntry(cfg.GearPool.Pool, name)
		if err != nil {
			// Skip missing entries rather than failing the whole character
			continue
		}

		gear = append(gear, gearItemFromEntry(entry, 1))
		added[entry.Name] = true
	}

	// Random slots, weighted by ProbabilityWeight, no duplicates
	pool := filterAvailableGear(cfg.GearPool.Pool, added)

	for i := 0; i < rule.RandomSlots && len(pool) > 0; i++ {
		pick := PickWeighted(pool, func(g DataGenerator.GearPoolEntry) int {
			return g.ProbabilityWeight
		})

		gear = append(gear, gearItemFromEntry(pick, 1))
		added[pick.Name] = true

		pool = filterAvailableGear(pool, added)
	}

	return gear
}
