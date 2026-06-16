package main

import "time"

// @title Issue 2078 Repro API
// @version 1.0

// Date is a custom type embedding time.Time.
type Date struct {
	time.Time
}

// MyModel has a slice of the custom type overridden via swaggertype.
type MyModel struct {
	Dates []Date `json:"dates" swaggertype:"array,string" format:"date" example:"2025-01-01"`
}

// Handle godoc
// @Summary  custom type array with swaggertype override
// @Produce  json
// @Success  200 {object} MyModel
// @Router   /dates [get]
func Handle() {}

func main() {}
