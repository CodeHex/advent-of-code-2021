package main

import (
	"adventofcode2021/pkg/fileparser"
	"adventofcode2021/pkg/maps"
	"fmt"
)

func main() {
	positionData := fileparser.ReadPairs[string, int]("day21/input.txt", ":")
	posPlayer1 := positionData[0].Value
	posPlayer2 := positionData[1].Value

	game := NewGame(posPlayer1, posPlayer2, NewDeterministicDie(), 1000)
	game.Play()
	fmt.Printf("[Part 1] %s wins (%d), %s loses with score %d after %d rolls, result:%d\n",
		game.FirstWin.Winner.Name,
		game.FirstWin.Winner.FinalScore,
		game.FirstWin.Loser.Name,
		game.FirstWin.Loser.FinalScore,
		game.FirstWin.RollCount,
		game.FirstWin.RollCount*game.FirstWin.Loser.FinalScore)

	game = NewGame(posPlayer1, posPlayer2, NewDiracDie(), 21)
	game.Play()
	fmt.Printf("[Part 2] Number of player 1 wins %d, number of player 2 wins %d\n",
		game.QuantumWins.NumP1Wins,
		game.QuantumWins.NumP2Wins)

}

type Die interface {
	RollScore() []int // Represents the values that the die can provide on any given turn
}

type DeterministicDie struct{ lastVal int }

func NewDeterministicDie() *DeterministicDie {
	return &DeterministicDie{lastVal: 1}
}

func (d *DeterministicDie) RollScore() []int {
	result := (d.lastVal + 1) * 3
	d.lastVal += 3
	if d.lastVal > 101 {
		d.lastVal -= 100
	}
	return []int{result}
}

type DiracDie struct{}

func NewDiracDie() *DiracDie {
	return &DiracDie{}
}

func (d *DiracDie) RollScore() []int {
	result := []int{}
	for r1 := 1; r1 <= 3; r1++ {
		for r2 := 1; r2 <= 3; r2++ {
			for r3 := 1; r3 <= 3; r3++ {
				result = append(result, r1+r2+r3)
			}
		}
	}
	return result
}

type Game struct {
	die           Die
	player1       *Player
	player2       *Player
	currentPlayer int
	turn          int
	maxScore      int
	FirstWin      struct {
		Winner    *Player
		Loser     *Player
		RollCount int
	}
	QuantumWins struct {
		NumP1Wins int64
		NumP2Wins int64
	}
}

type Player struct {
	Name       string
	States     map[int]map[PlayerGameState]int
	FinalScore int
}

type PlayerGameState struct {
	Pos   int
	Score int
}

func NewGame(p1StartPos, p2StartPos int, die Die, maxScore int) *Game {
	return &Game{
		die:           die,
		player1:       NewPlayer("player 1", p1StartPos),
		player2:       NewPlayer("player 2", p2StartPos),
		currentPlayer: 1,
		maxScore:      maxScore,
	}
}

func NewPlayer(name string, startPos int) *Player {
	states := make(map[int]map[PlayerGameState]int)
	// Add the starting position at turn 0
	states[0] = map[PlayerGameState]int{{Pos: startPos}: 1}
	return &Player{Name: name, States: states}
}

func (g *Game) CurrentPlayer() *Player {
	if g.currentPlayer == 1 {
		return g.player1
	}
	return g.player2
}

func (g *Game) SwitchPlayer() {
	if g.currentPlayer == 1 {
		g.currentPlayer = 2
	} else {
		g.currentPlayer = 1
	}
}

func (p *Player) IncrementState(turn int, state PlayerGameState, count int) {
	if _, ok := p.States[turn]; !ok {
		p.States[turn] = make(map[PlayerGameState]int)
	}
	p.States[turn][state] += count
}

func (g *Game) Play() {
	running := true

	// Evolve each universes independently for each player
	for running {
		g.turn++
		running = false                                        // If we progress any universe, we are still running
		for s, c := range g.CurrentPlayer().States[g.turn-1] { // Only progress the previous turns states
			if s.Score >= g.maxScore {
				// If score is a winner, then stop progressing this state
				continue
			}
			running = true
			for _, dieVal := range g.die.RollScore() {
				newPos := ModPos(s.Pos + dieVal)
				newState := PlayerGameState{Pos: newPos, Score: s.Score + newPos}
				g.CurrentPlayer().IncrementState(g.turn, newState, c)
			}
		}

		g.SwitchPlayer()

		for s, c := range g.CurrentPlayer().States[g.turn-1] {
			if s.Score < g.maxScore {
				g.CurrentPlayer().IncrementState(g.turn, s, c)
			}
		}
	}
	g.scoreGame()
}

func (g *Game) scoreGame() {
	// Check each turn for states that have winning scores. If a state has a winning score
	// then the number of universes in which that play wins is multiples by all the universes
	// in which the other player is still playing
	for i := 0; i < g.turn; i++ {
		for s, c := range g.player1.States[i] {
			if s.Score >= g.maxScore {
				g.QuantumWins.NumP1Wins += int64(c) * int64(maps.SumValues(g.player2.States[i]))
				if g.FirstWin.Winner == nil {
					g.FirstWin.Winner = g.player1
					g.FirstWin.Loser = g.player2
					g.FirstWin.Winner.FinalScore = s.Score
					g.FirstWin.Loser.FinalScore = maps.AnyKey(g.player2.States[i]).Score
					g.FirstWin.RollCount = i * 3
				}
			}
		}

		for s, c := range g.player2.States[i] {
			if s.Score >= g.maxScore {
				g.QuantumWins.NumP2Wins += int64(c) * int64(maps.SumValues(g.player1.States[i]))
				if g.FirstWin.Winner == nil {
					g.FirstWin.Winner = g.player2
					g.FirstWin.Loser = g.player1
					g.FirstWin.Winner.FinalScore = s.Score
					g.FirstWin.Loser.FinalScore = maps.AnyKey(g.player1.States[i]).Score
					g.FirstWin.RollCount = i * 3
				}
			}
		}
	}
}

func ModPos(val int) int {
	result := ((val - 1) % 10) + 1
	if result < 0 {
		panic("invalid")
	}
	return result
}
