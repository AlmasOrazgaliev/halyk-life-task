package models

type Request struct {
	Id      uint    `json:"id"`
	Method  string  `json:"method"`
	Url     string  `json:"url"`
	Headers Headers `json:"headers"`
}

type Headers struct {
	Authentication string `json:"Authentication"`
}

type Response struct {
	Id      int      `json:"id"`
	Status  int      `json:"status"`
	Headers []string `json:"headers"`
	Length  int      `json:"length"`
}
