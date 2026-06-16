package main

// @title Issue 2080 Repro API
// @version 1.0

// NickType is an enum string type.
type NickType string

// ENUM values for NickType.
const (
	NickOne NickType = "one"
	NickTwo NickType = "two"
)

// Wrapper is a generic struct whose field is an array of the type parameter.
type Wrapper[T any] struct {
	Values []T `json:"values"`
}

// Handle godoc
// @Summary  generic struct of enum array
// @Produce  json
// @Success  200 {object} Wrapper[NickType]
// @Router   /nick [get]
func Handle() {}

func main() {}
