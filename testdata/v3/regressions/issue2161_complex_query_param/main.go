package main

// @title Issue 2161 Repro API
// @version 1.0

// ComplexType is a struct used as a query parameter; each field must expand to
// its own query parameter (v1 behavior).
type ComplexType struct {
	One string `json:"one"`
	Two string `json:"two"`
}

// Result is the response type.
type Result struct {
	OK bool `json:"ok"`
}

// Handle godoc
// @Produce json
// @Success 200 {array} Result
// @Param   root query ComplexType false "descr"
// @Router  /api/v1/orgstruct [get]
func Handle() {}

func main() {}
