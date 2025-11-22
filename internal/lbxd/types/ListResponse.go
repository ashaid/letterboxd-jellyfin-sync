package types

type ListItem struct {
	Id         uint32 `json:"id"`
	BoxdItCode string `json:"boxdItCode"`
	SharingUrl string `json:"sharingUrl"`
	Name       string `json:"name"`
	Version    uint8  `json:"version"`
}

type ListResponse struct {
	Result      bool       `json:"result"`
	Csrf        string     `json:"csrf"`
	Messages    []string   `json:"messages"`
	ErrorCodes  []string   `json:"errorCodes"`
	ErrorFields []string   `json:"errorFields"`
	NewList     bool       `json:"newList"`
	ListId      uint32     `json:"listId"`
	Version     uint8      `json:"version"`
	Name        string     `json:"name"`
	List        []ListItem `json:"list"`
	EditFormURL string     `json:"editFormURL"`
	URL         string     `json:"url"`
}
