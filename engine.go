package spell

import (
	"github.com/evanw/esbuild/pkg/api"
	"net/http"
)

type Engine struct {
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
	return &Engine{
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
		Logger:        NewDefaultLogger(),
		CookieManager: NewCookieManager("qTvsuYaLZTOmaRz8", "qTvsuYaLZTOmaRz8"),
	}
}
