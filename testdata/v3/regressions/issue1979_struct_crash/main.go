package main

import "time"

// @title Issue 1979 Repro API
// @version 1.0

// TODO: the upstream crash points at an external repo's `sourceAge` struct that
// is not published in the issue. This is a best-effort minimal struct combining
// a duration field, a pointer, and an embedded anonymous struct to exercise the
// same v3 schema-generation path. Refine if the real trigger surfaces.
type sourceAge struct {
	Age     time.Duration `json:"age"`
	Updated *time.Time    `json:"updated"`
	Meta    struct {
		Tag string `json:"tag"`
	} `json:"meta"`
}

// Handle godoc
// @Summary  struct that crashed v3 generation
// @Produce  json
// @Success  200 {object} sourceAge
// @Router   /source [get]
func Handle() {}

func main() {}
