package main

import "fmt"
import "sort"
import "os"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

type Factory struct{
    id int
    owner int
    pop int
    prod int
}

type Troop struct{
    id int
    owner int
    from int
    to int
    pop int
    eta int
}

func main() {
    // factoryCount: the number of factories
    var factoryCount int
    fmt.Scan(&factoryCount)

    // linkCount: the number of links between factories
    var linkCount int
    fmt.Scan(&linkCount)

    // Consturcting map
    factoryMap := make([][]int, factoryCount)
    allLinks := make([]int, factoryCount*factoryCount)
    for i := range factoryMap {
        factoryMap[i], allLinks = allLinks[:factoryCount], allLinks[factoryCount:]
    }

    for i := 0; i < linkCount; i++ {
        var factory1, factory2, distance int
        fmt.Scan(&factory1, &factory2, &distance)
        factoryMap[factory1][factory2] = distance
        factoryMap[factory2][factory1] = distance
    }
    for {
        // entityCount: the number of entities (e.g. factories and troops)
        var entityCount int
        fmt.Scan(&entityCount)

        myFactories := make(map[int]Factory)
        theirFactories := make(map[int]Factory)
        //neutralFactories := make(map[int]Factory)

        myTroops := make(map[int][]Troop)
        theirTroops := make(map[int][]Troop)

        for i := 0; i < entityCount; i++ {
            var entityId int
            var entityType string
            var arg1, arg2, arg3, arg4, arg5 int
            fmt.Scan(&entityId, &entityType, &arg1, &arg2, &arg3, &arg4, &arg5)

            if entityType == "FACTORY" {
                fmt.Fprintln(os.Stderr, "Factory number:", entityId)
                if arg1 == 1 {
                    myFactories[entityId] = Factory{entityId, arg1, arg2, arg3}
                } else {
                    theirFactories[entityId] = Factory{entityId, arg1, arg2, arg3}
                //} else {
                //    neutralFactories[entityId] = Factory{entityId, arg1, arg2, arg3}
                }
            } else if entityType == "TROOP" {
                if arg1 == 1 {
                    myTroops[arg2] = append(myTroops[arg2], Troop{entityId, arg1, arg2, arg3, arg4, arg5})
                } else if arg1 == -1 {
                    theirTroops[arg2] = append(myTroops[arg2], Troop{entityId, arg1, arg2, arg3, arg4, arg5})
                }
            }
        }

        var bestGuess [][3]int

        for _, base := range myFactories {
            links := make([]int, factoryCount)
            copy(links, factoryMap[base.id])

            mapDistances := make(map[int]int)
            for index, value := range links {
                mapDistances[value] = index
            }

            sort.Ints(links)
            //fmt.Fprintln(os.Stderr, "Factory number:", base.id)

            for _, distance := range links {

                otherId := mapDistances[distance]

                // check if mine
                _, mine := myFactories[otherId]
                if mine {
                    //fmt.Fprintln(os.Stderr, "Is Mine:", otherId)
                    continue
                }

                // Check if possible to invade
                other := theirFactories[otherId]
                if other.pop >= base.pop - 2 {
                    //fmt.Fprintln(os.Stderr, "Can't invade:", otherId)
                    continue
                }

                // if not mine, check if troops going there from here already
                troopGoingThere := false
                for _, troop := range myTroops[base.id] {
                    if troop.to == otherId {
                        troopGoingThere = true
                        break
                    }
                }
                if troopGoingThere {
                    //fmt.Fprintln(os.Stderr, "Already Invading:", otherId)
                    continue
                }

                //fmt.Fprintln(os.Stderr, "Best guess!", otherId, "with", other.pop + 1)
                bestGuess = append(bestGuess, [3]int{base.id, otherId, other.pop + 1})
                base.pop -= other.pop + 1
                //break
            }

            //if bestGuess[0] != -1 { break }
        }

        // fmt.Fprintln(os.Stderr, "Debug messages...")

         if len(bestGuess) != 0  {
            action := ""
            for i, guess := range bestGuess{
                if i == len(bestGuess)-1 {
                    action = fmt.Sprint(action, "MOVE ", guess[0], guess[1], guess[2])
                } else {
                    action = fmt.Sprint(action, "MOVE ", guess[0], guess[1], guess[2], " ; ")
                }

            }
            fmt.Println(action)
         } else {
            // Any valid action, such as "WAIT" or "MOVE source destination cyborgs"
            fmt.Println("WAIT")
         }


    }
}
