package main

import (
	"fmt"
	"math/rand"
)

// Factory represents all the properties of a factory
type Factory struct {
	id               int
	owner            int
	pop              int
	prod             int
	turnsWithoutProd int
	baryDistance     int
}

// Troop represents all the properties of a single troop
type Troop struct {
	id    int
	owner int
	from  int
	to    int
	pop   int
	eta   int
}

// Bomb represents all the properties of a single bomb
type Bomb struct {
	id    int
	owner int
	from  int
	to    int
	eta   int
}

/****************************************************
                GAME object and methods
****************************************************/

// Game is a single object containing all the data representing the state of
// the game, allowing for easier access and the creation of methods.
type Game struct {
	playerID       int             // 1 for me, -1 for opponent
	bombMe         int             // My number of bombs
	bombOponnent   int             // The oponnent number of bombs
	currentTurn    int             // Turn of the simulation. 0 for actual turn
	score          int             // Score for current simulation. MyCyborgs - TheirCyborgs
	firstMove      [][4]int        // Move for turn 1 that allowed this configuration
	distances      [][]int         // Matrix of distances between factory i and factory j
	proximities    [][]int         // Matrix of the proximities of each factory. Line i is the ordered list of factories closest to factory i
	factories      []Factory       // All the factories in the game
	myFactories    []*Factory      // The Factories I own at this turn
	theirFactories []*Factory      // All other factories. No difference between ennemy and neutral
	troops         map[int][]Troop // All the troops in the game, referenced by their destination
	bombs          map[int][]Bomb  // All the bombs in the game, referenced by their destination
	seers          [][20]int       // Compute all the populations of the game for the 20 next turns
}

// Compute the proximity matrix
func (game *Game) initializeProximities() {
	game.proximities = make([][]int, len(game.distances[0]))
	for id := range game.distances[0] {
		closestToID := make([]int, len(game.distances[0])-1)
		for index := range closestToID {
			closestToID[index] = index
		}
		closestToID = game.quicksortFrom(id, closestToID)
		game.proximities[id] = closestToID[1:]
	}
}

func (game *Game) computeSeers() {
	game.seers = make([][20]int, len(game.factories))

	for i := range game.seers {
		game.seers[i] = game.computeSeer(&game.factories[i])
	}
}

func (game *Game) computeBaryDistances() {
	// fmt.Fprintln(os.Stderr, "BaryDistances:")
	for index := range game.factories {
		baryDistance := 0
		for _, myFactory := range game.myFactories {
			baryDistance += game.distances[index][myFactory.id]
		}
		baryDistance /= len(game.myFactories)
		game.factories[index].baryDistance = baryDistance
		// fmt.Fprintln(os.Stderr, index, "->", baryDistance)
	}
}

// Implementation of QuickSort in our specific case
func (game *Game) quicksortFrom(currentID int, closests []int) []int {
	if len(closests) < 2 {
		return closests
	}

	left, right := 0, len(closests)-1
	pivotIndex := rand.Int() % len(closests)

	closests[pivotIndex], closests[right] = closests[right], closests[pivotIndex]

	for i := range closests {
		if game.distances[currentID][closests[i]] <
			game.distances[currentID][closests[right]] {
			closests[i], closests[left] = closests[left], closests[i]
			left++
		}
	}

	closests[left], closests[right] = closests[right], closests[left]

	game.quicksortFrom(currentID, closests[:left])
	game.quicksortFrom(currentID, closests[left+1:])

	return closests
}

func (game *Game) quicksortBary(baryFactories []int) []int {
	if len(baryFactories) < 2 {
		return baryFactories
	}
	left, right := 0, len(baryFactories)-1
	pivotIndex := rand.Int() % len(baryFactories)

	baryFactories[pivotIndex], baryFactories[right] = baryFactories[right], baryFactories[pivotIndex]

	for index := range baryFactories {
		if game.factories[baryFactories[index]].baryDistance <
			game.factories[baryFactories[right]].baryDistance {
			baryFactories[index], baryFactories[left] = baryFactories[left], baryFactories[index]
			left++

		}
	}

	baryFactories[left], baryFactories[right] = baryFactories[right], baryFactories[left]

	game.quicksortBary(baryFactories[:left])
	game.quicksortBary(baryFactories[left+1:])

	return baryFactories
}

// Compute the different combinaisons of orders to compute from
func (game *Game) computeFactoryOrder() []int {
	// orders := make([]int, len(game.factories))

	orderedFactories := make([]int, len(game.theirFactories))
	for i := range game.theirFactories {
		orderedFactories[i] = game.theirFactories[i].id
	}
	orderedFactories = game.quicksortBary(orderedFactories)
	return orderedFactories
}

