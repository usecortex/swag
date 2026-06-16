package main

// @title Issue 1855 Repro API
// @version 1.0

// Color is an enum reused across multiple fields; it should be emitted once as
// a reusable component schema and referenced, not inlined per field.
type Color string

// ENUM values for Color.
const (
	ColorRed   Color = "red"
	ColorGreen Color = "green"
	ColorBlue  Color = "blue"
)

// Palette references the reusable enum twice (one nullable).
type Palette struct {
	Primary   Color  `json:"primary"`
	Secondary *Color `json:"secondary"`
}

// Handle godoc
// @Produce json
// @Success 200 {object} Palette
// @Router  /palette [get]
func Handle() {}

func main() {}
