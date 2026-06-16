package main

// @title Issue 2153 Repro API
// @version 1.0

// CreateInput has a float64 field with a fractional maximum, which failed to
// parse (strconv.Atoi on "0.1").
type CreateInput struct {
	Factor *float64 `json:"factor" validate:"required" minimum:"0" maximum:"0.1" example:"0.04"`
}

// Create godoc
// @Accept  json
// @Produce json
// @Param   Payload body CreateInput true "Request Body"
// @Success 200 {string} string "ok"
// @Router  /something/create [post]
func Create() {}

func main() {}
