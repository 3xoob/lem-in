package lemin

import "sort"

func BFS(Farm *AntFarm) []*Path {
	var paths []*Path
	waiting := []Path{
		{
			rooms:         []*Room{Farm.start},
			RoomsNum: 1,
		},
	}
	visited := map[string]bool{Farm.start.name: true}
	for len(waiting) != 0 {
		path := waiting[0]
		waiting = waiting[1:]
		lastRoom := path.rooms[len(path.rooms)-1]
		if lastRoom == Farm.end {
			visited = map[string]bool{Farm.start.name: true}
			paths = append(paths, &path)
			for _, path := range paths {
				for _, room := range path.rooms {
					visited[room.name] = true
				}
			}
			waiting = []Path{
				{
					rooms:         []*Room{Farm.start},
					RoomsNum: 1,
				},
			}
			continue
		}
		for _, link := range lastRoom.links {
			if !visited[link.name] || (link == Farm.end && lastRoom != Farm.start) {
				newPathRooms := make([]*Room, len(path.rooms))
				copy(newPathRooms, path.rooms)
				newPathRooms = append(newPathRooms, link)
				newPath := Path{
					rooms:         newPathRooms,
					RoomsNum: len(newPathRooms),
				}
				waiting = append(waiting, newPath)
				visited[link.name] = true
				if Farm.edgeCase {
					break
				}
			}
		}
	}
	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i].RoomsNum < paths[j].RoomsNum
	})
	return paths
}
