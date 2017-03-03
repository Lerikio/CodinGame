package main

import "math/rand"

// Factory represents all the properties of a factory
type Factory struct {
	id               int
	owner            int
	pop              int
	prod             int
	turnsWithoutProd int
	baryToMe         int
	baryToThem       int
}

// Troop represents all the properties of a single troop
type Troop struct {
	owner int
	from  int
	to    int
	pop   int
	eta   int
}

// Bomb represents all the properties of a single bomb
type Bomb struct {
	owner int
	from  int
	to    int
	eta   int
}

// Turn is a structure of everything that is turn-dependant in the game,
// also representing a node in the alpha/beta tree
type Turn struct {
	playerID       int             // 1 for me, -1 for opponent
	bombMe         int             // My number of bombs
	score          int             // Score for current simulation. MyCyborgs - TheirCyborgs
	firstMove      [][4]int        // Move for turn 1 that allowed this configuration
	factories      []Factory       // All the factories in the game
	myFactories    []int           // The Factories I own at this turn
	theirFactories []int           // All other factories. No difference between ennemy and neutral
	troops         map[int][]Troop // All the troops in the game, referenced by their destination
	bombs          map[int][]Bomb  // All the bombs in the game, referenced by their destination
	seers          [][20]int       // Compute all the populations of the game for the 20 next turns
	father         Turn            // The father node
	children       []Turn          // All the children nodes
}

/****************************************************
                GAME object and methods
****************************************************/

// Game is a single object containing all the static data representing the state of
// the game, allowing for easier access and the creation of methods.
type Game struct {
	dists [][]int // Matrix of distances between factory i and factory j
	prox  [][]int // Matrix of the proximities of each factory. Line i is the ordered list of factories closest to factory i
	t0    Turn    // Access to the first turn's data, root of the tree
	start Time
}

// Compute the proximity matrix
func (g *Game) initializeProximities() {
	g.prox = make([][]int, len(g.dists[0]))
	for id := range g.prox {
		closestToID := make([]int, len(g.dists[0])-1)
		count := 0
		for index := range g.prox {
			if index != id {
				closestToID[count] = index
				count++
			}
		}
		closestToID = g.quicksortFrom(id, closestToID)
		g.prox[id] = closestToID[1:]
	}
}

// For each factory, compute the distsance toward the barycenter of allies and ennemies
func (g *Game) computeBaryDistances() {
	for index := range g.t0.factories {
		baryToMe := 0
		baryToThem := 0
		for _, myFactory := range g.t0.myFactories {
			baryToMe += g.dists[index][myFactory]
		}
		for _, theirFactory := range g.t0.theirFactories {
			baryToThem += g.dists[index][theirFactory]
		}
		if len(g.t0.myFactories) == 0 {
			baryToMe = 0
		} else {
			baryToMe /= len(g.t0.myFactories)
		}
		if len(g.t0.theirFactories) == 0 {
			baryToThem = 0
		} else {
			baryToThem /= len(g.t0.theirFactories)
		}
		g.factories[index].baryToMe = baryToMe
		g.factories[index].baryToThem = baryToThem
	}
}

// Implementation of QuickSort in our specific case
func (g *Game) quicksortFrom(currentID int, closests []int) []int {
	if len(closests) < 2 {
		return closests
	}

	left, right := 0, len(closests)-1
	pivotIndex := rand.Int() % len(closests)

	closests[pivotIndex], closests[right] = closests[right], closests[pivotIndex]

	for i := range closests {
		if g.dists[currentID][closests[i]] <
			g.dists[currentID][closests[right]] {
			closests[i], closests[left] = closests[left], closests[i]
			left++
		}
	}

	closests[left], closests[right] = closests[right], closests[left]

	g.quicksortFrom(currentID, closests[:left])
	g.quicksortFrom(currentID, closests[left+1:])

	return closests
}

func (g *Game) quicksortBary(baryFactories []int) []int {
	if len(baryFactories) < 2 {
		return baryFactories
	}
	left, right := 0, len(baryFactories)-1
	pivotIndex := rand.Int() % len(baryFactories)

	baryFactories[pivotIndex], baryFactories[right] = baryFactories[right], baryFactories[pivotIndex]

	for index := range baryFactories {
		if g.factories[baryFactories[index]].baryToMe <
			g.factories[baryFactories[right]].baryToMe {
			baryFactories[index], baryFactories[left] = baryFactories[left], baryFactories[index]
			left++

		}
	}

	baryFactories[left], baryFactories[right] = baryFactories[right], baryFactories[left]

	g.quicksortBary(baryFactories[:left])
	g.quicksortBary(baryFactories[left+1:])

	return baryFactories
}
