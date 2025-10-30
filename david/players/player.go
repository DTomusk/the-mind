package players

import (
	"fmt"
	"sync"
)

type Player struct {
	Id         int
	Hand       []int
	PlayChan   chan int
	NotifyChan chan struct{}
}

func (p *Player) Play(wg *sync.WaitGroup) {
	defer wg.Done()
	for len(p.Hand) > 0 {
		card := p.Hand[0]
		p.Hand = p.Hand[1:]
		fmt.Printf("Player %d is playing card %d\n", p.Id, card)
		p.PlayChan <- card
		<-p.NotifyChan
	}
	fmt.Printf("Player %d has no more cards to play.\n", p.Id)
}
