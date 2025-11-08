package rawg

type RAWGGame struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Released string `json:"released"`
}

type RAWGSearchResponse struct {
	Results []RAWGGame `json:"results"`
}
