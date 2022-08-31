package spell

import (
	"encoding/gob"
	"github.com/evanw/esbuild/pkg/api"
	"net/http"
)

type Engine struct {
	contexts        map[string]*Context
	Cookies         CookiesOptions
	EnableCSRFToken bool
	EsbuildOptions  api.BuildOptions
}

func NewEngine() *Engine {
	gob.Register(map[string]interface{}{})

	return &Engine{
		contexts:        make(map[string]*Context),
		Cookies:         defaultCookies,
		EnableCSRFToken: true,
		EsbuildOptions: api.BuildOptions{
			EntryPoints: []string{
				"web/assets/src/index.ts",
				"web/assets/stylesheets/main.css",
			},
			Outdir:            "dist/public",
			Bundle:            true,
			Write:             true,
			Sourcemap:         api.SourceMapNone,
			MinifyWhitespace:  true,
			MinifyIdentifiers: true,
			MinifySyntax:      true,
		},
	}
}

func (e *Engine) addContext(c *Context) {
	e.contexts[c.id] = c
}

func (e *Engine) removeContext(id string) {
	delete(e.contexts, id)
}

func (e *Engine) GetContext(r *http.Request) *Context {
	id := r.Header.Get(HeaderId)
	return e.contexts[id]
}
