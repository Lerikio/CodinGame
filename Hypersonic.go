package main

import "fmt"
import "os"
import "strconv"
import "math"
import "time"

const SEER = 9 // NEEDS to be strictly over 8

/*************** COORDINATES *****************/
type Coordinates struct{
    x int
    y int
}

/*************** PLAYER *****************/
type Player struct{
    coord Coordinates
    id int
    bombsAvailable int
    currentRange int
}

/*************** BOMB *****************/
type Bomb struct{
    coord Coordinates
    id int
    explodesIn int
    explosionRange int
}

/*************** BOOLGRID *****************/
type BoolGrid [11][13]bool
func (grid *BoolGrid) elem(x int, y int) *bool{ return &grid[y][x] }
func (grid *BoolGrid) elemC(coord Coordinates) *bool{ return &grid[coord.y][coord.x] }

/*************** COORDLIST *****************/
type CoordList []Coordinates

/*************** GAME *****************/
type Game struct{
    boxes [SEER]BoolGrid
    walls BoolGrid
    safe [SEER]BoolGrid
    players [4]Player
    bombs [SEER][]Bomb
    bonuses CoordList
    myID int
    directions []Coordinates
}

func (p *Game) removeBomb(i, index int) {
	s := *p
	s.bombs[i] = append(s.bombs[i][:index], s.bombs[i][index+1:]...) // perfectly fine if i is the last element
	*p = s
}
/*************** PATH *****************/
type Path struct{
    coord Coordinates
    father *Path
    score float64
    turn int
    sons [5]*Path
    game *Game
}

func (node *Path) buildTree(lastPaths *[]*Path) {
    if node.turn < SEER - 1 {
        for index, direction := range(append((*node.game).directions, Coordinates{0, 0})) {
            nextCoord := Coordinates{node.coord.x + direction.x, node.coord.y + direction.y}
            if nextCoord.x >= 0 && nextCoord.x < 13 && nextCoord.y >= 0 && nextCoord.y < 11 {
                if *(node.game).safe[node.turn + 1].elemC(nextCoord) {
                    node.sons[index] = &Path{ coord: nextCoord,
                                        father: node,
                                        score: node.score,
                                        turn: node.turn + 1,
                                        game: node.game }
                    (*node.sons[index]).computeScore()
                    (*node.sons[index]).buildTree(lastPaths)
                }
            }
        }
    } else {
        *lastPaths = append(*lastPaths, node)
    }
}

func (node *Path) computeScore() {
    var testCase Coordinates
    for _, direction := range((*node.game).directions){
        for step := 1; step < node.game.players[node.game.myID].currentRange; step ++ {
            testCase = Coordinates{node.coord.x + step * direction.x, node.coord.y + step * direction.y}
            if testCase.x >= 0 && testCase.x < 13 && testCase.y >= 0 && testCase.y < 11 {
                if *(node.game).walls.elemC(testCase){
                    break
                } else if *(node.game).boxes[SEER - 1].elemC(testCase) && node.father != nil && node.father.coord != node.coord{
                    node.score += math.Pow(0.99, float64(node.turn))
                    break
                }
            } else {
                break
            }
        }
    }

    for _, bonus := range(node.game.bonuses) {
        if bonus == node.coord && node.father != nil  && node.father.coord != node.coord{
            node.score += math.Pow(0.7, float64(node.turn)) * 3
            break
        }
    }

    for t := SEER - 1; t >= 0; t -- {
        if !*(node.game).safe[t].elemC(node.coord){
            node.score += math.Pow(0.8, float64(node.turn)) * -1
        }
    }
}

/*************** FUNCTIONS *****************/
func min(a int, b int) int{
    if a < b {
        return a
    } else {
        return b
    }
}

func max(a int, b int) int{
    if a > b {
        return a
    } else {
        return b
    }
}

func distE(a Coordinates, b Coordinates) int{
    return int(math.Abs(float64(a.x - b.x)) + math.Abs(float64(a.y - b.y)))
}

