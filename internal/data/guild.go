package data

type Mafia struct {
	SelectedUsersMap map[string][]string
}

type Wink struct {
	CheckedUsers      map[string]bool
	TotalParticipants int
	MessageIDMap      map[string]string
	SelectedUsersMap  map[string][]string
}

type Guild struct {
	ID string

	// embed struct
	Wink
	Mafia
}
