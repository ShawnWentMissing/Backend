package main

type Area int

const (
	Floor Area = iota
	FrontWithinBoundaryServe
	FrontWithinBoundary
	OutsideBoundary
)

type Player int

const (
	Player1 Player = 1
	Player2 Player = 2
)

type RallyState int

const (
	NewRally RallyState = iota
	NoBounces
	BouncedOffFloor
)

type SquashGame struct {
	Player1Score  int
	Player2Score  int
	MaxRoundScore int
	TurnPlayer    Player
	State         RallyState
}

type GameStorage struct {
	// Map to store games with raspberry pi's mac address
	games map[string]SquashGame
}

func NewGameStorage() *GameStorage {
	return &GameStorage{
		games: make(map[string]SquashGame),
	}
}

func (gs *GameStorage) AddGame(id string) {
	gs.games[id] = SquashGame{Player1Score: 0, Player2Score: 0, MaxRoundScore: 15, State: NewRally}
}

func (gs *GameStorage) GetGame(id string) (SquashGame, bool) {
	game, ok := gs.games[id]
	return game, ok
}

func (gs *GameStorage) UpdateGame(id string, game SquashGame) {
	gs.games[id] = game
}

func (gs *GameStorage) IncrementScore(id string, score int, player Player) (endGame, swapTurn, ok bool) {
	game, ok := gs.GetGame(id)
	if !ok {
		return endGame, swapTurn, false
	}

	switch player {
	case Player1:
		game.Player1Score++

		if game.Player1Score == game.MaxRoundScore {
			endGame = true
		}

		if player != game.TurnPlayer {
			game.TurnPlayer = Player1
			swapTurn = true
		}

	case Player2:
		game.Player2Score++

		if game.Player2Score == game.MaxRoundScore {
			endGame = true
		}

		if player != game.TurnPlayer {
			game.TurnPlayer = Player2
			swapTurn = true
		}
	}

	gs.UpdateGame(id, game)
	return endGame, swapTurn, true
}

func (gs *GameStorage) BallBounce(id string, hitArea Area) (endRally, ok bool) {
	game, ok := gs.GetGame(id)
	if !ok {
		return endRally, false
	}

	switch game.State {
	case NewRally:
		if hitArea == FrontWithinBoundaryServe {
			game.State = NoBounces
		} else {
			game.State = NewRally
		}
	case NoBounces:
		if hitArea == Floor {
			game.State = BouncedOffFloor
		} else if hitArea == FrontWithinBoundary {
			game.State = NoBounces
		}
	case BouncedOffFloor:
		if hitArea == Floor {
			game.State = NewRally
		} else if hitArea == FrontWithinBoundary {
			game.State = NoBounces
		}
	}

	if game.State == NewRally {
		endRally = true
	}

	gs.UpdateGame(id, game)
	return endRally, true
}
