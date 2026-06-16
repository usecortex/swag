package main

// @title Issue 2142 / 2086 Repro API
// @version 1.0
// @BasePath /api

// CreateRequest is the request body type.
type CreateRequest struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

// Response is the response body type.
type Response struct {
	ID int `json:"id"`
}

// CreateThing godoc
// @Summary Create thing
// @Accept  json
// @Produce json
// @Param   body body CreateRequest true "Request"
// @Success 200 {object} Response
// @Router  /things [post]
func CreateThing() {}

func main() {}
