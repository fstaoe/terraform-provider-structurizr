package model

// GenericResponse represents a response from structurizr
type GenericResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Revision int64  `json:"revision"`
}
