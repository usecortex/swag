package extension

import "net/http"

// @Summary Get items
// @x-custom-ext {"key":"value"}
// @Success 200 {string} string "ok"
// @Router /items [get]
func GetItems(w http.ResponseWriter, r *http.Request) {}
