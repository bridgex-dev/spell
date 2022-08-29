package spell

import (
	"encoding/gob"
	"net/http"
)

type Engine struct {
	contexts        map[string]*Context
	Cookies         CookiesOptions
	EnableCSRFToken bool
}

func NewEngine() *Engine {
	gob.Register(map[string]interface{}{})

	return &Engine{
		contexts:        make(map[string]*Context),
		Cookies:         defaultCookies,
		EnableCSRFToken: true,
	}
}

func (hw *Engine) addContext(c *Context) {
	hw.contexts[c.id] = c
}

func (hw *Engine) removeContext(id string) {
	delete(hw.contexts, id)
}

func (hw *Engine) GetContext(r *http.Request) *Context {
	id := r.Header.Get(HeaderId)
	return hw.contexts[id]
}
