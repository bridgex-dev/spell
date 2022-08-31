package spell

import (
	"github.com/justinas/nosurf"
	"net/http"
)

func (e *Engine) Handler(next http.Handler) http.Handler {
	if e.EnableCSRFToken {
		return e.withCsrf(next)
	}

	return e.middleware(next)
}

func (e *Engine) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, ContextOptions{e.Cookies, e.EnableCSRFToken})
		e.addContext(ctx)

		next.ServeHTTP(w, r)

		ctx.makeResponse()
		e.removeContext(ctx.id)
	})
}

func (e *Engine) withCsrf(next http.Handler) http.Handler {
	return nosurf.NewPure(e.middleware(next))
}
