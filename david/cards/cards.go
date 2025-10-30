package cards

import (
	"fmt"
	"math/rand"
	"time"
)

func GetHands(numPlayers int, cardsPerPlayer int) [][]int {
	deck := createDeck()
	fmt.Println("Deck of cards:", deck)
	shuffleDeck(deck)
	fmt.Println("Shuffled deck of cards:", deck)
	return dealCards(deck, numPlayers, cardsPerPlayer)
}

// initialize a new deck of cards numbered 1 to 100
func createDeck() []int {
	fmt.Println("Creating a new deck of cards...")
	deck := make([]int, 100)
	for i := 0; i < 100; i++ {
		deck[i] = i + 1
	}
	return deck
}

func shuffleDeck(deck []int) []int {
	fmt.Println("Shuffling the deck...")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}

func dealCards(deck []int, numPlayers int, cardsPerPlayer int) [][]int {
	fmt.Printf("Dealing %d cards to %d players...\n", cardsPerPlayer, numPlayers)
	hands := make([][]int, numPlayers)
	for i := 0; i < numPlayers; i++ {
		hands[i] = make([]int, cardsPerPlayer)
		for j := 0; j < cardsPerPlayer; j++ {
			hands[i][j] = deck[i*cardsPerPlayer+j]
		}
	}
	return hands
}