// Predicts the population of a factory for the next 20 turn
func (game *Game) computeSeer(this *Factory) [20]int {
	var seer [20]int
	for i := range seer {
		if this.owner == -game.playerID {
			seer[i] = -(this.pop + (i+1)*this.prod)
		} else if this.owner == 0 {
			seer[i] = -this.pop
		} else {
			seer[i] = this.pop + (i+1)*this.prod
		}
		for _, troop := range game.troops[this.id] {
			if troop.eta <= i+1 {
				if troop.owner != game.playerID {
					seer[i] -= troop.pop
				} else {
					seer[i] += troop.pop
				}
			}
		}
		for _, bomb := range game.bombs[this.id] {
			if bomb.eta == i-1 {
				if seer[i]/2 < 10 {
					if seer[i] < 10 {
						seer[i] = 0
					} else {
						seer[i] -= 10
					}
				} else {
					seer[i] /= 2
				}
			}
		}
	}
	return seer
}

/*
													FULL METAL TURN
    Just send everything to the closest factory. All the time. Waaargh.
*/
func (game *Game) computeFullMetalTurn() Game {
	resultingTurn := *game

	// Try to increase production of current factories
	for _, mine := range game.myFactories {
		if mine.prod < 3 && mine.pop > 11 {
			resultingTurn.firstMove = append(resultingTurn.firstMove, [4]int{2, mine.id, -1, -1})
			game.factories[mine.id].pop -= 10
		}
	}
	order := game.computeFactoryOrder()
	for _, targetID := range order {
		if game.factories[targetID].prod < 1 {
			continue
		}

		if game.seers[targetID][19] > 0 {
			continue
		}

		totalAttack := 0
		for _, mine := range game.myFactories {
			if mine.pop > 1 {
				attackPower := mine.pop - 1
				totalAttack += attackPower
				resultingTurn.firstMove = append(resultingTurn.firstMove, [4]int{0, mine.id, targetID, attackPower})
				if totalAttack > game.seers[targetID][19] {
					continue
				}
			}
		}
	}
	return resultingTurn
}

/*
													AGRESSIVE TURN
    A strategy based on always attacking the closest factory that still needs
		more forces to take it.
*/
func (game *Game) computeAgressiveTurn() Game {
	resultingTurn := *game

	// Try to increase production of current factories
	for _, mine := range game.myFactories {
		if mine.prod < 3 && mine.pop > 11 {
			resultingTurn.firstMove = append(resultingTurn.firstMove, [4]int{2, mine.id, -1, -1})
			game.factories[mine.id].pop -= 10
		}
	}

	order := game.computeFactoryOrder()
	//fmt.Fprintln(os.Stderr, "Factory orders:", orders)

	// Attack closes factory that has a non-zero production
	for _, targetID := range order {
		if game.factories[targetID].prod < 1 {
			continue
		}
		totalAttackPotential := 0
		var attackStrategy [][2]int

		criticalTurn, populationNeeded := -1, -1
		for index, seer := range game.seers[targetID] {
			if populationNeeded < -seer {
				criticalTurn, populationNeeded = index, -seer
			}
		}
		if criticalTurn == -1 {
			continue
		}

		for _, mineID := range game.proximities[targetID] {
			var mine *Factory
			mine = &game.factories[mineID]
			if mine.owner == game.playerID {
				troopSize := mine.pop
				if populationNeeded > 0 && populationNeeded < troopSize {
					troopSize = populationNeeded + 1
				}
				totalAttackPotential += troopSize
				if troopSize > 0 {
					attackStrategy = append(attackStrategy, [2]int{mineID, troopSize})
					mine.pop -= troopSize
				}
				if populationNeeded < totalAttackPotential {
					break
				}
			}
		}

		// if populationNeeded > totalAttackPotential {
		// 	for _, move := range attackStrategy {
		// 		game.factories[move[0]].pop += move[1]
		// 	}
		// } else {
		for _, move := range attackStrategy {
			// fmt.Fprintln(os.Stderr, "For factory N°", targetID, ", move is as follow:", attackStrategy)
			resultingTurn.firstMove = append(resultingTurn.firstMove, [4]int{0, move[0], targetID, move[1]})
			// fmt.Fprintln(os.Stderr, resultingTurn.firstMove)
		}
		// }
	}
	return resultingTurn
}

/*
													NETWORK TURN
    A strategy based more on a network system, where cyborgs are always
		pushed toward the frontiers of my 'empire', to give more assault
		powers and reinforce the borders.
*/
func (game *Game) computeNetworkTurn() Game {
	resultingTurn := *game

	// Try to increase production of current factories
	for _, mine := range game.myFactories {
		if mine.prod < 3 && mine.pop > 11 {
			resultingTurn.firstMove = append(resultingTurn.firstMove, [4]int{2, mine.id, -1, -1})
			game.factories[mine.id].pop -= 10
		}
	}

	for _, mine := range game.myFactories {
		if mine.pop == 0 {
			continue
		}
		// Check if one of the 3 closest is an ennemy
		closestEnnemy := -1
		var outwardFriends []int
		for _, currentClosest := range game.proximities[mine.id] {
			if game.factories[currentClosest].owner != game.playerID &&
				game.factories[currentClosest].prod > 0 {
				willBecomeMine := false
				for _, seer := range game.seers[currentClosest] {
					if seer > 0 {
						willBecomeMine = true
						break
					}
				}
				if !willBecomeMine {
					closestEnnemy = currentClosest
					break
				}
			} else if mine.baryDistance < game.factories[currentClosest].baryDistance {
				outwardFriends = append(outwardFriends, currentClosest)
				if len(outwardFriends) >= 3 {
					break
				}
			}
		}

		if closestEnnemy != -1 {
			necessaryForce := -game.seers[closestEnnemy][game.distances[mine.id][closestEnnemy]] + 1
			if necessaryForce > 0 {
				if mine.pop > necessaryForce {
					if mine.pop > necessaryForce+10 {
						resultingTurn.firstMove = append(resultingTurn.firstMove, [4]int{0, mine.id, closestEnnemy, necessaryForce + 10})
					} else {
						resultingTurn.firstMove = append(resultingTurn.firstMove, [4]int{0, mine.id, closestEnnemy, necessaryForce})
					}
				} else {
					resultingTurn.firstMove = append(resultingTurn.firstMove, [4]int{0, mine.id, closestEnnemy, mine.pop})
					// mine.pop = 0
				}
			}
		} else {
			for _, currentOutward := range outwardFriends {
				resultingTurn.firstMove = append(resultingTurn.firstMove, [4]int{0, mine.id, currentOutward, mine.pop / 3})
			}
		}
	}

	return resultingTurn
}

