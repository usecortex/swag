package main

// @title Issue 2139 Repro API
// @version 1.0

// MotdMessage is the response type.
type MotdMessage struct {
	Message string `json:"message"`
}

// AddHandlers registers routes; the swagger comments live inside the function
// body and must be parsed under --v3.1 + --parseFuncBody.
func AddHandlers() {
	// @Summary  Get all Message of the Day objects
	// @Tags     Motd
	// @Success  200 {array} MotdMessage
	// @Router   /motd [get]
	_ = "get"

	// @Summary  Create a new Message of the Day object
	// @Tags     Motd
	// @Param    request body MotdMessage true "Partial Motd values for creating"
	// @Success  201 {object} MotdMessage
	// @Router   /motd [post]
	_ = "post"
}

func main() {
	AddHandlers()
}
