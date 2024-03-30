type BoardActionConditionFn<S> = Option<Box<dyn Fn(&S) -> bool>>;
type BoardActionFn<S> = Box<dyn Fn(&mut S)>;

pub enum BoardActionTrigger {
    BeforeTurn,
    AfterTurn,
    BeforeRound,
    AfterRound,
}

pub struct BoardAction<S> {
    trigger: BoardActionTrigger,
    condition: BoardActionConditionFn<S>,
    action: BoardActionFn<S>,
}

pub struct Board<S> {
    state: S,
    actions: Vec<BoardAction<S>>,
}

pub struct BoardBuilder<S> {
    state: S,
    actions: Vec<BoardAction<S>>,
}

impl<S> BoardBuilder<S> {
    pub fn new(state: S) -> BoardBuilder<S> {
        BoardBuilder {
            state,
            actions: Vec::new(),
        }
    }

    pub fn add_action<F>(mut self, trigger: BoardActionTrigger, action: F) -> BoardBuilder<S>
    where
        F: Fn(&mut S) + 'static,
    {
        let action = BoardAction {
            trigger,
            condition: None,
            action: Box::new(action),
        };

        self.actions.push(action);
        self
    }

    pub fn add_cond_action<C, F>(
        mut self,
        trigger: BoardActionTrigger,
        cond: C,
        action: F,
    ) -> BoardBuilder<S>
    where
        C: Fn(&S) -> bool + 'static,
        F: Fn(&mut S) + 'static,
    {
        let action = BoardAction {
            trigger,
            condition: Some(Box::new(cond)),
            action: Box::new(action),
        };

        self.actions.push(action);
        self
    }

    pub fn build(self) -> Board<S> {
        Board {
            state: self.state,
            actions: self.actions,
        }
    }
}

trait Player {}

struct Game<S, P: Player> {
    board: Board<S>,
    players: Vec<P>,
}
