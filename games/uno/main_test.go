package main

import (
	"testing"
)

func MockUnoDeck() UnoDeck {
	return UnoDeck{
		UnoCard{c: RED, t: ctype(0)},
		UnoCard{c: GREEN, t: ctype(1)},
		UnoCard{c: BLUE, t: ctype(2)},
		UnoCard{c: YELLOW, t: ctype(3)},
	}
}

func TestUnoDeckPopsCard(t *testing.T) {
	d := FullDeck()
	l := len(d)

	expected := d[l-1]
	c := d.Pop()

	if len(d) != l-1 {
		t.Fatal("Deck length did not shrink")
	}

	if c.c != expected.c {
		t.Fatal("Card colour does not match expected")
	}
}

func TestUnoDeckShuffles(t *testing.T) {
	d1, d2 := FullDeck(), FullDeck()

	d1.Shuffle()

	for i := range d1 {
		if d1[i].c != d2[i].c {
			return
		}
	}

	t.Fatal("D1 and D2 are the same order")
}

func TestUnoDeckResorts(t *testing.T) {
	g := UnoGame{
		pickup:  UnoDeck{},
		putdown: FullDeck(),
	}

	top := g.putdown.Top()
	l := len(g.putdown)

	g.resortDecks()

	if top != g.putdown.Top() {
		t.Errorf("Tops are different: %v != %v", top, g.putdown.Top())
	}

	if len(g.pickup) != l-1 {
		t.Errorf("Pickup pile unexpected length: expected %d got %d", l-1, len(g.pickup))
	}
}

func TestNextPlayerIndex(t *testing.T) {
	g := UnoGame{
		players: []*UnoPlayer{nil, nil, nil},
		reverse: false,
	}

	curr := 0

	curr = g.nextPlayerIdx(curr)
	curr = g.nextPlayerIdx(curr)
	curr = g.nextPlayerIdx(curr)

	if curr != 0 {
		t.Fatalf("Current player index is %d, expected 0", curr)
	}

	g.reverse = true

	curr = g.nextPlayerIdx(curr)
	curr = g.nextPlayerIdx(curr)

	if curr != 1 {
		t.Fatalf("Current player index is %d, expected 1", curr)
	}
}
