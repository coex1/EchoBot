package data

// game states
const (
  NONE = iota + 0
  INITIATED
  IN_PROGRESS
  LAST_PLAYER
  ENDED
)

type Wink struct {
  // game state
  State               int

  // every possible player information
  IDList              map[string]string
  NameList            map[string]string
  MaxPossiblePlayers  int

  // users selected to play the game
	SelectedUsersID     []string

  // all selected detail

  // users that have confirmed their target
  // ID -> BOOL
	ConfirmedUsers    map[string]bool
  ConfirmedCount    int

  // king's ID
  KingID string

  // final player
  FinalPlayerID string

  // ID -> BOOL
	CheckedUsers      map[string]bool
	TotalParticipants int

  // ID -> ID
	UserSelection     map[string]string
	UserSelectionFinal     map[string]string
}
