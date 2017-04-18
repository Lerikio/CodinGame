package main

import "fmt"

/* Helpers **************************/
func abs(a int) int {
	if a < 0 {
		a = -a
	}
	return a
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
	mine bool // True if owned by me
	id   int  // ID
}

func (ship *Ship) closestBarrel(barrels *[]Point) Point {
	result := Point{0, 0}
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
	ghost := Ship{ship.x, ship.y, ship.d, ship.s, ship.mine, ship.id}
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
	hasAttacked := make(map[int]bool)
	for {
		var myShips []Ship
		var enShips []Ship
		var barrels []int
		var mines []int
		var targetedPosition []int
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
					enShips = append(enShips, Ship{x, y, arg1, arg2, false, entityID})
				} else {
					myShips = append(myShips, Ship{x, y, arg1, arg2, true, entityID})
				}
			} else if entityType == "BARREL" {
				barrels = append(barrels, coordToID(x, y))
			} else if entityType == "CANNONBALL" {
				if _, ok := balls[arg2]; !ok {
					balls[arg2] = make(map[int]bool)
				}
				balls[arg2][coordToID(x, y)] = true
			} else if entityType == "MINE" {
				mines = append(mines, coordToID(x, y))
			}
		}
		for i := 0; i < myShipCount; i++ {
			this := myShips[i]
			firing := false
			if !hasAttacked[this.id] {
				for _, mine := range mines {
					distance := distance(Point{this.x, this.y}, IDToCoord(mine))
					if distance > 2 && distance < 10 {
						alreadyDestroyed := false
						for _, array := range balls {
							if array[mine] {
								alreadyDestroyed = true
								break
							}
						}
						if !alreadyDestroyed {
							firing = true
							minePoint := IDToCoord(mine)
							fmt.Println("FIRE", minePoint.x, minePoint.y)
							break
						}
					}
				}
				if !firing {
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
						firing = true
						fmt.Println("FIRE", closest.x, closest.y)
					}
				}
			}
			if firing {
				hasAttacked[this.id] = true
			} else {
				ghost := Ship(this)
				nextCenter := ghost.nextPosition()
				ghost.x = nextCenter.x
				ghost.y = nextCenter.y
				nextFront := ghost.nextPosition()
				nextBack := Point{this.x, this.y}
				if balls[2][nextCenter.toID()] {
					if this.s == 2 {
						fmt.Println("SLOWER")
					} else {
						fmt.Println("FASTER")
					}
				} else if balls[2][nextFront.toID()] {
					if entityCount%2 == 0 {
						fmt.Println("PORT")
					} else {
						fmt.Println("STARBOARD")
					}
				} else if balls[2][nextBack.toID()] && this.s == 1 {
					fmt.Println("FASTER")
				} else {
					closest := Point{-1, -1}
					smallest := 10000
					for _, barrel := range barrels {
						distance := distance(Point{this.x, this.y}, IDToCoord(barrel))
						if distance < smallest {
							alreadyTargeted := false
							for _, point := range targetedPosition {
								if barrel == point {
									alreadyTargeted = true
									break
								}
							}
							if !alreadyTargeted {
								closest = IDToCoord(barrel)
								smallest = distance
							}
						}
					}
					hasAttacked[this.id] = false
					if closest.x == -1 {
						closest.x = enShips[0].x
						closest.y = enShips[0].y
					} else {
						targetedPosition = append(targetedPosition, closest.toID())
					}
					fmt.Println("MOVE", closest.x, closest.y)
				}
			}

			// fmt.Fprintln(os.Stderr, "Debug messages...")
			//fmt.Printf("MOVE 11 10\n") // Any valid action, such as "WAIT" or "MOVE x y"
		}
	}
}
