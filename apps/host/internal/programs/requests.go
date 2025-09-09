package programs

type CreateProgramRequest struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	HostID      string `json:"hostId"`
}