/****************************************************
                        MAIN
****************************************************/
func main() {
	// factoryCount: the number of factories
	var factoryCount int
	fmt.Scan(&factoryCount)

	// linkCount: the number of links between factories
	var linkCount int
	fmt.Scan(&linkCount)

	// Creating the Game object
	var actualGame Game
	actualGame.currentTurn = 0
	actualGame.playerID = 1
	actualGame.bombMe = 2
	actualGame.bombOponnent = 2
	actualGame.currentTurn = 0

	// Constructing map
	actualGame.distances = make([][]int, factoryCount)
	allLinks := make([]int, factoryCount*factoryCount)
	for i := range actualGame.distances {
		actualGame.distances[i], allLinks = allLinks[:factoryCount], allLinks[factoryCount:]
	}
	for i := 0; i < linkCount; i++ {
		var factory1, factory2, distance int
		fmt.Scan(&factory1, &factory2, &distance)
		actualGame.distances[factory1][factory2] = distance
		actualGame.distances[factory2][factory1] = distance
	}
	actualGame.initializeProximities()

	for {
		// entityCount: the number of entities (e.g. factories and troops)
		var entityCount int
		fmt.Scan(&entityCount)

		actualGame.factories = make([]Factory, factoryCount)
		actualGame.myFactories = actualGame.myFactories[:0]
		actualGame.theirFactories = actualGame.theirFactories[:0]
		actualGame.troops = make(map[int][]Troop)
		actualGame.bombs = make(map[int][]Bomb)

		for i := 0; i < entityCount; i++ {
			var entityID int
			var entityType string
			var arg1, arg2, arg3, arg4, arg5 int
			fmt.Scan(&entityID, &entityType, &arg1, &arg2, &arg3, &arg4, &arg5)

			if entityType == "FACTORY" {
				if arg1 == 1 {
					actualGame.factories[entityID] = Factory{entityID, arg1, arg2, arg3, arg4, 100}
					actualGame.myFactories = append(actualGame.myFactories, &actualGame.factories[entityID])
				} else {
					actualGame.factories[entityID] = Factory{entityID, arg1, arg2, arg3, arg4, 100}
					actualGame.theirFactories = append(actualGame.theirFactories, &actualGame.factories[entityID])
				}
			} else if entityType == "TROOP" {
				actualGame.troops[arg3] = append(actualGame.troops[arg3], Troop{entityID, arg1, arg2, arg3, arg4, arg5})
			} else if entityType == "BOMB" {
				actualGame.bombs[arg3] = append(actualGame.bombs[arg3], Bomb{entityID, arg1, arg2, arg3, arg4})
			}
		}
		/*************************End of Initialization*****************************/

		bestMove := actualGame.firstMove

		if len(actualGame.myFactories) > 0 {
			actualGame.computeBaryDistances()
			actualGame.computeSeers()
			resultingGame := actualGame.computeFullMetalTurn()
			bestMove = resultingGame.firstMove
		}
		// fmt.Fprintln(os.Stderr, "Debug messages...")

		if len(bestMove) != 0 {
			action := ""
			for i, move := range bestMove {
				if move[0] == 0 {
					action = fmt.Sprint(action, "MOVE ", move[1], move[2], move[3])
				} else if move[0] == 1 {
					action = fmt.Sprint(action, "BOMB ", move[1], move[2])
				} else if move[0] == 2 {
					action = fmt.Sprint(action, "INC ", move[1])
				}
				if i != len(bestMove)-1 {
					action = fmt.Sprint(action, " ; ")
				}
			}
			if actualGame.currentTurn == 0 {
				for _, factory := range actualGame.theirFactories {
					if factory.owner == -1 {
						action = fmt.Sprint(action, "; BOMB ", actualGame.myFactories[0].id, factory.id)
					}
				}
			}
			fmt.Println(action)
		} else {
			// Any valid action, such as "WAIT" or "MOVE source destination cyborgs"
			fmt.Println("WAIT")
		}
		actualGame.currentTurn++
	}
}
