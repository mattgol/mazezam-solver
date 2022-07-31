package main

import (
	"fmt"
	"mazezamsolver/pkg/game"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mazezam-solver file")
		os.Exit(1)
	}
	g := game.NewGame()
	f, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Could not read file", os.Args[1])
		os.Exit(1)
	}
	_, err = g.Solve(string(f))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Print(g.Stats())
}
