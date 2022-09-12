package spell

import (
	"encoding/gob"
	"github.com/evanw/esbuild/pkg/api"
	"net/http"
)

type Engine struct {
	contexts        map[string]*Context
	Cookies         CookiesOptions
	CookieManager   CookieManager
	EnableCSRFToken bool
	EsbuildOptions  api.BuildOptions
	Logger          Logger
}

var defaultCookies = CookiesOptions{
	Path:     "/",
	Domain:   "",
	SameSite: http.SameSiteStrictMode,
	Secure:   true,
	HttpOnly: true,
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
		Logger: NewDefaultLogger(),
		CookieManager: NewCookieManager(
			"hash-key",
			"block-key",
		),
	}
}

func (e *Engine) addContext(c *Context) {
	e.Logger.Logf(DebugLevel, "Adding context with id: %s", c.id)
	e.contexts[c.id] = c
}

func (e *Engine) removeContext(id string) {
	e.Logger.Logf(DebugLevel, "Removing context with id: %s", id)
	delete(e.contexts, id)
}

func (e *Engine) GetContext(r *http.Request) *Context {
	id := r.Header.Get(HeaderId)
	return e.contexts[id]
}
