use turbogamf::{BoardActionTrigger, BoardBuilder};
use UnoCardColour::*;

fn main() {
    let mut board = BoardBuilder::new(UnoState::default())
        .add_cond_action(
            BoardActionTrigger::AfterTurn,
            |s| s.pickup.is_empty(),
            |s| s.reshuffle(),
        )
        .build();
}

struct UnoState {
    pickup: Vec<UnoCard>,
    putdown: Vec<UnoCard>,
}

impl Default for UnoState {
    fn default() -> Self {
        Self {
            pickup: UnoCard::full_deck(),
            putdown: Vec::new(),
        }
    }
}

impl UnoState {
    fn reshuffle(&mut self) {}
}

struct UnoPlayer {
    hand: Vec<UnoCard>,
}

#[derive(Copy, Clone)]
enum UnoCardColour {
    Red,
    Green,
    Blue,
    Yellow,
}

#[derive(Copy, Clone)]
enum UnoCard {
    Number { colour: UnoCardColour, number: u8 },
    Plus2 { colour: UnoCardColour },
    Reverse { colour: UnoCardColour },
    Skip { colour: UnoCardColour },
    Wild,
    WildPlus4,
}

impl UnoCard {
    /* For each rbgy:
     *      19 number cards (1 0 and 2 of 1-9)
     *      2 plus2
     *      2 Reverse
     *      2 skip
     * Also
     * 4 wild cards
     * 4 +4 wild cards
     */
    fn full_deck() -> Vec<UnoCard> {
        let mut deck = Vec::new();

        for colour in [Red, Green, Blue, Yellow] {
            deck.push(UnoCard::Number { colour, number: 0 });

            for _ in 0..2 {
                for number in 1..=9 {
                    deck.push(UnoCard::Number { colour, number });
                }

                deck.push(UnoCard::Plus2 { colour });
                deck.push(UnoCard::Reverse { colour });
                deck.push(UnoCard::Skip { colour });
            }
        }

        for _ in 0..4 {
            deck.push(UnoCard::Wild);
            deck.push(UnoCard::WildPlus4);
        }

        deck
    }
}
