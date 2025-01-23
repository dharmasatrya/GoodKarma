package dto

type ResponseOK struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseError struct {
	Message string `json:"message"`
}
