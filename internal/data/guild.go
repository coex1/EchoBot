package data

type Guild struct {
	ID string

	// embed struct
	Wink
	Mafia
}

func Initialize(g *Guild){
  g.Wink.SelectedUsers = make(map[string][]string)
}
