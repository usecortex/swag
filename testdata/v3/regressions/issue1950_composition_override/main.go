package main

// @title Issue 1950 Repro API
// @version 1.0

// Response is a generic-style envelope whose Data field is overridden in the
// @Success annotation via Response{data=GetData1Response}.
type Response struct {
	Code string      `json:"code" example:"0"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

// GetData1Response is the concrete type overriding Data.
type GetData1Response struct {
	Field1 string `json:"field1" example:"Field1"`
	Field2 string `json:"field2" example:"Field2"`
}

// GetData1 godoc
// @Produce json
// @Success 200 {object} Response{data=GetData1Response} "desc"
// @Router  /testapi/data1 [get]
func GetData1() {}

func main() {}
