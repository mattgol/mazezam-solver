
Solver for Mazezam Levels
=========================

Mazezam is a block moving puzzle where all movable blocks in a row move together. The player has to find a combination of moves which allow the player to reach the exit. The [website of the inventor](https://sites.google.com/site/malcolmsprojects/mazezam-home-page) contains more background and a long list of implementations.

This repository contains a solver which performs an brute-force exhaustive search through all reachable positions to find an optimal solution.

And yes: I'm aware that there's already a [solver](https://github.com/Malcohol/mzm-designer)... I've built my own solver to solve a level which I couldn't find a solution myself and typed in the levels with an editor. Only when searching for the license of the levels I came across everything else: level designer, solver, source code to the levels. And I almost used the same format for the levels except two characters and the comments which made it easy to make my solver compatible.

Building instructions
=====================

The solver is written in go and only interacts using text files. It has no dependencies for building other than the standard library except for an assertion library used only for testing.

Building:

```
go build
```

Usage
=====

The solver uses text files with the level setup as input - the levels can be edited using Malcom Tyrell's [level editor](https://github.com/Malcohol/mzm-designer).

Each of these characters represents an entity in the level:
- `#`: Wall
- `+`: Start position of the player (only one allowed)
- `$`: Movable object
- `*`: Exit

Unrecognized characters will be interpreted as walls. Empy lines and lines starting with an `;` will be ignored.

The level will only be recognized as valid when it is surrounded by a Wall (with out without an exit) and only contains a single player position).

If the starting position is inside the left wall (as in the original levels) the player will be immediately moved right of the wall as this is the actual starting position in the game.

To start the solver start the compiled binary with the filename of the level file:

```
./mazezamsolver LEVELFILE
```

It prints the level as interpreted by the solver and one of the shortest solutions along with some statistical information.

Examples
========

The directory `examples` contains the original levels published with the [1k version for the ZX](https://github.com/Malcohol/1kMazezaM/blob/master/levels.mzm).



