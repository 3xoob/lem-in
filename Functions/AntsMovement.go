package lemin

import "fmt"

func AntsMovement(Paths []*Path, antFarm *AntFarm) {
	antPath := make(map[int]int)
	antPosition := make(map[int]int)
	InPath := make([]int, len(Paths))
	for antID := 1; antID <= antFarm.numberOfAnts; antID++ {
		pathIndex := 0
		minCost := InPath[0] + Paths[0].numberOfRooms
		for i := 1; i < len(Paths); i++ {
			cost := Paths[i].numberOfRooms + InPath[i]
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
	var output string

	for step := 1; ; step++ {
		antMoving = false
		for pathIndex := 0; pathIndex < len(Paths); pathIndex++ {
			for j := 0; j < len(antsInside[pathIndex]); j++ {
				if antPosition[antsInside[pathIndex][j]] < Paths[pathIndex].numberOfRooms-1 {
					antMoving = true
					antPosition[antsInside[pathIndex][j]]++
					output += "L"
					output += fmt.Sprint(antsInside[pathIndex][j])
					output += "-"
					output += Paths[pathIndex].rooms[antPosition[antsInside[pathIndex][j]]].name
					output += " "

				}
			}
		}
		for pathIndex := 0; pathIndex < len(Paths); pathIndex++ {
			for len(antsOutside[pathIndex]) != 0 {
				if antPosition[antsOutside[pathIndex][0]] < Paths[pathIndex].numberOfRooms-1 {
					antMoving = true
					antPosition[antsOutside[pathIndex][0]]++
					output += "L"
					output += fmt.Sprint(antsOutside[pathIndex][0])
					output += "-"
					output += Paths[pathIndex].rooms[antPosition[antsOutside[pathIndex][0]]].name
					output += " "
					antID := antsOutside[pathIndex][0]
					antsInside[pathIndex] = append(antsInside[pathIndex], antID)
					antsOutside[pathIndex] = antsOutside[pathIndex][1:]
					break
				}
			}
		}
		if !antMoving {
			break
		}
		fmt.Println(output)
		output = ""
	}
}
