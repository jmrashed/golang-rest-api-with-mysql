package model

type Employee struct {
	Id   string `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
	City string `form:"city" json:"city"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Employee
}
