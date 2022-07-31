package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyPrintSolution(t *testing.T) {
	s := PrettyPrintSolution("0123456789012345678901234567890")
	r := "0123 4567 8901 2345 \n6789 0123 4567 890 \n"
	assert.Equal(t, r, s)
}

func TestSolve(t *testing.T) {
	g := NewGame()

	// simple solution
	solution, err := g.Solve("####\n#+ E\n####")
	assert.Nil(t, err)
	fmt.Println("Solution", solution)
	assert.Equal(t, "rre", solution)

	// unsolvable
	_, err = g.Solve("####\n#+ #\n####")
	assert.NotNil(t, err)
	assert.Equal(t, "no solution found", err.Error())

	// more complex solution
	solution, err = g.Solve(`
	## o o E
	#+ o   #
	`)
	assert.Nil(t, err)
	assert.Equal(t, "rrruldrrurre", solution)
}
