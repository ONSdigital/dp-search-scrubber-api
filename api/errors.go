package api

type ErrorResp struct {
	Errors  []Errors
	TraceID string
}

type Errors struct {
	ErrorCode string
	Message   string
}
