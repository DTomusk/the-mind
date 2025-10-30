package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"the-mind/cards"
	"the-mind/players"
)

type GameConfig struct {
	NumPlayers     int
	CardsPerPlayer int
}

func main() {
	config, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}

	playChan := make(chan players.Move)
	cardsPlayedChan := make(chan []int)
	done := make(chan struct{})
	var wg sync.WaitGroup

	players := setUpGame(config, playChan, cardsPlayedChan)

	// Start each player's play routine in the background
	for _, player := range players {
		wg.Add(1)
		go player.Play(&wg, done)
	}

	// Run a goroutine in the background to listen for played cards
	// Gets cards from playChan and notifies the player to continue
	var gameWg sync.WaitGroup
	var cardsPlayed []int
	gameWg.Add(1)
	go func() {
		defer gameWg.Done()
		for move := range playChan {
			fmt.Printf("Main received card %d from Player %d\n", move.Card, move.PlayerId)

			// Send an empty struct to notify the player that a card has been played
			if len(cardsPlayed) > 0 {
				last := cardsPlayed[len(cardsPlayed)-1]
				if move.Card < last {
					fmt.Printf("Card %d played after %d! Game over.\n", move.Card, last)
					// Notify all players to stop playing
					cardsPlayedChan <- cardsPlayed
					close(done)
					return
				}
			}
			cardsPlayed = append(cardsPlayed, move.Card)
			fmt.Printf("Cards played so far: %v\n", cardsPlayed)
			select {
			case cardsPlayedChan <- cardsPlayed:
			case <-done:
				return
			}
		}
		// Once playChan is closed, we finish
		fmt.Println("No more cards to receive. Ending game.")
	}()

	// Wait for the player to finish, each player adds 1 to the wait group
	wg.Wait()
	// Once all players are done, close playChan to end the game goroutine
	close(playChan)
	gameWg.Wait()
	select {
	case <-done:
		fmt.Println("Game ended early due to invalid card.")
	default:
		fmt.Println("Game over — all cards played successfully.")
	}
}

func parseArgs() (*GameConfig, error) {
	// TODO: consider using flags
	args := os.Args[1:]
	if len(args) != 2 {
		return nil, fmt.Errorf("usage: go run main.go <numPlayers> <cardsPerPlayer>")

	}
	numPlayers, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, fmt.Errorf("error: numPlayers must be an integer")
	}
	cardsPerPlayer, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, fmt.Errorf("error: cardsPerPlayer must be an integer")
	}

	if numPlayers*cardsPerPlayer > 100 {
		return nil, fmt.Errorf("error: Not enough cards in the deck for the given number of players and cards per player")
	}
	return &GameConfig{
		NumPlayers:     numPlayers,
		CardsPerPlayer: cardsPerPlayer,
	}, nil
}

func setUpGame(config *GameConfig, playChan chan players.Move, cardsPlayedChan chan []int) []*players.Player {
	hands := cards.GetHands(config.NumPlayers, config.CardsPerPlayer)
	players := players.CreatePlayers(config.NumPlayers, hands, playChan, cardsPlayedChan)
	return players
}
