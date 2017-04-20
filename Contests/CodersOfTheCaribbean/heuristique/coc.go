package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

/* Helpers **************************/
func abs(a int) int {
	if a < 0 {
		a = -a
	}
	return a
}

type Action int

var Actions = [5]string{
	"SLOWER",
	"FASTER",
	"PORT",
	"STARBOARD",
	"WAIT",
}

/* Point Class */
type Point struct {
	x int
	y int
}

func distance(a Point, b Point) int {
	return (abs(a.x-b.x) + abs(a.x+a.y-b.x-b.y) + abs(a.y-b.y)) / 2
}

func coordToID(x int, y int) int {
	return x*100 + y
}

func (a *Point) toID() int {
	return a.x*100 + a.y
}

func IDToCoord(i int) Point {
	x := i / 100
	y := i % 100
	return Point{x, y}
}

func nextCase(x int, y int, d int) Point {
	delta := Point{0, 0}
	if d == 0 {
		delta = Point{1, 0}
	} else if d == 1 {
		if y%2 == 0 {
			delta = Point{0, -1}
		} else {
			delta = Point{1, -1}
		}
	} else if d == 2 {
		if y%2 == 0 {
			delta = Point{-1, -1}
		} else {
			delta = Point{0, -1}
		}
	} else if d == 3 {
		delta = Point{-1, 0}
	} else if d == 4 {
		if y%2 == 0 {
			delta = Point{-1, 1}
		} else {
			delta = Point{0, 1}
		}
	} else if d == 5 {
		if y%2 == 0 {
			delta = Point{0, 1}
		} else {
			delta = Point{1, 1}
		}
	}
	newX := x + delta.x
	newY := y + delta.y
	if newX < 23 && newX > 0 && newY < 21 && newY > 0 {
		return Point{newX, newY}
	}
	return Point{x, y}
}

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

func (ship *Ship) futureShip(startID int, dir int, speed int, act Action) Ship {
	result := Ship{ship.x, ship.y, dir, speed, ship.r, ship.mine, ship.id}
	switch act {
	case 0:
		if result.s > 0 {
			result.s--
		}
		arrival := result.nextPosition()
		result.x = arrival.x
		result.y = arrival.y
	case 1:
		if result.s < 2 {
			result.s++
		}
		arrival := result.nextPosition()
		result.x = arrival.x
		result.y = arrival.y
	case 2:
		arrival := result.nextPosition()
		result.x = arrival.x
		result.y = arrival.y
		result.d = (result.d + 1) % 6
	case 3:
		arrival := result.nextPosition()
		result.x = arrival.x
		result.y = arrival.y
		result.d = (6 + result.d - 1) % 6
	case 4:
		arrival := result.nextPosition()
		result.x = arrival.x
		result.y = arrival.y
	}
	return result
}

func (ship *Ship) nextPosition() Point {
	result := Point{ship.x, ship.y}
	for i := 0; i < ship.s; i++ {
		result = nextCase(result.x, result.y, ship.d)
	}
	return result
}

func (ship *Ship) possibleCollision(firingRange map[int]map[int]bool) []Point {
	var result []Point
	ghost := Ship{ship.x, ship.y, ship.d, ship.s, ship.r, ship.mine, ship.id}
	for k := 1; k <= 4; k++ {
		nextPosition := ghost.nextPosition()
		ghost.x = nextPosition.x
		ghost.y = nextPosition.y
		next := nextPosition.toID()
		if firingRange[k][next] {
			result = append(result, nextPosition)
		}
	}
	return result
}

func (ship *Ship) badPosition(turn int, mines map[int]bool, balls map[int]map[int]bool, act Action) bool {
	result := false
	front := nextCase(ship.x, ship.y, ship.d)
	back := nextCase(ship.x, ship.y, (ship.d+3)%6)
	previousFront := Point{front.x, front.y}
	switch act {
	case 2:
		previousFront = nextCase(ship.x, ship.y, (ship.d+1)%6)
	case 3:
		previousFront = nextCase(ship.x, ship.y, (6+ship.d-1)%6)
	}
	positions := [4]int{back.toID(), coordToID(ship.x, ship.y), front.toID(), previousFront.toID()}

	for _, position := range positions {
		if mines[position] || balls[turn][position] {
			result = true
			break
		}
	}
	return result
}

