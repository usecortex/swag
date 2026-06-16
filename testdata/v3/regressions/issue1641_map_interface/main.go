package main

// @title Issue 1641 Repro API
// @version 1.0

// HandleA godoc
// @Produce json
// @Success 200 {object} map[string]interface{} "result A"
// @Router  /a [get]
func HandleA() {}

// HandleB godoc
// @Produce json
// @Success 200 {object} map[string]interface{} "result B"
// @Router  /b [get]
func HandleB() {}

func main() {}
