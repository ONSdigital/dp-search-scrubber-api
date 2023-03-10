package api

type ErrorResp struct {
	errors   []Errors
	trace_id string
}

type Errors struct {
	error_code string
	message    string
}
