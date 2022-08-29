package spell

import (
	"github.com/justinas/nosurf"
	"net/http"
)

func (hw *Engine) Handler(next http.Handler) http.Handler {
	if hw.EnableCSRFToken {
		return hw.withCsrf(next)
	}

	return hw.middleware(next)
}

func (hw *Engine) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, ContextOptions{hw.Cookies, hw.EnableCSRFToken})
		hw.addContext(ctx)

		next.ServeHTTP(w, r)

		ctx.makeResponse()
		hw.removeContext(ctx.id)
	})
}

func (hw *Engine) withCsrf(next http.Handler) http.Handler {
	return nosurf.NewPure(hw.middleware(next))
}
