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
		ctx, err := NewContext(
			w,
			r,
			ContextOptions{e.Cookies, e.EnableCSRFToken},
			e.Logger,
			e.CookieManager,
		)
		if err != nil {
			e.Logger.Logf(ErrorLevel, "Error creating context: %s", err)
			return
		}

		r = withSpellContext(r, ctx)

		next.ServeHTTP(w, r)

		if err = ctx.WriteHeader(); err != nil {
			e.Logger.Logf(ErrorLevel, "Error making response: %s", err)
		}
	})
}

func (e *Engine) withCsrf(next http.Handler) http.Handler {
	return nosurf.NewPure(e.middleware(next))
}