func (ship *Ship) pathfinder(depth int, mines map[int]bool, balls map[int]map[int]bool) []map[int]int {
	paths := make([]map[int][2]int, depth)
	visited := make(map[int][2]int)
	seen := make([]map[int]int, depth)
	seen[0] = make(map[int]int)
	seen[0][coordToID(ship.x, ship.y)] = -1
	visited[coordToID(ship.x, ship.y)] = [2]int{ship.d, ship.s}
	paths[0] = make(map[int][2]int)
	paths[0][coordToID(ship.x, ship.y)] = [2]int{ship.d, ship.s}

	for k := 1; k < depth; k++ {
		paths[k] = make(map[int][2]int)
		seen[k] = make(map[int]int)
		atLeastOneSolution := false
		for startID, val := range paths[k-1] {
			for act := range Actions {
				ghost := ship.futureShip(startID, val[0], val[1], Action(act))
				ghostID := coordToID(ghost.x, ghost.y)
				sameAsBefore := false
				if _, ok := visited[ghostID]; ok &&
					visited[ghostID][0] == ghost.d &&
					visited[ghostID][1] == ghost.s {
					sameAsBefore = true
				}
				if !sameAsBefore && !ghost.badPosition(k, mines, balls, Action(act)) {
					paths[k][ghostID] = [2]int{ghost.d, ghost.s}
					front := nextCase(ghost.x, ghost.y, ghost.d)
					back := nextCase(ghost.x, ghost.y, (ghost.d+3)%6)
					if k == 1 {
						seen[k][front.toID()] = act
						seen[k][ghostID] = act
						seen[k][back.toID()] = act
					} else {
						seen[k][front.toID()] = seen[k][startID]
						seen[k][ghostID] = seen[k][startID]
						seen[k][back.toID()] = seen[k][startID]
					}
					atLeastOneSolution = true
					visited[ghostID] = [2]int{ghost.d, ghost.s}
				}
			}
		}
		if !atLeastOneSolution {
			break
		}
	}
	return seen
}

/*************** Breadth First Search *****************/
func firingRange(start Point, breadth int) map[int]map[int]bool {
	fringes := make(map[int]map[int]bool)
	visited := make(map[int]bool)
	visited[start.toID()] = true
	fringes[0] = make(map[int]bool)
	fringes[0][start.toID()] = true

	for k := 1; k <= breadth; k++ {
		if _, ok := fringes[1+k/3]; !ok {
			fringes[1+k/3] = make(map[int]bool)
		}
		for start := range fringes[k/3] {
			startP := IDToCoord(start)
			for dir := 0; dir < 6; dir++ {
				neighborP := nextCase(startP.x, startP.y, dir)
				neighbor := neighborP.toID()
				if !visited[neighbor] {
					visited[neighbor] = true
					fringes[1+k/3][neighbor] = true
				}
			}
		}
	}
	return fringes
}

/*************** Main Function *****************/
func main() {
	/* Persistent Values */
	// hasAttacked := make(map[int]bool)
	for {
		var myShips []Ship
		var enShips []Ship
		var barrels []int
		mines := make(map[int]bool)
		targetedPosition := make(map[int]bool)
		balls := make(map[int]map[int]bool)

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
					enShips = append(enShips, Ship{x, y, arg1, arg2, arg3, false, entityID})
				} else {
					myShips = append(myShips, Ship{x, y, arg1, arg2, arg3, true, entityID})
				}
			} else if entityType == "BARREL" {
				barrels = append(barrels, coordToID(x, y))
			} else if entityType == "CANNONBALL" {
				if _, ok := balls[arg2]; !ok {
					balls[arg2] = make(map[int]bool)
				}
				balls[arg2][coordToID(x, y)] = true
			} else if entityType == "MINE" {
				mines[coordToID(x, y)] = true
			}
		}

		for index, this := range myShips {
			fmt.Fprint(os.Stderr, "Ship n°", index, " : ")
			paths := this.pathfinder(40, mines, balls)

			closestBarrel := -1
			smallestTurn := 100
			for _, barrel := range barrels {
				for turn, achievable := range paths {
					if turn >= smallestTurn {
						break
					}

					if _, exists := achievable[barrel]; exists && !targetedPosition[barrel] {
						closestBarrel = barrel
						smallestTurn = turn
					}
				}
			}
			fmt.Fprint(os.Stderr, "Coordinates of closest barrel are: ", closestBarrel, ", ")
			shouldAttack := false

			if closestBarrel != -1 {
				targetedPosition[closestBarrel] = true
				act := Actions[paths[smallestTurn][closestBarrel]]
				if act != "WAIT" {
					fmt.Fprint(os.Stderr, "Action is: ", act, "\n")
					fmt.Println(act)
				} else {
					shouldAttack = true
					fmt.Fprint(os.Stderr, "Should attack...\n")
				}
			} else {
				rand.Seed(time.Now().UTC().UnixNano())
				max := rand.Intn(len(paths[1]))
				var act string
				count := 0
				for _, val := range paths[1] {
					count++
					if count == max {
						act = Actions[val]
						break
					}
				}
				if act != "WAIT" {
					fmt.Fprint(os.Stderr, "Action is: ", act, "\n")
					fmt.Println(act)
				} else {
					fmt.Fprint(os.Stderr, "Should attack...\n")
					shouldAttack = true
				}
			}

			if shouldAttack {
				closest := Point{-1, -1}
				smallest := 100000
				explosions := firingRange(Point{this.x, this.y}, 10)
				for _, ennemy := range enShips {
					collisions := ennemy.possibleCollision(explosions)
					for _, collision := range collisions {
						if distance(collision, Point{this.x, this.y}) < smallest {
							closest = collision
						}
					}
				}
				if closest.x >= 0 {
					fmt.Fprint(os.Stderr, "ATTACK!")
					fmt.Println("FIRE", closest.x, closest.y)
				} else {
					fmt.Fprint(os.Stderr, "\nWAIT\n")
					fmt.Println("WAIT")
				}
			}
			// fmt.Fprintln(os.Stderr, "Debug messages...")
			//fmt.Printf("MOVE 11 10\n") // Any valid action, such as "WAIT" or "MOVE x y"
		}
	}
}
