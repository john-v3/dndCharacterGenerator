package utilities

import "math/rand"

type DiceConfig struct {
	DiceCount  int `yaml:"diceCount"`
	DiceSides  int `yaml:"diceSides"`
	Multiplier int `yaml:"multiplier,omitempty"`
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
