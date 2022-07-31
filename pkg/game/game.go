package game

import (
	"fmt"
	"mazezamsolver/pkg/playfield"
	"strings"
)

const (
	Up = iota
	Left
	Down
	Right
	Exit
)

//
// Game recursively tries moves to find a solution

type Game struct {
	// map of playfield.Playfield
	field map[string]string
	start *playfield.Playfield
}

func NewGame() *Game {
	return &Game{
		field: make(map[string]string),
		start: nil,
	}
}

func (g *Game) Clear() {
	g.field = make(map[string]string)
	g.start = nil
}

var moveMarker map[int]string = map[int]string{
	Up:    "u",
	Down:  "d",
	Left:  "l",
	Right: "r",
	Exit:  "e",
}

func (g *Game) Stats() string {
	stats := ""
	stats += fmt.Sprintf("Positions reached: %d\n", len(g.field))
	longest := 0
	for _, s := range g.field {
		if len(s) > longest {
			longest = len(s)
		}
	}
	// stats += fmt.Sprintf("Longest path without repeat: %d\n", longest)
	solution, count := g.GetProgressSolution()
	stats += fmt.Sprintf("Number of solutions found: %d\n", count)
	stats += fmt.Sprintf("Shortest solution: \n%s\n", PrettyPrintSolution(solution))
	return stats
}

func PrettyPrintSolution(s string) string {
	linecount := 0
	out := ""
	for x := 0; x < len(s); x += 4 {
		start := x
		end := x + 4
		if end > len(s) {
			end = len(s)
		}
		out += s[start:end] + " "
		linecount++
		if linecount >= 4 {
			linecount = 0
			out += "\n"
		}
	}
	return out
}

// Print all positions reached by player
func (g *Game) AllPlayerPositions() {
	pp := g.start.Copy()
	for s := range g.field {
		pf := playfield.FromString(s)
		x, y := pf.GetPlayer()
		pp.Set(x, y, playfield.Player)
	}
	fmt.Println(pp.Print())
}

// iterate over field entries and expand list
func (g *Game) ProgressList(thisPositions []string) []string {
	nextPositions := []string{}
	for _, s := range thisPositions {
		pf := playfield.FromString(s)
		next := g.ProgressPosition(pf, g.field[s])
		nextPositions = append(nextPositions, next...)
	}
	return nextPositions
}

// expand on a single position
func (g *Game) ProgressPosition(pf *playfield.Playfield, prev string) []string {
	next := pf.PossibleMoves()
	nextPosition := []string{}
	for _, n := range next {
		moves := prev + moveMarker[n]
		nextPf := pf.Copy()
		exited := nextPf.MovePlayer(n)
		s := nextPf.Stringify()
		_, beenHere := g.field[s]
		if !beenHere {
			if exited {
				moves += moveMarker[Exit]
			} else {
				nextPosition = append(nextPosition, s)
			}
			g.field[s] = moves
		}
	}
	return nextPosition
}

func (g *Game) GetProgressSolution() (string, int) {
	count := 0
	solution := ""
	for _, p := range g.field {
		if strings.HasSuffix(p, moveMarker[Exit]) {
			count++
			if count == 1 || len(p) < len(solution) {
				solution = p
			}
		}
	}
	return solution, count
}

// Solve Flat
func (g *Game) Solve(start string) (string, error) {
	// clear field
	g.Clear()
	// read start field
	pf := playfield.ReadPlayField(start)
	fmt.Println(pf.Print())
	if !pf.IsValid() {
		return "", fmt.Errorf("not a valid start playfield.Playfield")
	}
	// iterate
	g.start = pf
	nextPos := []string{pf.Stringify()}
	g.field[nextPos[0]] = ""
	for {
		nextPos = g.ProgressList(nextPos)
		if len(nextPos) == 0 {
			break
		}
	}
	solution, count := g.GetProgressSolution()
	if count > 0 {
		return solution, nil
	}
	return "", fmt.Errorf("no solution found")
}
