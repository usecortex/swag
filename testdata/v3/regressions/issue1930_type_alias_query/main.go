package main

// @title Issue 1930 Repro API
// @version 1.0

// PriceType is a defined string type (alias-like) used inside a query struct.
type PriceType string

// RequestParam embeds a defined type as a field, which triggered a nil deref in
// parseStructFieldV3 on --v3.1.
type RequestParam struct {
	Price PriceType `form:"price" json:"price"`
}

// Handle godoc
// @Summary  query with a defined-type field
// @Param    request query RequestParam false "Query Params"
// @Success  200 {string} string "ok"
// @Router   /price [get]
func Handle() {}

func main() {}
