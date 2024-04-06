package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	numPlayers := 2

	g := NewUnoGame(numPlayers)

	go g.Run()

	winner := <-g.winner
	log.Printf("Winner is player %d\n", winner)
}

type colour int

const (
	RED colour = iota
	GREEN
	BLUE
	YELLOW
	WILD
)

var colourStrings map[colour]string = map[colour]string{
	RED:    "RED",
	GREEN:  "GREEN",
	BLUE:   "BLUE",
	YELLOW: "YELLOW",
}

func (c colour) String() string {
	s, ok := colourStrings[c]
	if !ok {
		panic(fmt.Sprintf("Unknown colour with value %d", int(c)))
	}
	return s
}

type ctype int

const (
	PLUS2 ctype = iota + 10
	REVERSE
	SKIP
	PLUS4
	CHANGE_COLOUR
)

var ctypeStrings map[ctype]string = map[ctype]string{
	PLUS2:         "PLUS2",
	REVERSE:       "REVERSE",
	SKIP:          "SKIP",
	PLUS4:         "PLUS4",
	CHANGE_COLOUR: "CHANGE_COLOUR",
}

func (t ctype) String() string {
	if int(t) < 10 {
		return fmt.Sprint(int(t))
	}

	s, ok := ctypeStrings[t]
	if !ok {
		panic(fmt.Sprintf("Unknown card type with value: %d", int(t)))
	}
	return s
}

type UnoCard struct {
	c colour
	t ctype
}

func (c UnoCard) CanPlaceOn(o UnoCard) bool {
	return c.c == WILD || c.c == o.c || c.t == o.t
}

type UnoDeck []UnoCard

/* For each rbgy:
 *      19 number cards (1 0 and 2 of 1-9)
 *      2 plus2
 *      2 Reverse
 *      2 skip
 * Also
 * 4 wild cards
 * 4 +4 wild cards
 */
func FullDeck() UnoDeck {
	d := UnoDeck{}

	for _, col := range []colour{RED, GREEN, BLUE, YELLOW} {
		d = append(d, UnoCard{c: col, t: ctype(0)})

		for i := 0; i < 2; i++ {
			for v := 1; v < 10; v++ {
				d = append(d, UnoCard{c: col, t: ctype(v)})
			}

			d = append(d, UnoCard{c: col, t: PLUS2})
			d = append(d, UnoCard{c: col, t: REVERSE})
			d = append(d, UnoCard{c: col, t: SKIP})
		}
	}

	for i := 0; i < 4; i++ {
		d = append(d, UnoCard{c: WILD, t: PLUS4})
		d = append(d, UnoCard{c: WILD, t: CHANGE_COLOUR})
	}

	return d
}

func (d *UnoDeck) Push(c UnoCard) {
	(*d) = append(*d, c)
}

func (d *UnoDeck) Pop() UnoCard {
	c := (*d)[len((*d))-1]
	*d = (*d)[:len((*d))-1]

	return c
}

func (d *UnoDeck) PopAt(idx int) UnoCard {
	c := (*d)[idx]
	*d = append((*d)[:idx], (*d)[idx+1:]...)

	return c
}

func (d UnoDeck) Top() UnoCard {
	return d[len(d)-1]
}

func (d *UnoDeck) Shuffle() {
	rand.Shuffle(len((*d)), func(i, j int) {
		(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
	})
}

func PromptCardChoice(options map[int]UnoCard) (int, error) {
	for i, c := range options {
		fmt.Printf("[%d]: %s\t%s\n", i, c.c, c.t)
	}

	for {
		fmt.Printf("Please choose an option [0 - %d]: ", len(options)-1)

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return 0, err
		}

		idx, err := strconv.Atoi(input)

		if err != nil || idx < 0 || idx >= len(options) {
			continue
		}

		return idx, nil
	}
}

type UnoGame struct {
	pickup, putdown UnoDeck
	players         []*UnoPlayer
	winner          chan int
	reverse         bool
}

func NewUnoGame(n int) *UnoGame {
	if n > 10 {
		panic("Too many players")
	}

	p := []*UnoPlayer{}
	for i := 0; i < n; i++ {
		p = append(p, NewUnoPlayer())
	}

	g := &UnoGame{
		pickup:  FullDeck(),
		putdown: UnoDeck{},
		players: p,
		winner:  make(chan int),
		reverse: false,
	}

	// start each player off with 7 cards each
	for i := range g.players {
		for j := 0; j < 7; j++ {
			g.Deal(i)
		}
	}

	// put the first non wild card down
	first := g.pickup.Pop()
	for first.c == WILD {
		first = g.pickup.Pop()
	}

	g.putdown.Push(first)

	return g
}

func (g *UnoGame) nextPlayerIdx(curr int) int {
	if g.reverse {
		curr -= 1
	} else {
		curr += 1
	}

	if curr < 0 {
		curr = len(g.players) - 1
	} else if curr >= len(g.players) {
		curr = 0
	}

	return curr
}

func (g *UnoGame) Run() {
	pIdx := 0

	for {
		card, ok := g.ChooseCard(pIdx)

		if !ok {
			g.Deal(pIdx)
			pIdx = g.nextPlayerIdx(pIdx)
			continue
		}

		g.putdown.Push(card)

		switch g.putdown.Top().t {
		case REVERSE:
			g.reverse = !g.reverse
		case SKIP:
			pIdx = g.nextPlayerIdx(pIdx)
		}

		pIdx = g.nextPlayerIdx(pIdx)
	}
}

func (g *UnoGame) ChooseCard(pIdx int) (UnoCard, bool) {
	t := g.putdown.Top()
	options := map[int]UnoCard{}

	for i, c := range g.players[pIdx].hand {
		if c.CanPlaceOn(t) {
			options[i] = c
		}
	}

	if len(options) == 0 {
		return UnoCard{}, false
	}

	idx, err := PromptCardChoice(options)
	if err != nil {
		panic(err)
	}

	return g.players[pIdx].hand.PopAt(idx), true
}

func (g *UnoGame) Pickup() UnoCard {
	if len(g.pickup) == 0 {
		g.resortDecks()
	}

	return g.pickup.Pop()
}

func (g *UnoGame) resortDecks() {
	// put all the cards but the top one from putdown deck into the pickup deck
	g.pickup = append(g.pickup, g.putdown[:len(g.putdown)-1]...)
	g.pickup.Shuffle()

	g.putdown = UnoDeck{g.putdown[len(g.putdown)-1]}
}

func (g *UnoGame) Deal(p int) {
	g.players[p].AddToHand(g.pickup.Pop())
}

type UnoPlayer struct {
	hand UnoDeck
}

func NewUnoPlayer() *UnoPlayer {
	return &UnoPlayer{hand: UnoDeck{}}
}

func (p *UnoPlayer) AddToHand(c UnoCard) {
	p.hand = append(p.hand, c)
}
