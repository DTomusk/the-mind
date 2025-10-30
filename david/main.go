package main

import (
	"fmt"
	"the-mind/cards"
)

func main() {
	hands := cards.GetHands(4, 10)
	for i, hand := range hands {
		fmt.Printf("Player %d's hand: %v\n", i+1, hand)
	}
}
