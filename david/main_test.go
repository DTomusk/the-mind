package main

import (
	"testing"
)

// TestDeckIntegrity verifies that the deck has 100 unique cards numbered 1–100.
func TestDeckIntegrity(t *testing.T) {
	deck := createDeck()
	shuffleDeck(deck)

	if len(deck) != 100 {
		t.Fatalf("expected deck length 100, got %d", len(deck))
	}

	seen := make(map[int]bool)

	for _, card := range deck {
		if card < 1 || card > 100 {
			t.Fatalf("card %d out of range (must be 1–100)", card)
		}
		if seen[card] {
			t.Fatalf("duplicate card found: %d", card)
		}
		seen[card] = true
	}
}
