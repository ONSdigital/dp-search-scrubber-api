package api

type ErrorResp struct {
	Errors   []Errors
	Trace_id string
}

type Errors struct {
	Error_code string
	Message    string
}
