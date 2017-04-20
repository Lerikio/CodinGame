package main

import "fmt"

/* Helpers */
func abs(a int) int {
	if a < 0 {
		a = -a
	}
	return a
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

type Action int

var Actions = [5]string{
	"SLOWER",
	"FASTER",
	"PORT",
	"STARBOARD",
	"WAIT",
}

func coordToID(x int, y int) int {
	return x*100 + y
}

/*****************************************************************************/
/* Tile Class */
type Tile struct {
	x int
	y int
}

func distance(a Tile, b Tile) int {
	return (abs(a.x-b.x) + abs(a.x+a.y-b.x-b.y) + abs(a.y-b.y)) / 2
}

func (a *Tile) toID() int {
	return a.x*100 + a.y
}

func nextTile(x int, y int, d int) Tile {
	delta := Tile{0, 0}
	if d == 0 {
		delta = Tile{1, 0}
	} else if d == 1 {
		if y%2 == 0 {
			delta = Tile{0, -1}
		} else {
			delta = Tile{1, -1}
		}
	} else if d == 2 {
		if y%2 == 0 {
			delta = Tile{-1, -1}
		} else {
			delta = Tile{0, -1}
		}
	} else if d == 3 {
		delta = Tile{-1, 0}
	} else if d == 4 {
		if y%2 == 0 {
			delta = Tile{-1, 1}
		} else {
			delta = Tile{0, 1}
		}
	} else if d == 5 {
		if y%2 == 0 {
			delta = Tile{0, 1}
		} else {
			delta = Tile{1, 1}
		}
	}
	newX := x + delta.x
	newY := y + delta.y
	if newX < 23 && newX > 0 && newY < 21 && newY > 0 {
		return Tile{newX, newY}
	}
	return Tile{x, y}
}

/*****************************************************************************/
/* Ship Class */
type Ship struct {
	x    int  // X coordinate
	y    int  // Y coordinate
	d    int  // Direction
	s    int  // Speed
	r    int  // Rhum
	mine bool // True if owned by me
	id   int  // ID
}

func (ship *Ship) copy() Ship {
	return Ship{ship.x, ship.y, ship.d, ship.s, ship.r, ship.mine, ship.id}
}

func (ship *Ship) applyAction(act Action) []int {
	var touchedTiles []int
	switch act {
	case 0:
		if ship.s > 0 {
			ship.s--
		}
		arrival := ship.nextPosition()
		front := nextTile(arrival.x, arrival.y, ship.d)
		back := nextTile(arrival.x, arrival.y, (ship.d+3)%6)
		touchedTiles = []int{arrival.toID(), front.toID(), back.toID()}
		ship.x = arrival.x
		ship.y = arrival.y
	case 1:
		if ship.s < 2 {
			ship.s++
		}
		arrival := ship.nextPosition()
		front := nextTile(arrival.x, arrival.y, ship.d)
		back := nextTile(arrival.x, arrival.y, (ship.d+3)%6)
		touchedTiles = []int{arrival.toID(), front.toID(), back.toID()}
		ship.x = arrival.x
		ship.y = arrival.y
	case 2:
		arrival := ship.nextPosition()
		front := nextTile(arrival.x, arrival.y, ship.d)
		back := nextTile(arrival.x, arrival.y, (ship.d+3)%6)
		ship.x = arrival.x
		ship.y = arrival.y
		ship.d = (ship.d + 1) % 6
		pFront := nextTile(arrival.x, arrival.y, ship.d)
		pBack := nextTile(arrival.x, arrival.y, (ship.d+3)%6)
		touchedTiles = []int{front.toID(), arrival.toID(), back.toID(), pFront.toID(), pBack.toID()}
	case 3:
		arrival := ship.nextPosition()
		front := nextTile(arrival.x, arrival.y, ship.d)
		back := nextTile(arrival.x, arrival.y, (ship.d+3)%6)
		ship.x = arrival.x
		ship.y = arrival.y
		ship.d = (6 + ship.d - 1) % 6
		pFront := nextTile(arrival.x, arrival.y, ship.d)
		pBack := nextTile(arrival.x, arrival.y, (ship.d+3)%6)
		touchedTiles = []int{front.toID(), arrival.toID(), back.toID(), pFront.toID(), pBack.toID()}
	case 4:
		arrival := ship.nextPosition()
		front := nextTile(arrival.x, arrival.y, ship.d)
		back := nextTile(arrival.x, arrival.y, (ship.d+3)%6)
		touchedTiles = []int{arrival.toID(), front.toID(), back.toID()}
		ship.x = arrival.x
		ship.y = arrival.y
	}
	return touchedTiles
}

func (ship *Ship) nextPosition() Tile {
	result := Tile{ship.x, ship.y}
	for i := 0; i < ship.s; i++ {
		result = nextTile(result.x, result.y, ship.d)
	}
	return result
}

/*****************************************************************************/
/* Game class */
type Game struct {
	turn        int
	firstAction [3]Action
	myShips     []Ship
	enShips     []Ship
	barrels     map[int]int
	mines       map[int]bool
	balls       map[int]map[int]bool
}

func (g *Game) init() {
	g.barrels = make(map[int]int)
	g.mines = make(map[int]bool)
	g.balls = make(map[int]map[int]bool)
}

func (g *Game) initNext() Game {
	return Game{g.turn + 1, g.firstAction, g.myShips, g.enShips, g.barrels, g.mines, g.balls}
}

func (g *Game) myScore() int {
	score := 0
	for _, ship := range g.myShips {
		score += ship.r
	}
	return score
}

func (g *Game) apply(myAct []Action) {
	for i := range g.myShips {
		g.myShips[i].r--
		touchedTiles := g.myShips[i].applyAction(myAct[2-i])
		for _, tile := range touchedTiles {
			if g.mines[tile] {
				g.myShips[i].r = min(g.myShips[i].r-25, 0)
				delete(g.mines, tile)
			}
			if rum, ok := g.barrels[tile]; ok {
				g.myShips[i].r = max(g.myShips[i].r+rum, 100)
				delete(g.barrels, tile)
			}
			if g.balls[g.turn][tile] && i != 3 && i != 4 {
				if tile == coordToID(g.myShips[i].x, g.myShips[i].y) {
					g.myShips[i].r = min(g.myShips[i].r-50, 0)
				} else {
					g.myShips[i].r = min(g.myShips[i].r-25, 0)
				}
			}
		}
	}
	for i := range g.enShips {
		next := g.enShips[i].nextPosition()
		g.enShips[i].x = next.x
		g.enShips[i].y = next.y
	}
}

func (g *Game) simulate(numberOfTurns int, allActions [][3]Action) [3]Action {
	var sims [][]Game
	var scores []int
	sims = [][]Game{[]Game{*g}}

	for k := 1; k < numberOfTurns; k++ {
		sims = append(sims, make([]Game, 10))
		for _, prevG := range sims[k-1] {
			for index, act := range allActions {
				if len(prevG.myShips) == 1 && index > 4 {
					break
				} else if len(prevG.myShips) == 2 && index > 24 {
					break
				}
				nextGame := prevG.initNext()
				nextGame.apply(act[:])
				if k == 1 {
					nextGame.firstAction = act
				}
				sims[k] = append(sims[k], nextGame)
				if k == numberOfTurns-1 {
					scores = append(scores, nextGame.myScore())
				}
			}
		}
	}

	maximum := -1
	max_i := -1
	for i, score := range scores {
		if score > maximum {
			maximum = score
			max_i = i
		}
	}

	return sims[len(sims)-1][max_i].firstAction
}

/*************** Main Function *****************/
func main() {
	/* Persistent Values */
	// hasAttacked := make(map[int]bool)

	// Compute all possible actions
	var allActions [][3]Action
	for i := range Actions {
		for j := range Actions {
			for k := range Actions {
				allActions = append(allActions, [3]Action{Action(i), Action(j), Action(k)})
			}
		}
	}

	for {
		var cGame Game
		cGame.init()

		// myShipCount: the number of remaining ships
		var myShipCount int
		fmt.Scan(&myShipCount)

		// entityCount: the number of entities (e.g. ships, mines or cannonballs)
		var entityCount int
		fmt.Scan(&entityCount)

		for i := 0; i < entityCount; i++ {
			var entityID int
			var entityType string
			var x, y, arg1, arg2, arg3, arg4 int
			fmt.Scan(&entityID, &entityType, &x, &y, &arg1, &arg2, &arg3, &arg4)
			if entityType == "SHIP" {
				if arg4 == 0 {
					cGame.enShips = append(cGame.enShips, Ship{x, y, arg1, arg2, arg3, false, entityID})
				} else {
					cGame.myShips = append(cGame.myShips, Ship{x, y, arg1, arg2, arg3, true, entityID})
				}
			} else if entityType == "BARREL" {
				cGame.barrels[coordToID(x, y)] = arg1
			} else if entityType == "CANNONBALL" {
				if _, ok := cGame.balls[arg2]; !ok {
					cGame.balls[arg2] = make(map[int]bool)
				}
				cGame.balls[arg2][coordToID(x, y)] = true
			} else if entityType == "MINE" {
				cGame.mines[coordToID(x, y)] = true
			}
		}

		actions := cGame.simulate(15, allActions)
		for index := range cGame.myShips {
			fmt.Println(Actions[int(actions[index])])
		}
	}
}
