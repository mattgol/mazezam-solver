package playfield

import (
	"fmt"
	"strings"
)

const (
	Space         = ' '
	Wall          = '#'
	MovableObject = '$'
	Player        = '+'
	Exit          = '*'
)

const (
	Up = iota
	Left
	Down
	Right
)

//
// PlayField is a byte array with these
// entries: [0] = width, [1] = height, [x+y*width+2] = contents
type Playfield struct {
	squares []byte
}

func NewPlayfield(width, height int) *Playfield {
	if width > 255 || height > 255 {
		fmt.Println("Playfield maximum size of 255x255 exceeded.")
		panic("playfield size exceeded")
	}
	f := make([]byte, width*height+2)
	f[0] = byte(width)
	f[1] = byte(height)
	for p := 0; p < width*height; p++ {
		f[p+2] = Wall
	}
	return &Playfield{
		squares: f,
	}
}

func (pf *Playfield) Copy() *Playfield {
	sq := make([]byte, pf.Width()*pf.Height()+2)
	for k, v := range pf.squares {
		sq[k] = v
	}
	return &Playfield{
		squares: sq,
	}
}

func (pf *Playfield) Stringify() string {
	// simply return playfield as string
	return string(pf.squares)
}

func FromString(input string) *Playfield {
	b := []byte(input)
	pf := Playfield{
		squares: b,
	}
	return &pf
}

func (pf *Playfield) Width() int {
	return int(pf.squares[0])
}

func (pf *Playfield) Height() int {
	return int(pf.squares[1])
}

func (pf *Playfield) squaresIndex(x, y int) int {
	pos := x + y*int(pf.squares[0]) + 2
	return pos
}

func (pf *Playfield) Get(x, y int) byte {
	if x >= pf.Width() {
		return Wall
	}
	if y >= pf.Height() {
		return Wall
	}
	return pf.squares[pf.squaresIndex(x, y)]
}

func (pf *Playfield) Set(x, y int, value byte) {
	pos := pf.squaresIndex(x, y)
	if pos < 2 || pos > len(pf.squares) {
		return // ignore outside values
	}
	pf.squares[pos] = value
}

func (pf *Playfield) getLine(y int) []byte {
	width := pf.Width()
	line := pf.squares[y*width+2 : (y+1)*width+2]
	return line
}

func (pf *Playfield) GetPlayer() (x, y int) {
	pos := strings.IndexRune(string(pf.squares[2:]), Player)
	if pos < 0 {
		return 0, 0
	}
	y = pos / pf.Width()
	x = pos % pf.Width()
	return
}

func (pf *Playfield) Print() string {
	out := ""
	height := pf.Height()
	for y := 0; y < height; y++ {
		line := pf.getLine(y)
		out += string(line) + "\n"
	}
	return out
}

func ReadPlayField(input string) *Playfield {
	lines := strings.Split(input, "\n")
	count := len(lines)
	var width int = 0
	var height int = 0
	for y := 0; y < count; y++ {
		line := strings.TrimSpace(lines[y])
		l := len(line)
		if l == 0 || strings.HasPrefix(line, ";") {
			continue
		}
		if l > width {
			width = l
		}
		height++
	}
	pf := NewPlayfield(width, height)
	y := 0
	for _, lineRaw := range lines {
		line := strings.TrimSpace(lineRaw)
		l := len(line)
		if l == 0 || strings.HasPrefix(line, ";") {
			continue
		}
		for x, s := range line {
			if s != Space && s != Wall && s != MovableObject && s != Player && s != Exit {
				s = Wall
			}
			pf.Set(x, y, byte(s))
		}
		y++
	}
	for y := 0; y < height; y++ {
		if pf.Get(0, y) == Player {
			pf.Set(0, y, Wall)
			pf.Set(1, y, Player)
		}
	}
	return pf
}

func (pf *Playfield) CanShiftLeft(y int) bool {
	if y >= pf.Height() {
		return false
	}
	line := pf.getLine(y)
	for x := pf.Width() - 1; x >= 1; x-- {
		atPos := rune(line[x])
		leftPos := rune(line[x-1])
		leftPosBlock := leftPos == Wall || leftPos == Exit
		atPosSomething := atPos == MovableObject || atPos == Player
		if leftPosBlock && atPosSomething {
			return false
		}
	}
	return true
}

