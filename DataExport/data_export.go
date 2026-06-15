package dataexport

import (
	"fmt"
	character "john-v3/dndCharacterGenerator/Character"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func PrintPDFTest() {
	pdf := gofpdf.New("P", "mm", "Letter", "")

	Page1Test(pdf)
	pdf.OutputFileAndClose("test.pdf")
}

func Page1Test(pdf *gofpdf.Fpdf) {

	pdf.AddPage()

	pdf.Image(
		"./DataExport/charactersheets/page1_v2.jpg",
		0,
		0,
		215.9,
		279.4,
		false,
		"",
		0,
		"",
	)

	pdf.SetFont("Arial", "", 10)

	for x := 0.0; x < 220; x += 10 {
		for y := 0.0; y < 280; y += 5 {
			pdf.Text(x, y, fmt.Sprintf("%d,%d", int(x), int(y)))
		}
	}

	pdf.AddPage()
	pdf.Image(
		"./DataExport/charactersheets/page2.png",
		0,
		0,
		215.9,
		279.4,
		false,
		"",
		0,
		"",
	)

	pdf.SetFont("Arial", "", 10)

	for x := 0.0; x < 220; x += 10 {
		for y := 0.0; y < 280; y += 5 {
			pdf.Text(x, y, fmt.Sprintf("%d,%d", int(x), int(y)))
		}
	}

}

func PrintPDF(c character.Character, savePath string) {
	pdf := gofpdf.New("P", "mm", "Letter", "")

	Page1(pdf, c)
	// Page1Test(pdf)
	Page2(pdf, c)

	pdf.OutputFileAndClose(savePath)
}

func Page1(pdf *gofpdf.Fpdf, c character.Character) {

	pdf.AddPage()

	pdf.Image(
		"./DataExport/charactersheets/page1_v2.jpg",
		0,
		0,
		215.9,
		279.4,
		false,
		"",
		0,
		"",
	)

	proficenciesAndLanguageArea := concat([]string{"Languages: "})
	for x := range c.Languages {
		proficenciesAndLanguageArea = concat([]string{proficenciesAndLanguageArea, c.Languages[x], ","})
	}

	proficenciesAndLanguageArea += "\n"

	proficenciesAndLanguageArea += concat([]string{"weapon spec. \n"})
	for x := range c.WeaponSpecializations {
		proficenciesAndLanguageArea += concat([]string{c.WeaponSpecializations[x].Name, " ", c.WeaponSpecializations[x].Rank, "\n"})
	}
	proficenciesAndLanguageArea += "\n"

	proficenciesAndLanguageArea += concat([]string{"nonweapon spec. \n"})
	for x := range c.NonWeaponSpecializations {
		proficenciesAndLanguageArea += concat([]string{c.NonWeaponSpecializations[x].Name, " ", c.NonWeaponSpecializations[x].BaseStat, " roll:", strconv.Itoa(c.NonWeaponSpecializations[x].TargetRoll), "\n"})
	}

	specialAbilities := concat([]string{"Spell Immunities INT: ", concat(c.DerivedIntelligenceStats.SpellImmunities), "\n"})
	specialAbilities += concat([]string{"Spell Immunities WIS: ", concat(c.DerivedWisdomStats.SpellImmunities), "\n"})

	if c.WizardData != nil && len(c.WizardData.MemorizedSpells) != 0 {
		specialAbilities += "memorized spells\n"

		for x := range c.WizardData.MemorizedSpells {
			specialAbilities = concat([]string{specialAbilities, c.WizardData.MemorizedSpells[x], "\n"})
		}
	}

	if c.ThiefData != nil && c.ThiefData.ClimbWalls != 0 {
		specialAbilities += "PickPocket: " + strconv.Itoa(c.ThiefData.PickPocket) + "\n"
		specialAbilities += "OpenLock: " + strconv.Itoa(c.ThiefData.OpenLock) + "\n"
		specialAbilities += "FindRemoveTrap: " + strconv.Itoa(c.ThiefData.FindRemoveTrap) + "\n"
		specialAbilities += "MoveSilently: " + strconv.Itoa(c.ThiefData.MoveSilently) + "\n"
		specialAbilities += "HideInShadows: " + strconv.Itoa(c.ThiefData.HideInShadows) + "\n"
		specialAbilities += "DetectNoise: " + strconv.Itoa(c.ThiefData.DetectNoise) + "\n"
		specialAbilities += "ClimbWalls: " + strconv.Itoa(c.ThiefData.ClimbWalls) + "\n"
		specialAbilities += "ReadLanguages: " + strconv.Itoa(c.ThiefData.ReadLanguages) + "\n"
	}

	pdf.SetFont("Arial", "", 10)

	DrawField(pdf, PositionsPage1["CharacterName"], c.Character.Name)
	DrawField(pdf, PositionsPage1["alignment"], c.Character.Alignment)
	DrawField(pdf, PositionsPage1["race"], c.Character.Race)
	DrawField(pdf, PositionsPage1["class"], c.Character.Class)
	DrawField(pdf, PositionsPage1["level"], strconv.Itoa(c.Character.Level))
	DrawField(pdf, PositionsPage1["playerName"], "Your Name Here")
	DrawField(pdf, PositionsPage1["family"], c.Biography.Family)
	DrawField(pdf, PositionsPage1["RaceClan"], c.Biography.RaceClan)
	DrawField(pdf, PositionsPage1["HomeLand"], c.Biography.Homeland)
	DrawField(pdf, PositionsPage1["leigePatron"], c.Biography.Patron)
	DrawField(pdf, PositionsPage1["religon"], c.Biography.Religion)
	DrawField(pdf, PositionsPage1["sex"], "male")
	DrawField(pdf, PositionsPage1["age"], strconv.Itoa(c.Physical.Age))
	DrawField(pdf, PositionsPage1["socialClass"], c.Biography.SocialClass)
	DrawField(pdf, PositionsPage1["status"], c.Biography.Status)
	DrawField(pdf, PositionsPage1["height"], strconv.Itoa(c.Physical.Height))
	DrawField(pdf, PositionsPage1["weight"], strconv.Itoa(c.Physical.Weight))
	DrawField(pdf, PositionsPage1["birthRank"], strconv.Itoa(c.Biography.BirthRank))
	DrawField(pdf, PositionsPage1["numberOfSiblings"], strconv.Itoa(c.Biography.NumberOfSiblings))
	DrawField(pdf, PositionsPage1["hair"], c.Physical.Hair)
	DrawField(pdf, PositionsPage1["eyes"], c.Physical.Eyes)
	DrawField(pdf, PositionsPage1["appearance"], c.Physical.Appearance)
	DrawField(pdf, PositionsPage1["honor"], strconv.Itoa(c.Biography.Honor))
	DrawField(pdf, PositionsPage1["BaseHonor"], strconv.Itoa(c.Biography.BaseHonor))
	DrawField(pdf, PositionsPage1["ReactionAdjustment1"], strconv.Itoa(c.DerivedDexterityStats.ReactionAdjustment))
	DrawField(pdf, PositionsPage1["STR"], strconv.Itoa(c.Attributes.STR.Score))
	DrawField(pdf, PositionsPage1["DEX"], strconv.Itoa(c.Attributes.DEX))
	DrawField(pdf, PositionsPage1["CON"], strconv.Itoa(c.Attributes.CON))
	DrawField(pdf, PositionsPage1["INT"], strconv.Itoa(c.Attributes.INT))
	DrawField(pdf, PositionsPage1["WIS"], strconv.Itoa(c.Attributes.WIS))
	DrawField(pdf, PositionsPage1["CHR"], strconv.Itoa(c.Attributes.CHA))
	DrawField(pdf, PositionsPage1["ProficenciesAndLanguageArea"], proficenciesAndLanguageArea)
	DrawField(pdf, PositionsPage1["hitprob"], strconv.Itoa(c.DerivedStrengthStats.HitProbability))
	DrawField(pdf, PositionsPage1["dmgadj"], strconv.Itoa(c.DerivedStrengthStats.DamageAdjustment))
	DrawField(pdf, PositionsPage1["wgtallow"], strconv.Itoa(c.DerivedStrengthStats.WeightAllowance))
	DrawField(pdf, PositionsPage1["maxpress"], strconv.Itoa(c.DerivedStrengthStats.MaxPress))
	DrawField(pdf, PositionsPage1["opdrs"], strconv.Itoa(c.DerivedStrengthStats.BendBarsLiftGates))
	DrawField(pdf, PositionsPage1["bblg"], strconv.Itoa(c.DerivedStrengthStats.BendBarsLiftGates))
	DrawField(pdf, PositionsPage1["rctnadj"], strconv.Itoa(c.DerivedDexterityStats.ReactionAdjustment))
	DrawField(pdf, PositionsPage1["missleAttackAdj"], strconv.Itoa(c.DerivedDexterityStats.MissileAttackAdjustment))
	DrawField(pdf, PositionsPage1["defenseAdj"], strconv.Itoa(c.DerivedDexterityStats.DefenseAdjustment))
	DrawField(pdf, PositionsPage1["hpAdj"], strconv.Itoa(c.DerivedConstitutionStats.HPAdjustment))
	DrawField(pdf, PositionsPage1["sysShk"], strconv.Itoa(c.DerivedConstitutionStats.SystemShock))
	DrawField(pdf, PositionsPage1["resSur"], strconv.Itoa(c.DerivedConstitutionStats.Resurrection))
	DrawField(pdf, PositionsPage1["PoisonSave"], strconv.Itoa(c.DerivedConstitutionStats.PoisonSave))
	DrawField(pdf, PositionsPage1["Regen"], strconv.Itoa(c.DerivedConstitutionStats.Regeneration))
	DrawField(pdf, PositionsPage1["noOfLang"], strconv.Itoa(c.DerivedIntelligenceStats.NumberOfLanguages))
	DrawField(pdf, PositionsPage1["SpellLvl"], strconv.Itoa(c.DerivedIntelligenceStats.SpellLevel))
	DrawField(pdf, PositionsPage1["LrnSp"], strconv.Itoa(c.DerivedIntelligenceStats.LearnSpell))
	DrawField(pdf, PositionsPage1["SpellsPerLevel"], strconv.Itoa(c.DerivedIntelligenceStats.SpellsPerLevel))
	DrawField(pdf, PositionsPage1["SpellImmunity"], strconv.Itoa(len(c.DerivedIntelligenceStats.SpellImmunities)))
	DrawField(pdf, PositionsPage1["MagDefAdjus"], strconv.Itoa(c.DerivedWisdomStats.MagicDefenseAdjustment))
	DrawField(pdf, PositionsPage1["BonusSpells"], strconv.Itoa(len(c.DerivedWisdomStats.BonusSpells)))
	DrawField(pdf, PositionsPage1["SpellFail"], strconv.Itoa(c.DerivedWisdomStats.SpellFailure))
	DrawField(pdf, PositionsPage1["MaxNoHench"], strconv.Itoa(c.DerivedCharismaStats.MaxNumberOfHenchmen))
	DrawField(pdf, PositionsPage1["LoyaltyBase"], strconv.Itoa(c.DerivedCharismaStats.LoyaltyBase))
	DrawField(pdf, PositionsPage1["ReactionAdjustment2"], strconv.Itoa(c.DerivedDexterityStats.ReactionAdjustment))
	DrawField(pdf, PositionsPage1["ParalyzePoisonSave"], strconv.Itoa(c.SavingThrows.ParalyzePoisonDeath))
	DrawField(pdf, PositionsPage1["rodStaffWandSave"], strconv.Itoa(c.SavingThrows.RodStaffOrWand))
	DrawField(pdf, PositionsPage1["PetrifyPolymorphSave"], strconv.Itoa(c.SavingThrows.PetrifyPolymorph))
	DrawField(pdf, PositionsPage1["BreathWeaponSave"], strconv.Itoa(c.SavingThrows.BreathWeapon))
	DrawField(pdf, PositionsPage1["SpellSave"], strconv.Itoa(c.SavingThrows.Spell))
	DrawField(pdf, PositionsPage1["SpecialAbilities"], specialAbilities)
	DrawField(pdf, PositionsPage1["hitpoints"], strconv.Itoa(c.DerivedConstitutionStats.Hitpoints))

	DrawField(pdf, PositionsPage1["Thaco1"], strconv.Itoa(c.Combat.THACO))
	DrawField(pdf, PositionsPage1["Thaco2"], strconv.Itoa(c.Combat.THACO))
	DrawField(pdf, PositionsPage1["Thaco3"], strconv.Itoa(c.Combat.THACO))
	DrawField(pdf, PositionsPage1["Thaco4"], strconv.Itoa(c.Combat.THACO))
	DrawField(pdf, PositionsPage1["Thaco5"], strconv.Itoa(c.Combat.THACO))
	DrawField(pdf, PositionsPage1["Thaco6"], strconv.Itoa(c.Combat.THACO))
}

func Page2(pdf *gofpdf.Fpdf, c character.Character) {
	pdf.AddPage()
	pdf.Image(
		"./DataExport/charactersheets/page2.png",
		0,
		0,
		215.9,
		279.4,
		false,
		"",
		0,
		"",
	)

	pdf.SetFont("Arial", "", 10)

	gear := "Gold: " + strconv.Itoa(c.Economy.Gold) + "\n"

	fmt.Println(c.Gear)
	if len(c.Gear) != 0 {
		for x := range c.Gear {
			gear = concat([]string{gear, c.Gear[x].Name, " quantity ", strconv.Itoa(c.Gear[x].Quantity), "wght: ", strconv.Itoa(c.Gear[x].Weight)})
			// for _, d := range c.Gear[x].Properties {
			// 	gear = concat([]string{gear, d,"\n"})
			// }

			if len(c.Gear[x].Properties) == 0 {
				gear += "\n"
			}
		}
	}

	DrawField(pdf, PositionsPage2["gear"], gear)

}

func DrawField(pdf *gofpdf.Fpdf, pos Position, value string) {
	lines := strings.Split(value, "\n")
	lineHeight := 4.0 // adjust to taste, in mm
	for i, line := range lines {
		pdf.Text(pos.X, pos.Y+float64(i)*lineHeight, line)
	}
}

func concat(values []string) string {
	var sb strings.Builder
	for _, value := range values {
		sb.WriteString(value) // Appends each string to the internal buffer
	}
	return sb.String() // Converts the entire buffer to a single string
}
