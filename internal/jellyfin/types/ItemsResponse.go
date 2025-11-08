package types

import "time"

type UserData struct {
	Played bool `json:"Played"`
}

type Item struct {
	Name            string    `json:"Name"`
	PremiereDate    time.Time `json:"PremiereDate"`
	CriticRating    uint8     `json:"CriticRating"`
	CommunityRating float32   `json:"CommunityRating"`
	UserData        UserData  `json:"UserData"`
}

type ItemsResponse struct {
	Items            []Item `json:"Items"`
	TotalRecordCount int    `json:"TotalRecordCount"`
	StartIndex       int    `json:"StartIndex"`
}
