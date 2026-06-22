// Provides composable HTTP middleware
package middleware

import (
	"net/http"
)

// Composes middlewares right-to-left so the call order matches the declared order (first declared = outermost)
func Chain(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
