package models

type Character struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Thumbnail   struct {
		Path      string `json:"path"`
		Extension string `json:"extension"`
	} `json:"thumbnail"`
	Comics struct {
		Available int `json:"available"`
	} `json:"comics"`
	Series struct {
		Available int `json:"available"`
	} `json:"series"`
	Stories struct {
		Available int `json:"available"`
	} `json:"stories"`
	Events struct {
		Available int `json:"available"`
	} `json:"events"`
	IsFavorite bool `json:"isFavorite"`
}

type CharacterResponse struct {
	Data struct {
		Results []Character `json:"results"`
		Total   int         `json:"total"`
	} `json:"data"`
}
