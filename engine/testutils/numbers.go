package testutils

import "math/rand/v2"

func NumbersAscending(max uint64) []uint64 {
	seed := []uint64{}
	for i := uint64(1); i <= max; i++ {
		seed = append(seed, i)
	}
	return seed
}

func NumbersDescending(max uint64) []uint64 {
	seed := []uint64{}
	for i := max; i >= 1; i-- {
		seed = append(seed, i)
	}
	return seed
}

func NumbersRandom(max uint64) []uint64 {
	seed := []uint64{}
	for i := uint64(1); i <= max; i++ {
		seed = append(seed, i)
	}

	rand.Shuffle(len(seed), func(i, j int) {
		seed[i], seed[j] = seed[j], seed[i]
	})
	return seed
}
