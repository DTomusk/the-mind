package players

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Player struct {
	Id              int
	Hand            []int
	PlayChan        chan Move
	CardsPlayedChan chan []int
}

type Move struct {
	PlayerId int
	Card     int
}

func CreatePlayers(numPlayers int, hands [][]int, playChan chan Move, cardsPlayedChan chan []int) []*Player {
	players := make([]*Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = &Player{
			Id:              i + 1,
			Hand:            hands[i],
			PlayChan:        playChan,
			CardsPlayedChan: cardsPlayedChan,
		}
	}
	return players
}

func (p *Player) Play(wg *sync.WaitGroup, done <-chan struct{}) {
	sort.Ints(p.Hand)
	fmt.Printf("Player %d's sorted hand: %v\n", p.Id, p.Hand)
	// r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(p.Id)))
	defer wg.Done()
	for len(p.Hand) > 0 {
		card := p.Hand[0]
		fmt.Printf("Player %d is preparing to play card %d\n", p.Id, card)
		delay := time.Duration(card) * time.Second / 100
		timer := time.NewTimer(delay)

		// time.Sleep(time.Duration(card) * time.Second / 10)
		select {
		case <-timer.C:
			fmt.Printf("Player %d is playing card %d\n", p.Id, card)
			p.Hand = p.Hand[1:]
			p.PlayChan <- Move{PlayerId: p.Id, Card: card}

			select {
			case _, ok := <-p.CardsPlayedChan:
				if !ok {
					fmt.Printf("Player %d: CardsPlayedChan closed. Stopping play.\n", p.Id)
					return
				}
			case <-done:
				fmt.Printf("Player %d received done signal. Stopping.\n", p.Id)
				return
			}
		case <-p.CardsPlayedChan:
			// Received notification that a card has been played
			fmt.Printf("Player %d received notification of played card while waiting to play %d\n", p.Id, card)
			if !timer.Stop() {
				<-timer.C
			}
			continue
		case <-done:
			fmt.Printf("Player %d received done signal. Stopping play.\n", p.Id)
			if !timer.Stop() {
				<-timer.C
			}
			return
		}
	}
	fmt.Printf("Player %d has no more cards to play.\n", p.Id)
}
