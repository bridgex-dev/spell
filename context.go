package spell

import (
	"github.com/justinas/nosurf"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type Context struct {
	w        http.ResponseWriter
	r        *http.Request
	id       string
	Flash    Flash
	Session  Session
	ViewData map[string]interface{}
	options  ContextOptions
}

type ContextOptions struct {
	cookies   CookiesOptions
	csrfToken bool
}

func NewContext(w http.ResponseWriter, r *http.Request, options ContextOptions) *Context {
	id := markRequest(r)
	flash := getFlash(r)
	session := getSession(r)
	viewData := make(map[string]interface{})

	viewData["flash"] = flash

	if options.csrfToken {
		token := nosurf.Token(r)
		viewData[CSRFToken] = token
		viewData["renderToken"] = renderToken(token)
	}

	return &Context{
		id:       id,
		w:        w,
		r:        r,
		options:  options,
		Flash:    NewFlash(),
		Session:  session,
		ViewData: viewData,
	}
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.w, c.r, url, http.StatusFound)
}

func markRequest(r *http.Request) string {
	id := uuid.NewV4().String()
	r.Header.Set(HeaderId, id)

	return id
}

func getFlash(r *http.Request) Flash {
	flash := NewFlash()

	value, err := getCookieValue(r, FlashCookie)
	if err != nil {
		return flash
	}

	_ = flash.decode(value)

	return flash
}

func getSession(r *http.Request) Session {
	session := NewSession()

	value, err := getCookieValue(r, SessionCookie)
	if err != nil {
		return session
	}

	_ = session.decode(value)

	return session
}

func (c *Context) makeResponse() {
	if len(c.Flash) > 0 {
		flash, err := c.Flash.encode()
		if err == nil {
			setCookies(c.w, FlashCookie, flash, c.options.cookies)
		}
	} else {
		setCookies(c.w, FlashCookie, "", c.options.cookies)
	}

	session, err := c.Session.encode()
	if err == nil {
		setCookies(c.w, SessionCookie, session, c.options.cookies)
	}
}
