package main

import "fmt"

import "os"

// Factory represents all the properties of a factory
type Factory struct {
	id               int
	owner            int
	pop              int
	prod             int
	turnsWithoutProd int
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
	id         int
	owner      int
	from       int
	to         int
	explodesIn int
}

// Game is a single object containing all the data representing the state of
// the game, allowing for easier access and the creation of methods.
type Game struct {
	distances      [][]int         // Matrix of distances between factory i and factory j
	factories      []Factory       // All the factories in the game
	myFactories    []*Factory      // The Factories I own at this turn
	theirFactories []*Factory      // All other factories. No difference between ennemy and neutral
	troops         map[int][]Troop // All the troops in the game, referenced by their destination
	bombs          map[int][]Bomb  // All the bombs in the game, referenced by their destination
}

func main() {
	// factoryCount: the number of factories
	var factoryCount int
	fmt.Scan(&factoryCount)

	// linkCount: the number of links between factories
	var linkCount int
	fmt.Scan(&linkCount)

	// Creating the Game object
	var game Game

	// Constructing map
	game.distances = make([][]int, factoryCount)
	allLinks := make([]int, factoryCount*factoryCount)
	for i := range game.distances {
		game.distances[i], allLinks = allLinks[:factoryCount], allLinks[factoryCount:]
	}
	for i := 0; i < linkCount; i++ {
		var factory1, factory2, distance int
		fmt.Scan(&factory1, &factory2, &distance)
		game.distances[factory1][factory2] = distance
		game.distances[factory2][factory1] = distance
	}

	for {
		// entityCount: the number of entities (e.g. factories and troops)
		var entityCount int
		fmt.Scan(&entityCount)

		game.factories = make([]Factory, factoryCount)
		game.troops = make(map[int][]Troop)

		for i := 0; i < entityCount; i++ {
			var entityID int
			var entityType string
			var arg1, arg2, arg3, arg4, arg5 int
			fmt.Scan(&entityID, &entityType, &arg1, &arg2, &arg3, &arg4, &arg5)

			if entityType == "FACTORY" {
				fmt.Fprintln(os.Stderr, "Factory number:", entityID)
				if arg1 == 1 {
					game.factories[entityID] = Factory{entityID, arg1, arg2, arg3, arg4}
					game.myFactories = append(game.myFactories, &game.factories[entityID])
				} else {
					game.factories[entityID] = Factory{entityID, arg1, arg2, arg3, arg4}
					game.theirFactories = append(game.theirFactories, &game.factories[entityID])
					//} else {
					//    neutralFactories[entityID] = Factory{entityID, arg1, arg2, arg3}
				}
			} else if entityType == "TROOP" {
				game.troops[arg3] = append(game.troops[arg3], Troop{entityID, arg1, arg2, arg3, arg4, arg5})
			} else if entityType == "BOMB" {
				game.bombs[arg3] = append(game.bombs[arg3], Bomb{entityID, arg1, arg2, arg3, arg4})
			}
		}

		var bestMove [][3]int

		// fmt.Fprintln(os.Stderr, "Debug messages...")

		if len(bestMove) != 0 {
			action := ""
			for i, move := range bestMove {
				if i == len(bestMove)-1 {
					action = fmt.Sprint(action, "MOVE ", move[0], move[1], move[2])
				} else {
					action = fmt.Sprint(action, "MOVE ", move[0], move[1], move[2], " ; ")
				}

			}
			fmt.Println(action)
		} else {
			// Any valid action, such as "WAIT" or "MOVE source destination cyborgs"
			fmt.Println("WAIT")
		}

	}
}
