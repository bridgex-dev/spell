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
		)
		if err != nil {
			e.Logger.Logf(ErrorLevel, "Error creating context: %s", err)
			return
		}

		e.addContext(ctx)
		defer e.removeContext(ctx.id)

		next.ServeHTTP(w, r)

		if err = ctx.makeResponse(); err != nil {
			e.Logger.Logf(ErrorLevel, "Error making response: %s", err)
		}
	})
}

func (e *Engine) withCsrf(next http.Handler) http.Handler {
	return nosurf.NewPure(e.middleware(next))
}
