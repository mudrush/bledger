package model

// StandardResponse is a struct for returning normalized api responses
type StandardResponse struct {
	APIVersion APIVersion `json:"api_version"`
	Response   any        `json:"response"`
}

// StandardErrorResponse is a struct for returning normalized api error responses
type StandardErrorResponse struct {
	APIVersion     APIVersion `json:"api_version"`
	Error          any        `json:"error"`
	InternalErrMsg string     `json:"internal_error_message"`
}
