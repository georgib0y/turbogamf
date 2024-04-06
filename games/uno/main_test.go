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
	d := MockUnoDeck()

	c := d.Pop()

	if len(d) != 3 {
		t.Fatal("Deck length is not 3")
	}

	if c.c != YELLOW {
		t.Fatal("Card colour is not YELLOW")
	}
}

func TestUnoDeckSuffles(t *testing.T) {
	d1, d2 := MockUnoDeck(), MockUnoDeck()

	d1.Shuffle()

	for i := range d1 {
		if d1[i].c != d2[i].c {
			return
		}
	}

	t.Fatal("D1 and D2 are the same order")
}
