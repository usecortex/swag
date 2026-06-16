package main

// @title Issue 1949 Repro API
// @version 1.0

// Thing has a field with an example, which on 3.1 should be emitted under the
// `examples` keyword rather than the deprecated singular `example`.
type Thing struct {
	Name string `json:"name" example:"widget"`
}

// Handle godoc
// @Produce json
// @Success 200 {object} Thing
// @Router  /thing [get]
func Handle() {}

func main() {}
