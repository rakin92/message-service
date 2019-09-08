package models

// Response : sample request response
type Response struct {
	Status  string
	code    int
	message string
	data    map[string]string
	failure error
}
