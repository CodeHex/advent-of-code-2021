package main

import (
	"adventofcode2021/pkg/maps"
	"adventofcode2021/pkg/slices"
	"fmt"
)

// Sample
//
// #############
// #...........#
// ###B#C#B#D###
//   #A#D#C#A#
//   #########

func main() {

	//gamePart1 := NewGame("B", "A", "C", "D", "B", "C", "D", "A")
	gamePart1 := NewGame("B", "D", "B", "A", "C", "A", "D", "C")
	gamePart1.RunLowest()
	fmt.Println("[Part 1] Lowest score is", gamePart1.LowestScore)

	//gamePart2 := NewGame("B", "D", "D", "A", "C", "C", "B", "D", "B", "B", "A", "C", "D", "A", "C", "A")
	gamePart2 := NewGame("B", "D", "D", "D", "B", "C", "B", "A", "C", "B", "A", "A", "D", "A", "C", "C")
	gamePart2.RunLowest()
	fmt.Println("[Part 2] Lowest score is", gamePart2.LowestScore)
}

type Game struct {
	Pods        []Pod
	HomeTiles   map[string][]Tile
	Moves       map[Tile][]*Move
	LowestScore int
	RunningList *GameState
	Played      map[string]int
}

type GameState struct {
	PodPositions  map[Pod]Tile
	TilePositions map[Tile]Pod
	hash          string
	Score         int
	Finished      bool
	PrevMoves     []string
	Next          *GameState
}

type Move struct {
	Steps      int
	EndTile    Tile
	BlockTiles []Tile
}

type Tile string
type Pod string

func (p Pod) PodType() string {
	return string(p[1])
}

func (t Tile) TileType() (string, string) {
	r1 := string(t[0])
	r2 := ""
	if r1 == "b" {
		r2 = string(t[1])
	}
	return r1, r2
}

func (p Pod) ScoreMultiplier() int {
	switch p.PodType() {
	case "A":
		return 1
	case "B":
		return 10
	case "C":
		return 100
	case "D":
		return 1000
	default:
		panic("unrecognized pod type")
	}
}

func NewGame(inPodTypes ...string) *Game {
	columnLength := len(inPodTypes) / 4
	podTypes := []string{"A", "B", "C", "D"}

	startingTiles := []Tile{}
	pods := []Pod{}
	homeTiles := make(map[string][]Tile)
	// Generate starting tiles in order of top left to bottom right, reading down then across
	for _, t := range podTypes {
		homes := []Tile{}
		for i := 0; i < columnLength; i++ {
			tile := Tile(fmt.Sprintf("b%s%d", t, i))
			homes = append(homes, tile)
			startingTiles = append(startingTiles, tile)
			pods = append(pods, Pod(fmt.Sprintf("p%s%d", t, i)))
		}
		homeTiles[t] = homes
	}

	podCounts := make(map[string]int)
	pos := make(map[Pod]Tile)
	tiles := make(map[Tile]Pod)
	for i := range startingTiles {
		pod := Pod(fmt.Sprintf("p%s%d", inPodTypes[i], podCounts[inPodTypes[i]]))
		podCounts[inPodTypes[i]]++
		pos[pod] = startingTiles[i]
		tiles[startingTiles[i]] = pod

	}
	startState := &GameState{
		PodPositions:  pos,
		TilePositions: tiles,
	}
	g := &Game{
		Pods:        pods,
		HomeTiles:   homeTiles,
		RunningList: startState,
		Played:      make(map[string]int),
	}
	g.Moves = g.InitMoves()
	return g
}

func (s *GameState) Clone() *GameState {
	newState := &GameState{
		PodPositions:  make(map[Pod]Tile),
		TilePositions: make(map[Tile]Pod),
		Score:         s.Score,
		Finished:      s.Finished,
	}

	for k, v := range s.PodPositions {
		newState.PodPositions[k] = v
	}
	for k, v := range s.TilePositions {
		newState.TilePositions[k] = v
	}
	for _, v := range s.PrevMoves {
		newState.PrevMoves = append(newState.PrevMoves, v)
	}
	return newState
}

func (s *GameState) Hash(g *Game) string {
	if s.hash != "" {
		return s.hash
	}

	out := "|"
	for _, p := range g.Pods {
		out += fmt.Sprintf("%s:%s|", p, s.PodPositions[p])
	}
	s.hash = out
	return s.hash
}

func (g *Game) RunLowest() {
	// Create initial game state
	count := 0
	for g.RunningList != nil {
		count++
		// Take the lowest score of current running states, removing it from the list
		state := g.RunningList
		g.RunningList = g.RunningList.Next
		if count%100000 == 0 {
			fmt.Printf("%d running...(score: %d)\n", g.LengthRunning(), state.Score)
		}

		// Check if we've already played this but with a lower or equal score, if so
		// ignore it
		//playedScore, ok := g.Played[state.Hash]
		//if ok && playedScore <= state.Score {
		//	continue
		//}
		//

		// If this state has finished then this must be the lowest as all other states
		// have a higher score and maybe still running
		if state.Finished {
			g.LowestScore = state.Score
			return
		}

		// Generate and inserts possible states
		g.NextPossibleStates(state)
	}
	// Lowest score was not found as all running states ended
	g.LowestScore = -1
}

