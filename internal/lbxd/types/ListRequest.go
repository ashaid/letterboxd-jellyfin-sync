package types

type Entries struct {
	Film             string `json:"film"`
	Action           string `json:"action"` // only see ADD for now, could be change to enum
	ContainsSpoilers bool   `json:"containsSpoilers"`
}

// this is what Update should be if lbxd knew what they were doing
type Update struct {
	Version     uint8     `json:"version"`
	Name        string    `json:"boxdItCode"`
	Description string    `json:"sharingUrl"`
	Tags        []string  `json:"tags"`
	Published   bool      `json:"published"`
	SharePolicy string    `json:"sharePolicy"`
	Ranked      bool      `json:"ranked"`
	Entries     []Entries `json:"entries"`
}

type ListRequest struct {
	Csrf       string `json:"csrf"`
	FilmListId string `json:"filmListId"`
	// lol lbxd this json object is passed as a full string
	Update string `json:"update"`
}
