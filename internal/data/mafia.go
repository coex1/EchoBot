package data

type Mafia struct {
	SelectedUsers map[string][]string // GuildID : []UserID

	MafiaList   []string
	PoliceList  []string
	DoctorList  []string
	CitizenList []string
}
