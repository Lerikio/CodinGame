package main

import "fmt"
import "os"
import "math"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

type Point struct{
    x int
    y int
}

type Target []Point

// Computes the gravitational center of the Target zone
func (t *Target) center() Point {
    var n int
    var c Point

    for _, point := range *t {
        n += 1
        c.x += point.x
        c.y += point.y
    }

    c.x = c.x / n
    c.y = c.y / n

    return c
}


// Updates the target with the knowledge of one point being closer (c) to the bomb than the other (f).
func (t *Target) update(c Point, f Point, now Point) Point{
    var newT Target
    var maxDist float64
    var furthestPoint Point

    for _, point := range *t {
        if dist(point, c) <= dist(point, f) {
            newT = append(newT, point)
             if dist(now, point) >= maxDist {
                furthestPoint = point
                maxDist = dist(now, furthestPoint)
            }
        }
    }

    *t = newT

    return furthestPoint
}

// Computes the imperfect symmetry of one point in respect to a center.
func (p *Point) symmmetry(c Point, width int, height int) Point {
    xDist := c.x - p.x
    yDist := c.y - p.y

    var newP Point
    newP.x = c.x + xDist + 1
    newP.y = c.y + yDist

    if newP.x > width - 1 {
        newP.x = width - 1
        newP.y += 1
    } else if newP.x < 0 {
        newP.x = 0
        newP.y -= 1
    }

    if newP.y > height - 1 {
        newP.y = height - 1
    } else if newP.y < 0 {
        newP.y = 0
    }

    return newP
}

// Computes the distance between two points the same way CG does
func dist(a Point, b Point) float64 {
    var xDist, yDist float64

    xDist = float64(b.x - a.x)
    yDist = float64(b.y - a.y)

    return math.Sqrt(xDist*xDist + yDist*yDist)
}

func main() {
    // W: width of the building.
    // H: height of the building.
    var W, H int
    fmt.Scan(&W, &H)

    // N: maximum number of turns before game over.
    var N int
    fmt.Scan(&N)

    var nowPos Point
    fmt.Scan(&nowPos.x, &nowPos.y)
    pastPos := nowPos

    var target Target
    var turn int

    for {
        // bombDir: Current distance to the bomb compared to previous distance (COLDER, WARMER, SAME or UNKNOWN)
        var bombDir string
        fmt.Scan(&bombDir)
        var furthest Point
        var newPosition Point
        fmt.Fprintln(os.Stderr, N)


        if bombDir == "COLDER" {
            furthest = target.update(pastPos, nowPos, nowPos)
            newPosition = furthest
        } else if bombDir == "WARMER" {
            furthest = target.update(nowPos, pastPos, nowPos)
            gravitationalCenter := target.center()
            newPosition = nowPos.symmmetry(gravitationalCenter, W, H)
        } else if bombDir == "SAME" {
                var newT Target
                var maxDist float64
                for _, point := range target {
                    if dist(point, nowPos) == dist(point, pastPos) {
                        newT = append(newT, point)
                        if dist(nowPos, point) >= maxDist {
                            furthest = point
                            maxDist = dist(nowPos, furthest)
                        }
                    }

                }
                target = newT
                newPosition = furthest
        } else if bombDir == "UNKNOWN" {
            // First turn : populate Target
            for i := 0; i < W; i ++ {
                for j := 0; j < H; j ++ {
                    target = append(target, Point{i, j})
                }
            }
        }

        /*

        if turn > N/2 {
            newPosition = furthest
        } else {
            gravitationalCenter := target.center()
            newPosition = nowPos.symmmetry(gravitationalCenter, W, H)
        }
*/
        pastPos = nowPos
        nowPos = newPosition

        turn += 1

        // fmt.Fprintln(os.Stderr, "Debug messages...")
        fmt.Println(newPosition.x, newPosition.y)// Write action to stdout
    }
}