func addBombToSeer(game *Game , bomb *Bomb, turn int) {
    var testCase Coordinates
    for _, direction := range(game.directions){
        //fmt.Fprintf(os.Stderr, "Testing direction around bombe...")
        for step := 0; step < bomb.explosionRange; step ++ {
            testCase = Coordinates{bomb.coord.x + step * direction.x, bomb.coord.y + step * direction.y}
            if testCase.x >= 0 && testCase.x < 13 && testCase.y >= 0 && testCase.y < 11 {
                if *game.walls.elemC(testCase){
                    break
                } else if *game.boxes[turn].elemC(testCase) {
                    for previousTurn := max(turn + 1, 0); previousTurn < SEER; previousTurn ++ {
                        *game.boxes[previousTurn].elemC(testCase) = false
                    }
                    break
                } else {
                    for t := turn + 2; t < SEER ; t ++ {
                        //fmt.Fprintf(os.Stderr, "Turntest: " + strconv.Itoa(t) + " | ")
                        //fmt.Fprintf(os.Stderr, "Number of bombs: " + strconv.Itoa(len(game.bombs[t])) + "\n")
                        for index, bombF := range(game.bombs[t]) {
                            if bombF.coord == testCase {
                                addBombToSeer(game, &bombF, turn)
                                game.removeBomb(t, index)
                            }
                        }
                    }
                    *game.safe[turn].elemC(testCase) = false
                }
            } else {
                break
            }
        }
    }
    for i := 0; i <= turn; i ++ {
        *game.safe[turn].elemC(bomb.coord) = false
    }
}

