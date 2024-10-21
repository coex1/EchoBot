package data

type Wink struct {
  // users selected to play the game
	SelectedUsers     map[string][]string

	CheckedUsers      map[string]bool
	TotalParticipants int
	MessageIDMap      map[string]string
}
