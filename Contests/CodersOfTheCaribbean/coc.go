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
	delta := Point{0, 0}
	if ship.d == 0 {
		delta = Point{1, 0}
	} else if ship.d == 1 {
		if ship.y%2 == 0 {
			delta = Point{0, -1}
		} else {
			delta = Point{1, -1}
		}
	} else if ship.d == 2 {
		if ship.y%2 == 0 {
			delta = Point{-1, -1}
		} else {
			delta = Point{0, -1}
		}
	} else if ship.d == 3 {
		delta = Point{-1, 0}
	} else if ship.d == 4 {
		if ship.y%2 == 0 {
			delta = Point{-1, 1}
		} else {
			delta = Point{0, 1}
		}
	} else if ship.d == 5 {
		if ship.y%2 == 0 {
			delta = Point{0, 1}
		} else {
			delta = Point{1, 1}
		}
	}

	for i := 0; i < ship.s; i++ {
		x := result.x + delta.x
		y := result.y + delta.y
		if x < 23 && x > 0 && y < 21 && y > 0 {
			result.x = x
			result.y = y
		}
	}

	return result
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
					for _, ennemy := range enShips {
						target := ennemy.nextPosition()
						distance := distance(Point{this.x, this.y}, target)
						if distance < 10 && distance < smallest {
							smallest = distance
							closest = target
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
				closest := Point{this.x, this.y}
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
				targetedPosition = append(targetedPosition, closest.toID())
				fmt.Println("MOVE", closest.x, closest.y)
			}

			// fmt.Fprintln(os.Stderr, "Debug messages...")
			//fmt.Printf("MOVE 11 10\n") // Any valid action, such as "WAIT" or "MOVE x y"
		}
	}
}
