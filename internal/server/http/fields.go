package server

type ResponseSuccess struct {
	Status     string `json:"status"`
	StatusCode int    `json:"code"`
}

type ResponseError struct {
	Error      string `json:"error"`
	StatusCode int    `json:"code"`
}

type CheckRequestIn struct {
	IP       string `json:"ip"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AddIPIn struct {
	Cidr string `json:"cidr"`
}
type DeleteIPIn struct {
	Cidr string `json:"cidr"`
}

type ClearBucketIn struct {
	IP    string `json:"ip"`
	Login string `json:"login"`
}
