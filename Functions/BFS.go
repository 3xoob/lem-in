package lemin

import "sort"

// flowNode splits each intermediate room into an "in" and "out" side joined
// by a capacity-1 edge, so no room (other than start/end) can be shared by
// more than one path. Start and end are left unsplit since ants may pass
// through them without the one-ant-per-room restriction.
type flowNode struct {
	room *Room
	part byte // 'S' unsplit (start/end), 'I' in-side, 'O' out-side
}

func inNode(farm *AntFarm, r *Room) flowNode {
	if r == farm.start || r == farm.end {
		return flowNode{r, 'S'}
	}
	return flowNode{r, 'I'}
}

func outNode(farm *AntFarm, r *Room) flowNode {
	if r == farm.start || r == farm.end {
		return flowNode{r, 'S'}
	}
	return flowNode{r, 'O'}
}

// BFS finds the maximum set of room-disjoint paths from start to end, using
// Edmonds-Karp max-flow (BFS augmenting paths keep each path as short as
// possible, which also minimizes the total path length of the resulting
// set) followed by a flow decomposition into individual paths.
func BFS(Farm *AntFarm) []*Path {
	capacity := make(map[flowNode]map[flowNode]int)
	adjacency := make(map[flowNode][]flowNode)

	addEdge := func(u, v flowNode) {
		if capacity[u] == nil {
			capacity[u] = make(map[flowNode]int)
		}
		if capacity[v] == nil {
			capacity[v] = make(map[flowNode]int)
		}
		if _, exists := capacity[u][v]; !exists {
			adjacency[u] = append(adjacency[u], v)
			capacity[u][v] = 0
		}
		capacity[u][v]++
		if _, exists := capacity[v][u]; !exists {
			adjacency[v] = append(adjacency[v], u)
			capacity[v][u] = 0
		}
	}

	roomNames := make([]string, 0, len(Farm.rooms))
	for name := range Farm.rooms {
		roomNames = append(roomNames, name)
	}
	sort.Strings(roomNames) // deterministic edge construction/iteration order

	for _, name := range roomNames {
		r := Farm.rooms[name]
		if r != Farm.start && r != Farm.end {
			addEdge(inNode(Farm, r), outNode(Farm, r))
		}
	}
	for _, name := range roomNames {
		r := Farm.rooms[name]
		for _, n := range r.links {
			addEdge(outNode(Farm, r), inNode(Farm, n))
		}
	}

	origCap := make(map[flowNode]map[flowNode]int, len(capacity))
	for u, m := range capacity {
		origCap[u] = make(map[flowNode]int, len(m))
		for v, c := range m {
			origCap[u][v] = c
		}
	}

	start, end := inNode(Farm, Farm.start), inNode(Farm, Farm.end)

	for {
		prev := make(map[flowNode]flowNode)
		visited := map[flowNode]bool{start: true}
		queue := []flowNode{start}
		found := false
		for i := 0; i < len(queue) && !found; i++ {
			cur := queue[i]
			for _, next := range adjacency[cur] {
				if capacity[cur][next] > 0 && !visited[next] {
					visited[next] = true
					prev[next] = cur
					if next == end {
						found = true
						break
					}
					queue = append(queue, next)
				}
			}
		}
		if !found {
			break
		}
		path := []flowNode{end}
		for cur := end; cur != start; {
			p := prev[cur]
			path = append(path, p)
			cur = p
		}
		for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
			path[i], path[j] = path[j], path[i]
		}
		for i := 0; i < len(path)-1; i++ {
			u, v := path[i], path[i+1]
			capacity[u][v]--
			capacity[v][u]++
		}
	}

	// Net flow on each real edge = original capacity - remaining residual.
	remaining := make(map[flowNode][]flowNode)
	for u, neighbors := range adjacency {
		for _, v := range neighbors {
			if origCap[u][v] > 0 && origCap[u][v]-capacity[u][v] > 0 {
				remaining[u] = append(remaining[u], v)
			}
		}
	}

	var paths []*Path
	for len(remaining[start]) > 0 {
		var roomPath []*Room
		cur := start
		for cur != end {
			next := remaining[cur][len(remaining[cur])-1]
			remaining[cur] = remaining[cur][:len(remaining[cur])-1]
			if len(roomPath) == 0 || roomPath[len(roomPath)-1] != cur.room {
				roomPath = append(roomPath, cur.room)
			}
			cur = next
		}
		if len(roomPath) == 0 || roomPath[len(roomPath)-1] != end.room {
			roomPath = append(roomPath, end.room)
		}
		paths = append(paths, &Path{rooms: roomPath, RoomsNum: len(roomPath)})
	}

	sort.SliceStable(paths, func(i, j int) bool {
		return paths[i].RoomsNum < paths[j].RoomsNum
	})
	return paths
}
