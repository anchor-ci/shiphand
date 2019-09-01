package app

type JobMetadata struct {
	Id        string `json:"id"`
	State     string `json:"state"`
	HistoryId string
}