func GetLowestScore(states []*GameState) *GameState {
	return states[0]
}

func (g *Game) NextPossibleStates(s *GameState) {
	// Filter the list of pods down to ones that are not home
	movablePods := slices.Filter(g.Pods, func(p Pod) bool { return !g.PodHome(s, p) })

	// All pods are home, so return the same state but marked as finished
	if len(movablePods) == 0 {
		s.Finished = true
		g.InsertGameState(s)
		return
	}
	for _, pod := range movablePods {
		moves := g.PossibleMoves(s, pod)
		for _, move := range moves {
			state := g.CreateNextState(s, pod, move)
			playedScore, ok := g.Played[state.Hash(g)]
			if ok && playedScore <= state.Score {
				continue
			}
			g.Played[state.Hash(g)] = state.Score
			g.InsertGameState(state)
		}
	}
}

func (g *Game) PodHome(s *GameState, p Pod) bool {
	podTile := s.PodPositions[p]
	podType := p.PodType()

	indexHome := slices.IndexOf(g.HomeTiles[podType], podTile)
	// Any pod that is in the hallway or is not in one of its home squares, its not home
	if indexHome == -1 {
		return false
	}

	// If it is home, then all previous home tiles have to be the same
	for i := indexHome + 1; i < len(g.HomeTiles[podType]); i++ {
		tile := g.HomeTiles[podType][i]
		if s.TilePositions[tile] != "" && s.TilePositions[tile].PodType() != podType {
			return false
		}
	}
	return true
}

func (g *Game) PossibleMoves(s *GameState, p Pod) []*Move {
	possibleMoves := g.Moves[s.PodPositions[p]] // Get the base list of possible moves

	// Filter out moves where the path is blocked by other pods
	inUseTiles := maps.Keys(s.TilePositions)
	possibleMoves = slices.Filter(possibleMoves, func(m *Move) bool {
		if slices.Contains(inUseTiles, m.EndTile) {
			return false
		}
		return !slices.ContainsAny(m.BlockTiles, inUseTiles)
	})

	// Filter out moves from the hallway to a room that doesn't match the pod
	possibleMoves = slices.Filter(possibleMoves, g.ValidMoveFunc(p, s))
	return possibleMoves
}

func (g *Game) ValidMoveFunc(p Pod, s *GameState) func(m *Move) bool {
	return func(m *Move) bool {
		targetType, targetPodType := m.EndTile.TileType()
		if targetType == "h" {
			return true
		}
		if targetPodType != p.PodType() {
			return false
		}

		homeTiles := g.HomeTiles[p.PodType()]
		targetIndex := slices.IndexOf(homeTiles, m.EndTile)

		// Check that all home squares before are free
		for i := 0; i < targetIndex; i++ {
			if s.TilePositions[homeTiles[i]] != "" {
				return false
			}
		}

		// Check if all home square after has a same type pod
		for i := targetIndex + 1; i < len(homeTiles); i++ {
			if s.TilePositions[homeTiles[i]] == "" || s.TilePositions[homeTiles[i]].PodType() != p.PodType() {
				return false
			}
		}
		return true

	}
}

func (g *Game) CreateNextState(prevState *GameState, p Pod, m *Move) *GameState {
	// First clone the previous state for the new state
	newState := prevState.Clone()

	// Remove the pods current locations
	delete(newState.TilePositions, newState.PodPositions[p])
	delete(newState.PodPositions, p)

	// Add the new locations
	newState.TilePositions[m.EndTile] = p
	newState.PodPositions[p] = m.EndTile

	// Increase the score
	newState.Score += (p.ScoreMultiplier() * m.Steps)
	newState.PrevMoves = append(newState.PrevMoves, fmt.Sprintf("%s:%s->%s", p, prevState.PodPositions[p], m.EndTile))
	return newState
}

type Path struct {
	StartTile, EndTile Tile
}

func (g *Game) InitMoves() map[Tile][]*Move {
	// Generate a lookup list of moves between any two tiles (no pod based restrictions)
	hallways := []Tile{"h0", "h1", "h3", "h5", "h7", "h9", "h10"}
	rooms := []Tile{}
	for _, v := range g.HomeTiles {
		rooms = append(rooms, v...)
	}

	paths := []Path{}
	for _, r := range rooms {
		for _, h := range hallways {
			paths = append(paths, Path{r, h})
		}
	}

	stepsMap := InitSteps()
	blockersMap := InitBlockTiles()
	results := make(map[Tile][]*Move)
	for _, path := range paths {
		stepOffset := 0
		addBlockers := []Tile{}
		startType, startPodType := path.StartTile.TileType()
		if startType == "b" {
			stepOffset = slices.IndexOf(g.HomeTiles[startPodType], path.StartTile)
			addBlockers = append(addBlockers, g.HomeTiles[startPodType][0:stepOffset]...)
		}
		mapPath := Path{g.HomeTiles[startPodType][0], path.EndTile}

		steps := stepsMap[mapPath] + stepOffset
		blockers := blockersMap[mapPath]
		blockers = append(blockers, addBlockers...)

		moveForward := &Move{
			Steps:      steps,
			EndTile:    path.EndTile,
			BlockTiles: blockers,
		}
		if results[path.StartTile] == nil {
			results[path.StartTile] = []*Move{}
		}
		results[path.StartTile] = append(results[path.StartTile], moveForward)

		moveBack := &Move{
			Steps:      steps,
			EndTile:    path.StartTile,
			BlockTiles: blockers,
		}
		if results[path.EndTile] == nil {
			results[path.EndTile] = []*Move{}
		}
		results[path.EndTile] = append(results[path.EndTile], moveBack)
	}
	return results
}

