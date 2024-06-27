package lemin

import "sort"

func BFS(Farm *AntFarm) []*Path {
	var paths []*Path
	waiting := []Path{
		{
			rooms:         []*Room{Farm.startRoom},
			numberOfRooms: 1,
		},
	}
	visited := map[string]bool{Farm.startRoom.name: true}
	for len(waiting) != 0 {
		path := waiting[0]
		waiting = waiting[1:]
		lastRoom := path.rooms[len(path.rooms)-1]
		if lastRoom == Farm.endRoom {
			visited = map[string]bool{Farm.startRoom.name: true}
			paths = append(paths, &path)
			for _, path := range paths {
				for _, room := range path.rooms {
					visited[room.name] = true
				}
			}
			waiting = []Path{
				{
					rooms:         []*Room{Farm.startRoom},
					numberOfRooms: 1,
				},
			}
			continue
		}
		for _, link := range lastRoom.links {
			if !visited[link.name] || (link == Farm.endRoom && lastRoom != Farm.startRoom) {
				newPathRooms := make([]*Room, len(path.rooms))
				copy(newPathRooms, path.rooms)
				newPathRooms = append(newPathRooms, link)
				newPath := Path{
					rooms:         newPathRooms,
					numberOfRooms: len(newPathRooms),
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
		return paths[i].numberOfRooms < paths[j].numberOfRooms
	})
	return paths
}
