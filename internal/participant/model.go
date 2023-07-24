package participant

type Participant struct {
	ID           int    `json:"id"`
	Uuid         string `json:"uuid"`
	Organization string `json:"organization"`
}
