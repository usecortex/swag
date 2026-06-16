package main

// @title Issue 1934 Repro API
// @version 1.0

// EC is an enum string type.
type EC string

// ENUM(a,b,c,d)
const (
	ECA EC = "a"
	ECB EC = "b"
	ECC EC = "c"
	ECD EC = "d"
)

// OtherEnum is a second enum string type.
type OtherEnum string

// ENUM(aa,bb,cc,dd)
const (
	OtherAA OtherEnum = "aa"
	OtherBB OtherEnum = "bb"
)

// Complex mixes a scalar enum, an enum slice, and an int slice.
type Complex struct {
	F1 EC          `json:"f_1"`
	F2 []OtherEnum `json:"f_2"`
	F3 []int       `json:"f_3"`
}

// Test2 godoc
// @Produce json
// @Success 200 {object} Complex
// @Router  /test2 [get]
func Test2() {}

func main() {}
