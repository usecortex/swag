package main

// @title Issue 1911 Repro API
// @version 1.0

// TODO: the original input for the getFuncDoc nil-deref was not published. This
// is a best-effort fixture with router comments inside nested function literals
// (parsed because the test sets ParseFuncBody = true), which exercises the
// getFuncDoc recursion path from the stack trace.
func setupRoutes() {
	register := func() {
		inner := func() {
			// @Summary  nested-funclit route
			// @Success  200 {string} string "ok"
			// @Router   /nested [get]
		}
		_ = inner
	}
	_ = register
}

func main() {
	setupRoutes()
}