func InitSteps() map[Path]int {
	// Only record from the top of rooms to hallway rooms
	return map[Path]int{
		{"bA0", "h0"}:  3,
		{"bA0", "h1"}:  2,
		{"bA0", "h3"}:  2,
		{"bA0", "h5"}:  4,
		{"bA0", "h7"}:  6,
		{"bA0", "h9"}:  8,
		{"bA0", "h10"}: 9,

		{"bB0", "h0"}:  5,
		{"bB0", "h1"}:  4,
		{"bB0", "h3"}:  2,
		{"bB0", "h5"}:  2,
		{"bB0", "h7"}:  4,
		{"bB0", "h9"}:  6,
		{"bB0", "h10"}: 7,

		{"bC0", "h0"}:  7,
		{"bC0", "h1"}:  6,
		{"bC0", "h3"}:  4,
		{"bC0", "h5"}:  2,
		{"bC0", "h7"}:  2,
		{"bC0", "h9"}:  4,
		{"bC0", "h10"}: 5,

		{"bD0", "h0"}:  9,
		{"bD0", "h1"}:  8,
		{"bD0", "h3"}:  6,
		{"bD0", "h5"}:  4,
		{"bD0", "h7"}:  2,
		{"bD0", "h9"}:  2,
		{"bD0", "h10"}: 3,
	}
}

func InitBlockTiles() map[Path][]Tile {
	// Only record paths from the top of rooms to hallway tiles
	return map[Path][]Tile{
		{"bA0", "h0"}:  {"h1"},
		{"bA0", "h1"}:  nil,
		{"bA0", "h3"}:  nil,
		{"bA0", "h5"}:  {"h3"},
		{"bA0", "h7"}:  {"h3", "h5"},
		{"bA0", "h9"}:  {"h3", "h5", "h7"},
		{"bA0", "h10"}: {"h3", "h5", "h7", "h9"},

		{"bB0", "h0"}:  {"h3", "h1"},
		{"bB0", "h1"}:  {"h3"},
		{"bB0", "h3"}:  nil,
		{"bB0", "h5"}:  nil,
		{"bB0", "h7"}:  {"h5"},
		{"bB0", "h9"}:  {"h5", "h7"},
		{"bB0", "h10"}: {"h5", "h7", "h9"},

		{"bC0", "h0"}:  {"h5", "h3", "h1"},
		{"bC0", "h1"}:  {"h5", "h3"},
		{"bC0", "h3"}:  {"h5"},
		{"bC0", "h5"}:  nil,
		{"bC0", "h7"}:  nil,
		{"bC0", "h9"}:  {"h7"},
		{"bC0", "h10"}: {"h7", "h9"},

		{"bD0", "h0"}:  {"h7", "h5", "h3", "h1"},
		{"bD0", "h1"}:  {"h7", "h5", "h3"},
		{"bD0", "h3"}:  {"h7", "h5"},
		{"bD0", "h5"}:  {"h7"},
		{"bD0", "h7"}:  nil,
		{"bD0", "h9"}:  nil,
		{"bD0", "h10"}: {"h9"},
	}
}

func (m *Move) String() string {
	out := fmt.Sprintf("|%s:%d", m.EndTile, m.Steps)
	return out
}

func (s *GameState) String() string {
	out := s.hash
	out += fmt.Sprintf("SCORE:%d\n", s.Score)
	for _, v := range s.PrevMoves {
		out += fmt.Sprintf("%s\n", v)
	}
	return out
}

func (g *Game) InsertGameState(insert *GameState) {
	if g.RunningList == nil {
		g.RunningList = insert
		return
	}
	// Insert in order
	g.insertInOrder(insert)
}

func (g *Game) insertInOrder(insert *GameState) {
	added := false
	current, prev := g.RunningList, (*GameState)(nil)
	for current != nil {
		// If the currrent item is greater insert it before that score
		if insert.Score < current.Score {
			if prev != nil {
				prev.Next = insert
			} else {
				g.RunningList = insert
			}
			added = true
			insert.Next = current
			break
		}
		// Move to the next item on the list
		prev = current
		current = current.Next
	}
	if !added {
		prev.Next = insert
	}
}

func (g *Game) LengthRunning() int {
	count := 0
	ptr := g.RunningList
	for ptr != nil {
		count++
		ptr = ptr.Next
	}
	return count
}
