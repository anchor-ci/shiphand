package main

type History struct {
    FailureText string `json:"failureText"`
	Text      string `json:"text"`
	Failed    bool   `json:"failed"`
	Succeeded bool   `json:"succeeded"`
}
