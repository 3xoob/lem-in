package lemin

type Room struct {
	name  string
	x, y  int
	links []*Room
}

type Path struct {
	rooms         []*Room
	numberOfRooms int
}

type AntFarm struct {
	rooms              map[string]*Room
	numberOfAnts       int
	startRoom, endRoom *Room
	edgeCase           bool
}
