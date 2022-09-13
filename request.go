package spell

import (
	"context"
	"net/http"
)

const (
	ContextKey = "spell"
)

func withSpellContext(r *http.Request, c *Context) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), ContextKey, c))
}

func FromRequest(r *http.Request) *Context {
	return r.Context().Value(ContextKey).(*Context)
}
