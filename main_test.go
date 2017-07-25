package main

import (
	"testing"
)

func TestMove(t *testing.T) {
	setupMaps(false)

	k := cardsToKey(initialCards)

	k.Log()

	k.move(cardToVal["10C"], 4)
	if k[1] != cardToVal["10C"] {
		t.Error("wrong pos")
	}

	k.Log()

}

func TestMove2(t *testing.T) {
	setupMaps(false)

	k := cardsToKey(`
7H 8H 9H 10H JH QH KH 3C 4C 5C 6C 7C
8C 9C 10C JC QC KC 2S 4S 3H 9S *A 10S JS
QS KS 1D 2D 3D 4D 5D 6D 7D 8D 9D 10D JD
QD KD 1C 2C 1H 3S 1S 2H 5H 5S 6S 7S 8S
4H *B 6H 
`)

	k.Log()

	k.move(ja, 1)

	k.move(jb, 2)
	if k[1] != jb {
		t.Error("wrong pos")
	}

	k.Log()

}

func TestMove3(t *testing.T) {
	setupMaps(false)

	k := cardsToKey(`
7H 8H 9H 10H JH QH KH 3C 4C 5C 6C 7C 8C
9C 10C JC QC KC 2S 4S 3H 9S *A 10S JS QS
KS 1D 2D 3D 4D 5D 6D 7D 8D 9D 10D JD QD
KD 1C 2C 1H 3S 1S 2H 5H 5S 6S 7S 8S 4H
*B 6H
`)

	k.Log()

	k.move(ja, 1)
	if k[23] != ja {
		t.Error("wrong pos")
	}

	k.Log()
}

func BenchmarkStream(b *testing.B) {
	setupMaps(false)

	k := cardsToKey(`
7H 8H 9H 10H JH QH KH 3C 4C 5C 6C 7C
8C 9C 10C JC QC KC 2S 4S 3H 9S *A 10S JS
QS KS 1D 2D 3D 4D 5D 6D 7D 8D 9D 10D JD
QD KD 1C 2C 1H 3S 1S 2H 5H 5S 6S 7S 8S
4H *B 6H 
`)
	// 354595 per 1000
	// 354 ns per round
	for i := 0; i < b.N; i++ {
		k.stream(100)
	}
}
