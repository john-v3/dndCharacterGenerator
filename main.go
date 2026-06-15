package main

import (
	"john-v3/dndCharacterGenerator/Character"
	dataexport "john-v3/dndCharacterGenerator/DataExport"
	"john-v3/dndCharacterGenerator/DataGenerator"
	"math"
	"path"
	"strconv"
)

var configPath string = "./rules/adnd2ed/"

func main() {

	ruleSet, err := DataGenerator.LoadGeneratorConfig(configPath)

	if err != nil {
		panic(err)
	}

	test := make([]character.Character, 100)

	meta := character.Metadata{
		Ruleset: configPath,
	}

	newMan := character.GenerateBaseCharacter(ruleSet)
	newMan.Character.Name = "character_1"
	character.ApplyDerivedStats(&newMan, ruleSet, 0)
	newMan.Metadata = meta

	saveDir := "./CharacterSheets/"

	for x := range 100 {
		newMan := character.GenerateBaseCharacter(ruleSet)

		name := "character_" + strconv.Itoa(x)
		newMan.Character.Name = name
		character.ApplyDerivedStats(&newMan, ruleSet, 0)
		newMan.Metadata = meta
		test[x] = newMan
		dataexport.PrintPDF(newMan, path.Join(saveDir, "character_"+strconv.Itoa(x)+".pdf"))
	}

}

func gaussianWeights(n int, mu, sigma float64) []float64 {
	w := make([]float64, n)
	var sum float64

	for i := 1; i <= n; i++ {
		x := float64(i)
		diff := x - mu
		w[i-1] = math.Exp(-(diff * diff) / (2 * sigma * sigma))
		sum += w[i-1]
	}

	for i := range w {
		w[i] /= sum
		w[i] *= 100
	}

	return w
}
