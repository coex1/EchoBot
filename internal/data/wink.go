package data

type Wink struct {
	CheckedUsers      map[string]bool
	TotalParticipants int
	MessageIDMap      map[string]string
	SelectedUsersMap  map[string][]string
}
