package lemin

import (
	"fmt"
	"strings"
)

// AntsMovement assigns each ant to a path (greedily minimizing, at each
// assignment, the projected turn on which that ant would finish) and
// simulates the turn-by-turn movement. It returns one string per turn,
// each a space-separated list of "Lant-room" moves.
func AntsMovement(Paths []*Path, antFarm *AntFarm) []string {
	antPath := make(map[int]int)
	antPosition := make(map[int]int)
	InPath := make([]int, len(Paths))
	for antID := 1; antID <= antFarm.AntsNum; antID++ {
		pathIndex := 0
		minCost := InPath[0] + Paths[0].RoomsNum
		for i := 1; i < len(Paths); i++ {
			cost := Paths[i].RoomsNum + InPath[i]
			if minCost > cost {
				minCost = cost
				pathIndex = i
			}
		}
		antPath[antID] = pathIndex
		InPath[pathIndex]++
	}

	antsOutside := make(map[int][]int)
	for i := 0; i < len(Paths); i++ {
		antsOutside[i] = make([]int, 0)
	}
	for i := 1; i <= len(antPath); i++ {
		antsOutside[antPath[i]] = append(antsOutside[antPath[i]], i)
	}

	antsInside := make(map[int][]int)
	var antMoving bool
	var turns []string

	for {
		var moves []string
		antMoving = false
		for pathIndex := 0; pathIndex < len(Paths); pathIndex++ {
			for j := 0; j < len(antsInside[pathIndex]); j++ {
				ant := antsInside[pathIndex][j]
				if antPosition[ant] < Paths[pathIndex].RoomsNum-1 {
					antMoving = true
					antPosition[ant]++
					moves = append(moves, fmt.Sprintf("L%d-%s", ant, Paths[pathIndex].rooms[antPosition[ant]].name))
				}
			}
		}
		for pathIndex := 0; pathIndex < len(Paths); pathIndex++ {
			for len(antsOutside[pathIndex]) != 0 {
				ant := antsOutside[pathIndex][0]
				if antPosition[ant] < Paths[pathIndex].RoomsNum-1 {
					antMoving = true
					antPosition[ant]++
					moves = append(moves, fmt.Sprintf("L%d-%s", ant, Paths[pathIndex].rooms[antPosition[ant]].name))
					antsInside[pathIndex] = append(antsInside[pathIndex], ant)
					antsOutside[pathIndex] = antsOutside[pathIndex][1:]
					break
				}
			}
		}
		if !antMoving {
			break
		}
		turns = append(turns, strings.Join(moves, " "))
	}
	return turns
}
