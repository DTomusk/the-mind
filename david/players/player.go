package players

import (
	"fmt"
	"sync"
)

type Player struct {
	Id              int
	Hand            []int
	PlayChan        chan int
	CardsPlayedChan chan []int
}

func CreatePlayers(numPlayers int, hands [][]int, playChan chan int, cardsPlayedChan chan []int) []*Player {
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

func (p *Player) Play(wg *sync.WaitGroup) {
	defer wg.Done()
	for len(p.Hand) > 0 {
		card := p.Hand[0]
		p.Hand = p.Hand[1:]
		fmt.Printf("Player %d is playing card %d\n", p.Id, card)
		p.PlayChan <- card
		<-p.CardsPlayedChan
	}
	fmt.Printf("Player %d has no more cards to play.\n", p.Id)
}
