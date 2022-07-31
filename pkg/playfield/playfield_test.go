package playfield

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//
// PlayField

const testplayfield = `##########
#        #
# $ $$$  #
# $$ $$$ #
# +      *
##########`

func TestNewPlayfield(t *testing.T) {
	pf := NewPlayfield(4, 3)
	assert.EqualValues(t, 4, pf.squares[0], "should put correct width in array")
	assert.EqualValues(t, 3, pf.squares[1], "should put correct height in array")
	assert.EqualValues(t, "############", pf.squares[2:], "should prefill as Wall")
	// test panic
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		assert.Equal(t, "playfield size exceeded", r)
	}()
	NewPlayfield(300, 0)
}

func TestCopy(t *testing.T) {
	pf := ReadPlayField("###\n# *\n###")
	cp := pf.Copy()
	cp.Set(1, 1, Player)
	old := pf.Get(1, 1)
	assert.Equal(t, Space, rune(old))
}

func TestStringification(t *testing.T) {
	testString := string([]byte{4, 3}) + "#####+$##*##"
	pf := FromString(testString)
	assert.Equal(t, 4, pf.Width())
	assert.Equal(t, 3, pf.Height())
	assert.Equal(t, "####\n#+$#\n#*##\n", pf.Print())

	s := pf.Stringify()
	assert.Equal(t, testString, s)
}

func TestReadPlayField(t *testing.T) {
	pf := ReadPlayField(testplayfield)
	assert.Equal(t, 10, pf.Width())
	assert.Equal(t, testplayfield[:10], string(pf.squares[2:12]))
	assert.Equal(t, testplayfield[12:21], string(pf.squares[13:22]))

	pf = ReadPlayField("§§\n  ")
	assert.Equal(t, Wall, rune(pf.Get(0, 0)))
}

func TestPrint(t *testing.T) {
	pf := ReadPlayField(testplayfield)
	out := pf.Print()
	assert.Equal(t, testplayfield+"\n", out)
}

func TestGetPlayerPosition(t *testing.T) {
	pf := ReadPlayField(testplayfield)
	x, y := pf.GetPlayer()
	assert.Equal(t, 2, x)
	assert.Equal(t, 4, y)

	pf = ReadPlayField("  \n  ")
	x, y = pf.GetPlayer()
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)
}

func TestGet(t *testing.T) {
	pf := ReadPlayField(testplayfield)
	assert.Equal(t, Wall, rune(pf.Get(20, 0)))
	assert.Equal(t, Wall, rune(pf.Get(0, 20)))
}

func TestSet(t *testing.T) {
	pf := NewPlayfield(2, 2)
	pf.Set(0, 0, Space)
	pf.Set(1, 1, MovableObject)
	pf.Set(3, 3, Space)
	assert.Equal(t, Space, rune(pf.squares[2]))
	assert.Equal(t, Wall, rune(pf.squares[3]))
	assert.Equal(t, Wall, rune(pf.squares[4]))
	assert.Equal(t, MovableObject, rune(pf.squares[5]))
}

func TestCanShift(t *testing.T) {
	testlines := "#  $#\n# $ #\n#$  #\n#+$ #\n# $+#"
	testLeft := []bool{true, true, false, false, true}
	testRight := []bool{false, true, true, true, false}
	pf := ReadPlayField(testlines)
	for y := 0; y < pf.Height(); y++ {
		canShiftLeft := pf.CanShiftLeft(y)
		assert.Equal(t, testLeft[y], canShiftLeft, "should test left shift of "+string(pf.getLine(y)))
		canShiftRight := pf.CanShiftRight(y)
		assert.Equal(t, testRight[y], canShiftRight, "should test right shift of "+string(pf.getLine(y)))
	}
	// out of bounds
	cannot := pf.CanShiftLeft(pf.Height())
	assert.False(t, cannot)
	cannot = pf.CanShiftRight(pf.Height())
	assert.False(t, cannot)
}

func TestPossibleMoves(t *testing.T) {
	pf := ReadPlayField(testplayfield)
	x, y := pf.GetPlayer()

	moves := pf.PossibleMoves()
	assert.ElementsMatch(t, moves, []int{Left, Right})
	pf.Set(x, y, Space)

	pf.Set(3, 2, Player)
	pf.Print()
	moves = pf.PossibleMoves()
	assert.ElementsMatch(t, moves, []int{Left, Right, Up})
	pf.Set(3, 2, Space)

	pf.Set(3, 1, Player)
	pf.Print()
	moves = pf.PossibleMoves()
	assert.ElementsMatch(t, moves, []int{Left, Right, Down})
}

func TestShiftLeft(t *testing.T) {
	pf := ReadPlayField("# $$ $$ *")
	pf.ShiftLeft(0)
	assert.Equal(t, []byte("#$$ $$  *"), pf.getLine(0))
	pf.ShiftLeft(1)
	assert.Equal(t, []byte("#$$ $$  *"), pf.getLine(0))
	pf = ReadPlayField("# $ $$*")
	pf.ShiftLeft(0)
	assert.Equal(t, []byte("#$ $$ *"), pf.getLine(0))
}

func TestShiftRight(t *testing.T) {
	pf := ReadPlayField("# $$+$$ *")
	pf.ShiftRight(0)
	assert.Equal(t, []byte("#  $$+$$*"), pf.getLine(0))
	pf.ShiftRight(1)
	assert.Equal(t, []byte("#  $$+$$*"), pf.getLine(0))
	pf = ReadPlayField("#$ $$ *")
	pf.ShiftRight(0)
	assert.Equal(t, []byte("# $ $$*"), pf.getLine(0))
}

func TestMovePlayer(t *testing.T) {
	pf := ReadPlayField("###*####\n# $$ $ #\n# $ $  #\n#+     #\n########")
	var moves []int = []int{Right, Right, Up, Down, Up, Right, Up, Left}
	for _, m := range moves {
		exited := pf.MovePlayer(m)
		assert.False(t, exited)
	}
	ref := "###*####\n#$$+$  #\n#  $ $ #\n#      #\n########\n"
	out := pf.Print()
	fmt.Print(out)
	assert.Equal(t, ref, out)
	// move out
	exited := pf.MovePlayer(Up)
	assert.True(t, exited)
}

func TestIsValid(t *testing.T) {
	// valid
	pf := ReadPlayField("####\n#+ #\n####")
	assert.True(t, pf.IsValid())

	// player
	// no player
	pf.Set(1, 1, Space)
	assert.False(t, pf.IsValid())
	pf.Set(1, 1, Player)
	// two players
	pf.Set(2, 1, Player)
	assert.False(t, pf.IsValid())
	pf.Set(2, 1, Space)

	// hole in the walls
	// left
	pf.Set(0, 1, Space)
	assert.False(t, pf.IsValid())
	pf.Set(0, 1, Wall)
	// right
	pf.Set(3, 1, Space)
	assert.False(t, pf.IsValid())
	pf.Set(3, 1, Wall)
	// top
	pf.Set(1, 0, Space)
	assert.False(t, pf.IsValid())
	pf.Set(1, 0, Wall)
	// bottom
	pf.Set(1, 2, Space)
	assert.False(t, pf.IsValid())
	pf.Set(1, 2, Wall)
}