func main() {

    

/***********************************************
        Initialization for Inputs
***********************************************/
    var width, height, myId int
    fmt.Scan(&width, &height, &myId)
    var row string

    directions := []Coordinates{Coordinates{1, 0}, Coordinates{-1, 0}, Coordinates{0, 1}, Coordinates{0, -1}}

/***********************************************
        Other Inits
***********************************************/
    game := Game{myID: myId, directions: directions}

    for {
        start := time.Now()
/***********************************************
        Turn based initialization
***********************************************/
        game.safe = *new([SEER]BoolGrid)
        game.bombs = *new([SEER][]Bomb)
        game.bonuses = nil

/***********************************************
        Gathering inputs
***********************************************/
        for y := 0; y < height; y++ {
            fmt.Scan(&row)
            for x, value := range(row) {
                for turn := 0; turn < SEER; turn ++ {
                    if string(value) != "." && string(value) != "X"{
                        *game.boxes[turn].elem(x, y) = true
                    } else {
                        *game.boxes[turn].elem(x, y) = false
                    }
                }
                if string(value) == "X" {
                    *game.walls.elem(x, y) = true
                } else {
                    *game.walls.elem(x, y) = false
                }
            }
        }
        var entities int
        fmt.Scan(&entities)

        for i := 0; i < entities; i++ {
            var entityType, owner, x, y, param1, param2 int
            fmt.Scan(&entityType, &owner, &x, &y, &param1, &param2)

            if entityType == 0 {
                game.players[owner] = Player{Coordinates{x, y}, owner, param1, param2}
            } else if entityType == 1 {
                game.bombs[param1] = append(game.bombs[param1], Bomb{Coordinates{x, y}, owner, param1, param2})
            } else if entityType == 2 {
                game.bonuses = append(game.bonuses, Coordinates{x, y} )
            }
        }


/***********************************************
        Starting analysis
***********************************************/

        action := "Move "
        var lastNodes []*Path
        // var computeNoBomb bool

        father := Path{ score: -1, turn: 0,
                        game: &game, father: nil,
                        sons: *new([5]*Path), coord: game.players[myId].coord}
        father.computeScore()
        maxPath := &father
        maxScore := float64(-1)

        if father.score > 0 && game.players[myId].bombsAvailable > 0 {

            gameBomb := game

            bomb := Bomb{   coord: game.players[myId].coord,
                            id: myId,
                            explodesIn: 8,
                            explosionRange: game.players[myId].currentRange }

            gameBomb.bombs[7] = append(gameBomb.bombs[7], bomb)

            for turn := 0; turn < SEER; turn ++ {
                currentSafe := &gameBomb.safe[turn]
                *currentSafe = *new(BoolGrid)
                for y := 0; y < height; y ++ {
                    for x := 0; x < width; x ++ {
                        if !*gameBomb.boxes[turn].elem(x, y) && !*gameBomb.walls.elem(x, y) {
                            *(*currentSafe).elem(x, y) = true
                        }
                        if turn < SEER -1 {
                            for index, bomb := range(gameBomb.bombs[turn + 1]){
                                addBombToSeer(&gameBomb, &bomb, turn)
                                gameBomb.bombs[turn] = append(gameBomb.bombs[turn][:index], gameBomb.bombs[turn][index+1:]...)
                            }
                        }
                    }
                }
            }

            fatherBombNow :=  Path{ score: -1, turn: 0,
                                    game: &gameBomb, father: nil,
                                    sons: *new([5]*Path), coord: gameBomb.players[myId].coord}
            fatherBombNow.computeScore()
            fatherBombNow.buildTree(&lastNodes)

            maxPath = &fatherBombNow
            maxScore = -1
            for _, path := range(lastNodes) {
                if path.score > maxScore {
                    maxPath = path
                    maxScore = path.score
                }
            }

            if maxScore != -1 && fatherBombNow.score > 0{
                action = "BOMB "
            } else {
                action = "MOVE "
            }
        } else {
        //if game.players[myId].bombsAvailable == 0 || computeNoBomb || father.score == 0{
            for turn := 0; turn < SEER; turn ++ {
                currentSafe := &game.safe[turn]
                *currentSafe = *new(BoolGrid)
                for y := 0; y < height; y ++ {
                    for x := 0; x < width; x ++ {
                        if !*game.boxes[turn].elem(x, y) && !*game.walls.elem(x, y) {
                            *(*currentSafe).elem(x, y) = true
                        }
                        if turn < SEER -1 {
                            for len(game.bombs[turn + 1]) > 0 { 
                                bomb := game.bombs[turn + 1][0]
                                fmt.Fprintf(os.Stderr, "Ajout de bombe ! " + strconv.Itoa(turn) + " " + strconv.Itoa(0))
                                addBombToSeer(&game, &bomb, turn + 1)
                               game.removeBomb(turn + 1, 0)
                            }
                        }
                    }
                }
            }
            lastNodes = nil
            father.buildTree(&lastNodes)
            maxPath = &father
            maxScore = -1
            for _, path := range(lastNodes) {
                if path.score > maxScore {
                    maxPath = path
                    maxScore = path.score
                }
            }

            action = "MOVE "
        }

        var pathPrinting []Coordinates
        var scorePrinting []int
        fmt.Fprintf(os.Stderr, strconv.Itoa(int(100 * maxScore)) + " : ")
        for maxPath != nil && maxPath.father != nil && maxPath.turn > 1{
            maxPath = maxPath.father
            pathPrinting = append(pathPrinting, maxPath.coord)
            scorePrinting = append(scorePrinting, int(100*maxPath.score))
            //fmt.Fprintf(os.Stderr, "[" + strconv.Itoa(maxPath.coord.x) + ", " + strconv.Itoa(maxPath.coord.y) + "]<-")
        }
        for i := len(pathPrinting) - 1; i > 0; i -- {
            fmt.Fprintf(os.Stderr, "-> " + strconv.Itoa(scorePrinting[i]) + " [" + strconv.Itoa(pathPrinting[i].x) + ", " + strconv.Itoa(pathPrinting[i].y) + "] ")
        }

        fmt.Println( action + strconv.Itoa(maxPath.coord.x) + " " + strconv.Itoa(maxPath.coord.y) )

        elapsed := time.Since(start)
        fmt.Fprintf(os.Stderr, "\n%s\n", elapsed)
    }
}