func (pf *Playfield) CanShiftRight(y int) bool {
	if y >= pf.Height() {
		return false
	}
	line := pf.getLine(y)
	for x := 0; x < int(pf.Width())-1; x++ {
		atPos := rune(line[x])
		rightPos := rune(line[x+1])
		rightPosBlock := rightPos == Wall || rightPos == Exit
		atPosSomething := atPos == MovableObject || atPos == Player
		if rightPosBlock && atPosSomething {
			return false
		}
	}
	return true
}

func (pf *Playfield) ShiftLeft(y int) {
	if y >= pf.Height() {
		return
	}
	line := pf.getLine(y)
	for x := 0; x < int(pf.Width())-1; x++ {
		atPos := rune(line[x])
		atPosNoBlock := atPos != Wall && atPos != Exit
		rightPos := rune(line[x+1])
		rightPosNoBlock := rightPos != Wall && rightPos != Exit
		// if atPos != Wall && atPos != Exit && rightPos != Wall && rightPos != Exit {
		// 	line[x] = line[x+1]
		// }
		if atPosNoBlock {
			if rightPosNoBlock {
				line[x] = line[x+1] // moves
			} else {
				line[x] = Space // moves away
			}
		}

	}
}

func (pf *Playfield) ShiftRight(y int) {
	if y >= pf.Height() {
		return
	}
	line := pf.getLine(y)
	for x := int(pf.Width() - 1); x > 0; x-- {
		atPos := rune(line[x])
		atPosNoBlock := atPos != Wall && atPos != Exit
		leftPos := rune(line[x-1])
		leftPosNoBlock := leftPos != Wall && leftPos != Exit
		if atPosNoBlock {
			if leftPosNoBlock {
				line[x] = line[x-1] // moves
			} else {
				line[x] = Space // moves away
			}
		}
	}
}

func (pf *Playfield) PossibleMoves() []int {
	x, y := pf.GetPlayer()

	//left
	target := pf.Get(x-1, y)
	canMoveLeft := target == Space || target == Exit ||
		(target == MovableObject && pf.CanShiftLeft(y))

	target = pf.Get(x+1, y)
	canMoveRight := target == Space || target == Exit ||
		(target == MovableObject && pf.CanShiftRight(y))

	target = pf.Get(x, y-1)
	canMoveUp := target == Space || target == Exit

	target = pf.Get(x, y+1)
	canMoveDown := target == Space || target == Exit

	moves := []int{}
	if canMoveLeft {
		moves = append(moves, Left)
	}
	if canMoveRight {
		moves = append(moves, Right)
	}
	if canMoveUp {
		moves = append(moves, Up)
	}
	if canMoveDown {
		moves = append(moves, Down)
	}
	return moves
}

// this does not check again if move is possible
func (pf *Playfield) MovePlayer(direction int) (exited bool) {
	x, y := pf.GetPlayer()
	pf.Set(x, y, Space)
	switch direction {
	case Up:
		exited = pf.Get(x, y-1) == Exit
		pf.Set(x, y-1, Player)
	case Down:
		exited = pf.Get(x, y+1) == Exit
		pf.Set(x, y+1, Player)
	case Left:
		leftPos := pf.Get(x-1, y)
		exited = leftPos == Exit
		if leftPos == MovableObject {
			pf.ShiftLeft(y)
		}
		pf.Set(x-1, y, Player)
	case Right:
		rightPos := pf.Get(x+1, y)
		exited = rightPos == Exit
		if rightPos == MovableObject {
			pf.ShiftRight(y)
		}
		pf.Set(x+1, y, Player)
	}
	return exited
}

func (pf *Playfield) IsValid() bool {
	// has one player
	player := 0
	for _, c := range pf.squares {
		if c == Player {
			player++
		}
	}
	if player != 1 {
		return false
	}
	// surrounded by Wall or Exit
	for _, c := range pf.getLine(0) {
		if c != Wall && c != Exit {
			return false
		}
	}
	for _, c := range pf.getLine(pf.Height() - 1) {
		if c != Wall && c != Exit {
			return false
		}
	}
	for k := 0; k < int(pf.Height()); k++ {
		c := pf.Get(0, k)
		if c != Wall && c != Exit {
			return false
		}
		c = pf.Get(pf.Width()-1, k)
		if c != Wall && c != Exit {
			return false
		}
	}
	return true
}
