package lemin

type Room struct {
	name  string
	x, y  int
	links []*Room
}

type Path struct {
	rooms         []*Room
	RoomsNum int
}

type AntFarm struct {
	rooms              map[string]*Room
	AntsNum       int
	start, end *Room
	edgeCase           bool
}
