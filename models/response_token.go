package models

type ResponseToken struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   Token  `json:"data"`
}
